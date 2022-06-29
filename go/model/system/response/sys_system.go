package response

import "github.com/qingfeng777/owls/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
