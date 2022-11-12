package deployment

import (
	"context"
	"github.com/gin-gonic/gin"
	"k8s-client-go-api-practice/core"
	"k8s-client-go-api-practice/initClient"
	"k8s-client-go-api-practice/util"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


// RegHandlers 请求的post路由
func RegHandlers(r *gin.Engine) {
	// 对副本缩阔容
	r.POST("/update/deployment/scale", IncrReplicas)
	r.POST("/core/deployments", ListAllDeployments)
	r.POST("/core/pods",ListPodsByDeployment)
	r.GET("/core/pods_json",GetPODJSON)
	r.DELETE("/core/pods",DeletePOD)
	r.POST("/create/deployments",CreateDeploymentHandler)
}

func CreateDeploymentHandler(c *gin.Context){
	depReq := &DeploymentRequest{}

	util.CheckError(c.ShouldBind(depReq))
	util.CheckError(CreateDeployment(depReq))

	c.Redirect(301,"/deployments")
}

// ListAllDeployment list 传入namespace 结果
func ListAllDeployments(c *gin.Context) {
	ns := c.DefaultQuery("namespace", "default")
	c.JSON(200, gin.H{"message":"Ok", "result": ListAllByWatchList(ns)})
}

// ListPodsByDeployment 根据 deployment名称 获取 pod列表
func ListPodsByDeployment(c *gin.Context){
	// 取得请求参数
	ns := c.DefaultQuery("namespace","default")
	depName := c.DefaultQuery("deployment","default")

	dep, err := core.DepMap.GetDeployment(ns, depName)
	util.CheckError(err)

	rsList,err:=core.RsMap.ListByNameSpace(ns)  // 根据namespace 取到 所有rs
	util.CheckError(err)

	labels, err:=GetRsLableByDeployment_ListWatch(dep, rsList) //根据deployment过滤出 rs，然后直接获取标签
	util.CheckError(err)

	c.JSON(200,gin.H{"message":"Ok","result": ListPodsByLabel(ns,labels)})
}

// incrReplicas 扩缩容副本数
func IncrReplicas(c *gin.Context) {
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
	if req.Dec { //dec==true代表是减少副本数
		scale.Spec.Replicas--
	} else {
		scale.Spec.Replicas++
	}

	// 更新数量
	_, err = initClient.K8sClient.AppsV1().Deployments(req.Namespace).
		UpdateScale(ctx,req.Deployment,scale,v1.UpdateOptions{})
	util.CheckError(err)
	util.Sunccess("Ok",c)
}

//删除POD
func DeletePOD(c *gin.Context){
	ns := c.DefaultQuery("namespace","default")
	podName := c.DefaultQuery("pod","")
	if podName == "" || ns == "" {
		panic("error ns or pod")
	}
	util.CheckError(DeletePod(ns, podName))
	c.JSON(200,gin.H{"message":"Ok"})

}


//获取POD的JSON详细内容
func GetPODJSON(c *gin.Context){
	// 传入参数
	ns := c.DefaultQuery("namespace","default")
	podName := c.DefaultQuery("pod","")
	if podName == "" || ns == "" {
		panic("error ns or pod")
	}
	// 用list-watch机制，存到内存里的数据
	if pod := core.PodMap.Get(ns,podName); pod == nil {
		panic("no such pod " + podName)
	} else {
 		c.JSON(200,pod)
	}

}
