package response

import (
	"github.com/nooncall/owls/go/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
