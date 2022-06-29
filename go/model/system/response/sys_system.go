package response

import "github.com/nooncall/owls/go/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
