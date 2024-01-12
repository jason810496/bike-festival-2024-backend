package controller

import (
	"net/http"
	"strconv"
	"strings"

	"bikefest/pkg/bootstrap"
	"bikefest/pkg/model"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userSvc model.UserService
	env     *bootstrap.Env
}

func NewUserController(userSvc model.UserService, env *bootstrap.Env) *UserController {
	return &UserController{
		userSvc: userSvc,
		env:     env,
	}
}

func (ctrl *UserController) Profile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	profile, err := ctrl.userSvc.GetUserByID(c, userID.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Msg: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Data: profile,
	})
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	userID := c.Param("user_id")
	user, err := ctrl.userSvc.GetUserByID(c, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Msg: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Msg:  "get user by id",
		Data: user,
	})
}

func (ctrl *UserController) RefreshToken(c *gin.Context) {
	var request model.RefreshTokenRequest

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Response{
			Msg: err.Error(),
		})
		return
	}

	identity, err := ctrl.userSvc.VerifyRefreshToken(c, request.RefreshToken, ctrl.env.JWT.RefreshTokenSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
			Msg: err.Error(),
		})
		return
	}

	// TODO: DO we need to check if the user is still active or any other things?

	accessToken, err := ctrl.userSvc.CreateAccessToken(c, identity, ctrl.env.JWT.AccessTokenSecret, ctrl.env.JWT.AccessTokenExpiry)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Msg: err.Error(),
		})
		return
	}

	refreshToken, err := ctrl.userSvc.CreateRefreshToken(c, identity, ctrl.env.JWT.RefreshTokenSecret, ctrl.env.JWT.RefreshTokenExpiry)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Msg: err.Error(),
		})
		return
	}

	loginResponse := model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, model.Response{
		Data: loginResponse,
	})
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
	// page, limit := RetrievePagination(c)

	c.JSON(http.StatusOK, model.Response{
		Msg:  "get users",
		Data: interface{}(nil),
	})
}

func (ctrl *UserController) Logout(c *gin.Context) {
	// TODO: need to discuss where to read the token from (header or body or cookie)
	authHeader := c.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{
			Msg: "Invalid token format (length different from 2)",
		})
		return
	}
	authToken := bearerToken[1]
	err := ctrl.userSvc.Logout(c, &authToken, ctrl.env.JWT.AccessTokenSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Msg: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Msg: "logout success",
	})
}

func (ctrl *UserController) FakeLogin(c *gin.Context) {
	userID := c.Param("user_id")

	accessToken, err := ctrl.userSvc.CreateAccessToken(c, &model.User{ID: userID}, ctrl.env.JWT.AccessTokenSecret, ctrl.env.JWT.AccessTokenExpiry)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Msg: err.Error(),
		})
		return
	}

	refreshToken, err := ctrl.userSvc.CreateRefreshToken(c, &model.User{ID: userID}, ctrl.env.JWT.RefreshTokenSecret, ctrl.env.JWT.RefreshTokenExpiry)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Msg: err.Error(),
		})
		return
	}

	loginResponse := model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, model.Response{
		Data: loginResponse,
	})

	// set to cookie
	c.SetCookie("access_token", strconv.FormatInt(ctrl.env.JWT.AccessTokenExpiry, 10), 3600, "/", "", false, true)
	c.SetCookie("refresh_token", strconv.FormatInt(ctrl.env.JWT.AccessTokenExpiry, 10), 3600, "/", "", false, true)
}

func (ctrl *UserController) FakeRegister(c *gin.Context) {
	var request model.CreateFakeUserRequest

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.Response{
			Msg: err.Error(),
		})
		return
	}

	user := &model.User{
		Name: request.Name,
	}

	err := ctrl.userSvc.CreateFakeUser(c, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.Response{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Msg: "fake register success",
	})
}
