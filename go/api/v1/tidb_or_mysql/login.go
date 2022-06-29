package tidb_or_mysql

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/nooncall/owls/go/service/tidb_or_mysql/admin"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/auth"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/task"
)

const (
	roleAdmin = "ADMIN"
	roleUser  = "USER"
)

//todo, delete this file

func Login(ctx *gin.Context) Resp {
	f := "Login()-->"
	var authParam auth.Claims
	if err := ctx.BindJSON(&authParam); err != nil {
		return Resp{Message: fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), Code: 500}
	}

	//token, err := auth.Login(authParam.Username, authParam.Password)
	//if err != nil {
	//	return Resp{Message: fmt.Sprintf("%s, login failed :%s", f, err.Error()), Code: code.InternalErr}
	//}
	return Resp{Data: "test token"}
}

func RoleGet(ctx *gin.Context) Resp {
	f := "RoleGet()-->"
	user := ctx.GetString("user")

	isAdmin, err := admin.IsAdmin(user)
	if err != nil {
		return Resp{Message: fmt.Sprintf("%s,check admin failed :%s ", f, err.Error()), Code: 500}
	}

	role := roleUser
	if isAdmin {
		role = roleAdmin
	}

	return Resp{Data: struct {
		Role string `json:"role"`
		Name string `json:"name"`
	}{
		Role: role,
		Name: user,
	}}
}

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*		token := c.GetHeader("token")
				claims, err := auth.ParseToken(token)
				if err != nil {
					c.JSON(http.StatusOK, "please login")
					c.Abort()
					return
				}

				if err := claims.Valid(); err != nil {
					c.JSON(http.StatusOK, "place login")
					c.Abort()
					return
				}*/

		c.Set("user", "test")
		c.Next()
	}
}

func OnlyDbaOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(string)
		isDba, err := task.IsDba(user)
		if err != nil {
			c.JSON(http.StatusOK, Resp{
				Code:    500,
				Message: fmt.Sprintf("check dba auth failed, %s", err.Error()),
			})
			c.Abort()
			return
		}

		if !isDba {
			c.JSON(http.StatusOK, Resp{
				Code:    500,
				Message: "only admin can operate cluster and administrator module",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
