package controller

import (
	"bikefest/pkg/bootstrap"
	"bikefest/pkg/model"
	"fmt"
	"github.com/gin-gonic/gin"
	social "github.com/kkdai/line-login-sdk-go"
	"log"
	"net/http"
	"strconv"
)

func NewOAuthController(lineSocialClient *social.Client, env *bootstrap.Env, userSvc model.UserService) *OAuthController {
	return &OAuthController{
		lineSocialClient: lineSocialClient,
		userSvc:          userSvc,
		env:              env,
	}
}

type OAuthController struct {
	lineSocialClient *social.Client
	userSvc          model.UserService
	env              *bootstrap.Env
}

// http://localhost:8000/line-login/auth
func (ctrl *OAuthController) LineLogin(c *gin.Context) {
	originalUrl := c.Request.Referer() + c.Query("redirect_path")

	// remove the duplicate slash
	originalUrl = strings.ReplaceAll(originalUrl, "//", "/")

	log.Println("originalUrl:", originalUrl)
	serverURL := ctrl.env.Line.ServerUrl
	scope := "profile openid" //profile | openid | email
	// TODO: add some random string to state
	state := originalUrl + "$" + social.GenerateNonce()
	if len(state) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Response{
			Msg: "login with the wrong way, please try again in official website",
		})
		return
	}
	nonce := social.GenerateNonce()
	redirectURL := fmt.Sprintf("%s/line-login/callback", serverURL)
	targetURL := ctrl.lineSocialClient.GetWebLoinURL(redirectURL, state, scope, social.AuthRequestOptions{Nonce: nonce, Prompt: "consent"})
	c.SetCookie("state", state, 3600, "/", "", false, true)
	c.Redirect(http.StatusFound, targetURL)
}

func (ctrl *OAuthController) LineLoginCallback(c *gin.Context) {
	serverURL := ctrl.env.Line.ServerUrl
	code := c.Query("code")
	state := c.Query("state")
	stateInCookie, err := c.Cookie("state")
	if err != nil || stateInCookie != state {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Response{
			Msg: "state cookie is invalid",
		})
		return
	}
	log.Println("code:", code, " stateInCookie:", stateInCookie)
	// TODO: remove the random string from state
	frontendURL := strings.Split(stateInCookie, "$")[0]
	token, err := ctrl.lineSocialClient.GetAccessToken(fmt.Sprintf("%s/line-login/callback", serverURL), code).Do()
	if err != nil {
		log.Println("RequestLoginToken err:", err)
		return
	}
	log.Println("access_token:", token.AccessToken, " refresh_token:", token.RefreshToken)

	var payload *social.Payload
	//if len(token.IDToken) == 0 {
	//	// User don't request openID, use access token to get user profile
	//	log.Println(" token:", token, " AccessToken:", token.AccessToken)
	//	res, err := ctrl.lineSocialClient.GetUserProfile(token.AccessToken).Do()
	//	if err != nil {
	//		log.Println("GetUserProfile err:", err)
	//		return
	//	}
	//	payload = &social.Payload{
	//		Name:    res.DisplayName,
	//		Picture: res.PictureURL,
	//	}
	//} else {
	//Decode token.IDToken to payload
	payload, err = token.DecodePayload(ctrl.env.Line.ChannelID)
	if err != nil {
		log.Println("DecodeIDToken err:", err)
		return
	}
	//}
	log.Printf("payload: %#v", payload)

	//c.JSON(http.StatusOK, gin.H{
	//	"status": "Success",
	//	"data":   payload,
	//})

	user := &model.User{
		ID:   payload.Sub,
		Name: payload.Name,
	}

	err = ctrl.userSvc.CreateFakeUser(c, user)

	if err != nil {
		log.Printf("user with id %s already exists", user.ID)
	}

	accessToken, err := ctrl.userSvc.CreateAccessToken(c, user, ctrl.env.JWT.AccessTokenSecret, ctrl.env.JWT.AccessTokenExpiry)
	if err != nil {
		log.Printf("failed to create access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Failed",
			"message": "failed to create access token",
		})
		return
	}

	refreshToken, err := ctrl.userSvc.CreateRefreshToken(c, user, ctrl.env.JWT.RefreshTokenSecret, ctrl.env.JWT.RefreshTokenExpiry)
	if err != nil {
		log.Printf("failed to create refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Failed",
			"message": "failed to create refresh token",
		})
		return
	}

	// set to cookie
	c.SetCookie("access_token", strconv.FormatInt(ctrl.env.JWT.AccessTokenExpiry, 10), 3600, "/", "", false, true)
	c.SetCookie("refresh_token", strconv.FormatInt(ctrl.env.JWT.AccessTokenExpiry, 10), 3600, "/", "", false, true)
	// redirect to frontend
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s?access_token=%s&refresh_token=%s&redirect_url=%s", frontendURL, accessToken, refreshToken, frontendURL))
}
