package response

import (
	"github.com/qingfeng777/owls/server/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
