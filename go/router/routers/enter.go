// all follow routers will be added here

package routers

import "github.com/gin-gonic/gin"

type Routers struct{}

func (s *Routers) InitApiRouter(Router *gin.RouterGroup) {
	s.InitAuthRouter(Router)
	s.InitTaskRouter(Router)
}

var Router Routers
