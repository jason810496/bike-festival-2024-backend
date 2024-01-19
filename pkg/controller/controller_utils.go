package controller

import (
	"bikefest/pkg/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RetrievePagination(c *gin.Context) (page, limit uint64) {
	pageStr := c.Query("page")

	limitStr := c.Query("limit")

	// convert string to uint64
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	limit, err = strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		limit = 10
	}

	return
}

// RetrieveIdentity retrieves the identity of the user from the context.
// raise: if true, raise a http error if the identity does not exist
func RetrieveIdentity(c *gin.Context, raise bool) (identity *model.Identity, exist bool) {
	id, exist := c.Get("identity")
	if !exist {
		if raise {
			// raise not login error
			c.AbortWithStatusJSON(401, model.Response{
				Msg: "not login",
			})
		}
		return nil, false
	}
	identity = id.(*model.Identity)
	return
}
