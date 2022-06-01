package auth

import (
	"encoding/json"
	"fmt"
	"github.com/qingfeng777/owls/server/utils"
	"net/http"
)

type NetAuthToolImpl struct {
}

var NetAuthService NetAuthToolImpl

var getReviewParam = `{
    "field": 5,
    "subordinate_op_name_list": ["%s"]
}`

type NetResp struct {
	Data struct {
		Ent struct {
			Items []struct {
				OpName string `json:"op_name"`
				Uid    string `json:"uid"`
			} `json:"items"`
		} `json:"ent"`
	} `json:"data"`
}

func (NetAuthToolImpl) GetReviewer(userName string) (reviewerName string, err error) {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Authorization", "not set")

	respData, err := utils.DoHttpReq(http.MethodPost, "", fmt.Sprintf(getReviewParam, userName), header)
	if err != nil {
		return "", err
	}

	var resp NetResp
	if err = json.Unmarshal(respData, &resp); err != nil {
		return "", fmt.Errorf("unmarshal reviewer api resp err: %s", err.Error())
	}
	if len(resp.Data.Ent.Items) < 1 {
		return "", fmt.Errorf("get reviewer by api no response content")
	}
	return resp.Data.Ent.Items[0].OpName, nil
}

var isDBAparam = `{"busiid":%d}`

func (NetAuthToolImpl) IsDba(userName string) (isDba bool, err error) {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Authorization", "")

	resData, err := utils.DoHttpReq(http.MethodPost, "", fmt.Sprintf(isDBAparam, ""), header)
	if err != nil {
		return false, err
	}

	var resp NetResp
	if err = json.Unmarshal(resData, &resp); err != nil {
		return false, fmt.Errorf("unmarshal dba api resp err: %s", err.Error())
	}
	if len(resp.Data.Ent.Items) < 1 {
		return false, fmt.Errorf("get dba member by api no response content")
	}

	for _, v := range resp.Data.Ent.Items {
		if userName == v.Uid {
			return true, nil
		}
	}
	return false, nil
}
