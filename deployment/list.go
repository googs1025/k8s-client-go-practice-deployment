package deployment

import (
	"context"
	"k8s-client-go-api-practice/core"
	"k8s-client-go-api-practice/initClient"
	"k8s-client-go-api-practice/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//ListAll list 所有deployment
// 每次都调用api-server，请求多次会每次都进行调用
func ListAll(namespace string) (ret []*Deployment){
	// 方法一：直接调用api-server的list接口
	ctx := context.Background()
	listOpt := metav1.ListOptions{}
	depList, err := initClient.K8sClient.AppsV1().Deployments(namespace).List(ctx,listOpt)
	util.CheckError(err)
	for _, item := range depList.Items{ //遍历所有deployment
		ret = append(ret,&Deployment{
			Name:item.Name,
			Replicas: [3]int32{item.Status.Replicas, item.Status.AvailableReplicas, item.Status.UnavailableReplicas},
			Images: GetImages(item),
		})
	}
	return
}

// ListAllByWatchList 使用watch-list方法来list deployment
// 只有在初始化才会全量list，后面开始监听deployment，有event才会通知。
func ListAllByWatchList(namespace string) (ret []*Deployment) {

	deploymentList, err := core.DepMap.ListDeploymentByNamespace(namespace)
	util.CheckError(err)
	for _, deployment := range deploymentList {
		ret = append(ret,&Deployment{
			Name: deployment.Name,
			Replicas: [3]int32{
				deployment.Status.Replicas,
				deployment.Status.AvailableReplicas,
				deployment.Status.UnavailableReplicas,
			},
			Images: GetImages(*deployment),
		})
	}
	return

}