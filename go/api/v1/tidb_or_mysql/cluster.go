package tidb_or_mysql

import (
	"fmt"
	"github.com/qingfeng777/owls/server/model/common/request"
	"github.com/qingfeng777/owls/server/utils"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qingfeng777/owls/server/model/common/response"
	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/db_info"
)

type ClusterApi struct{}

func (clusterApi *ClusterApi) ListDB(ctx *gin.Context) {
	f := "listDB()-->"

	var filter bool
	filterStr := ctx.Query("filter")
	if filterStr == "true" {
		filter = true
	}

	cluster := ctx.Query("cluster")
	if cluster == "" {
		response.FailWithMessage("need cluster param: cluster name", ctx)
		return
	}

	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.FailWithMessage("get user err: "+err.Error(), ctx)
		return
	}

	dbInfo, err := db_info.ListDB(cluster, claims.ID, filter)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,list db failed :%s ", f, err.Error()), ctx)
		return
	}

	response.OkWithData(dbInfo, ctx)
}

func (clusterApi *ClusterApi) ListTable(ctx *gin.Context) {
	f := "ListTable()-->"

	cluster := ctx.Query("cluster")
	if cluster == "" {
		response.FailWithMessage("need cluster param: cluster", ctx)
		return
	}
	db := ctx.Query("db")
	if db == "" {
		response.FailWithMessage("need db param: db", ctx)
		return
	}

	table, err := db_info.ListTable(cluster, db)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,list table failed :%s ", f, err.Error()), ctx)
		return
	}

	response.OkWithData(table, ctx)
}

func (clusterApi *ClusterApi) ListCluster(ctx *gin.Context) {
	f := "ListCluster()-->"

	var pageInfo request.SortPageInfo
	ctx.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	clusters, err := db_info.ListClusterForUI(pageInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,list cluster failed :%s ", f, err.Error()), ctx)
		return
	}

	response.OkWithData(ListData{
		List:     clusters,
		Total:    int64(len(clusters)),
		More:     false,
		Offset:   0,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, ctx)
}

func (clusterApi *ClusterApi) ListClusterName(ctx *gin.Context) {
	f := "ListClusterName()-->"

	var filter bool
	filterStr := ctx.Query("filter")
	if filterStr == "true" {
		filter = true
	}

	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.FailWithMessage("get user err: "+err.Error(), ctx)
		return
	}

	clusters, err := db_info.ListClusterName(claims.ID, filter)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,list cluster failed :%s ", f, err.Error()), ctx)
		return
	}

	response.OkWithData(clusters, ctx)
}

func (clusterApi *ClusterApi) AddCluster(ctx *gin.Context) {
	f := "AddCluster()-->"

	var cluster db_info.OwlCluster
	if err := ctx.BindJSON(&cluster); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	id, err := db_info.AddCluster(&cluster)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,add cluster failed :%s ", f, err.Error()), ctx)
		return
	}

	response.OkWithData(id, ctx)
}

func (clusterApi *ClusterApi) UpdateCluster(ctx *gin.Context) {
	f := "UpdateCluster()-->"

	var cluster db_info.OwlCluster
	if err := ctx.BindJSON(&cluster); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	if err := db_info.UpdateCluster(&cluster); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,update cluster failed :%s ", f, err.Error()), ctx)
		return
	}

	response.Ok(ctx)
}

func (clusterApi *ClusterApi) DelCluster(ctx *gin.Context) {
	f := "DelCluster()-->"

	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, get param failed :%s, id: %s ", f, err.Error(), idStr), ctx)
		return
	}

	if err := db_info.DelCluster(id); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,del cluster failed :%s ", f, err.Error()), ctx)
		return
	}

	response.Ok(ctx)
}
