package deployment

import (
	"k8s-client-go-api-practice/core"
	"k8s-client-go-api-practice/util"
)

//ListAll list 所有deployment
// 每次都调用api-server，请求多次会每次都进行调用
//func ListAll(namespace string) (ret []*Deployment){
//	// 方法一：直接调用api-server的list接口
//	ctx := context.Background()
//	listOpt := metav1.ListOptions{}
//	depList, err := initClient.K8sClient.AppsV1().Deployments(namespace).List(ctx,listOpt)
//	util.CheckError(err)
//	for _, item := range depList.Items{ //遍历所有deployment
//		ret = append(ret,&Deployment{
//			Name:item.Name,
//			Replicas: [3]int32{item.Status.Replicas, item.Status.AvailableReplicas, item.Status.UnavailableReplicas},
//			Images: GetImages(item),
//		})
//	}
//	return
//}

// ListAllByWatchList 使用watch-list方法来list deployment
// 只有在初始化才会全量list，后面开始监听deployment，有event才会通知。
// ListAll ListAllByWatchList
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
			IsComplete: GetDeploymentIsComplete(deployment),
			//Message: core.EventMap.GetMessage(deployment.Namespace, "Deployment", deployment.Name),
			Message: GetDeploymentCondition(deployment),
		})
	}
	return

}

//这里的函数好比 DTO . 把原生的 deployment 或pod转换为  自己的 实体对象
func ListPodsByLabel(ns string, labels []map[string]string) (ret []*Pod){
	podList,err :=core.PodMap.ListByLabels(ns,labels)
	util.CheckError(err)
	for _, pod := range podList {
		ret = append(ret, &Pod{
			Name: pod.Name,
			NameSpace: pod.Namespace,
			Images: GetImagesByPod(pod.Spec.Containers),
			NodeName: pod.Spec.NodeName,
			Phase: string(pod.Status.Phase),// 阶段
			IsReady: GetPodIsReady(*pod), //是否就绪
			IP: []string{pod.Status.PodIP,pod.Status.HostIP},
			Message: core.EventMap.GetMessage(pod.Namespace,"Pod", pod.Name),
			CreateTime: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return
}