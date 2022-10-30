package debug

import (
	"github.com/gin-gonic/gin"
	"k8s-client-go-api-practice/core"
)

func RegHandlers(r *gin.Engine){
	r.GET("/debug/pods",ListAllPODS)

}
//列出所有PODS
func ListAllPODS(c *gin.Context){
	ns:=c.DefaultQuery("namespace","default")

	c.JSON(200,gin.H{"message":"Ok","result":	core.PodMap.DEBUG_ListByNS(ns)})
}

