package example

import "github.com/nooncall/owls/go/model/system"

type ExcelInfo struct {
	FileName string               `json:"fileName"` // 文件名
	InfoList []system.SysBaseMenu `json:"infoList"`
}
