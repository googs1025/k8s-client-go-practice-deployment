package deployment

import (
	"github.com/gin-gonic/gin"
	"k8s-client-go-api-practice/initClient"
	"k8s-client-go-api-practice/util"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"context"
)


//
func RegHandlers(r *gin.Engine) {
	// 对副本缩阔容
	r.POST("/update/deployment/scale", incrReplicas)
}


func incrReplicas(c *gin.Context) {
	// request
	req := struct {
		Namespace string `json:"ns" binding:"required,min=1"`
		Deployment string `json:"deployment" binding:"required,min=1"`
		Dec bool `json:"dec"` //是否减少一个
	}{}
	util.CheckError(c.ShouldBindJSON(&req))

	ctx := context.Background()
	// 取到当前scale
	scale, err := initClient.K8sClient.AppsV1().Deployments(req.Namespace).
		GetScale(ctx,req.Deployment,v1.GetOptions{})
	util.CheckError(err)
	if req.Dec{ //dec==true代表是减少副本数
		scale.Spec.Replicas--
	}else{
		scale.Spec.Replicas++
	}

	// 更新数量
	_, err = initClient.K8sClient.AppsV1().Deployments(req.Namespace).
		UpdateScale(ctx,req.Deployment,scale,v1.UpdateOptions{})
	util.CheckError(err)
	util.Sunccess("Ok",c)
}
