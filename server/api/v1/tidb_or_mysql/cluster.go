package tidb_or_mysql

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/db_info"
)

type ClusterApi struct {}

func (clusterApi *ClusterApi)ListDB(ctx *gin.Context) {
	f := "ListDB()-->"

	dbInfo, err := db_info.ListAllDB()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,list db failed :%s ", f, err.Error()), ctx)
		return
	}

	response.OkWithData(dbInfo, ctx)
}

func (clusterApi *ClusterApi)ListCluster(ctx *gin.Context) {
	f := "ListCluster()-->"

	clusters, err := db_info.ListClusterForUI()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,list cluster failed :%s ", f, err.Error()), ctx)
		return
	}

	response.OkWithData(ListData{
		Items:  clusters,
		Total:  int64(len(clusters)),
		More:   false,
		Offset: 0,
	}, ctx)
}

func (clusterApi *ClusterApi)AddCluster(ctx *gin.Context) {
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

func (clusterApi *ClusterApi)UpdateCluster(ctx *gin.Context) {
	f := "UpdateCluster()-->"

	var cluster db_info.OwlCluster
	if err := ctx.BindJSON(&cluster); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()),  ctx)
		return
	}

	if err := db_info.UpdateCluster(&cluster); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,update cluster failed :%s ", f, err.Error()), ctx)
		return
	}

	response.Ok(ctx)
}

func (clusterApi *ClusterApi)DelCluster(ctx *gin.Context) {
	f := "DelCluster()-->"

	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, get param failed :%s, id: %s ", f, err.Error(), idStr),ctx)
		return
	}

	if err := db_info.DelCluster(id); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,del cluster failed :%s ", f, err.Error()), ctx)
		return
	}

	response.Ok(ctx)
}
