package deployment

import (
	"context"
	"k8s-client-go-api-practice/initClient"
	"k8s-client-go-api-practice/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//显示 所有deployment
func ListAll(namespace string ) (ret []*Deployment){
	ctx := context.Background()
	listopt := metav1.ListOptions{}
	depList, err := initClient.K8sClient.AppsV1().Deployments(namespace).List(ctx,listopt)
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