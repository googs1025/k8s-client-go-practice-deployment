package deployment

import (
	"context"
	"k8s-client-go-api-practice/initClient"
	"k8s-client-go-api-practice/util"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//
func GetPodsByDep(namespace string, dep *v1.Deployment) []*Pod {
	ctx := context.Background()
	// 由label找到deployment
	//listOpt := metav1.ListOptions{
	//	LabelSelector: GetLabels(dep.Spec.Selector.MatchLabels),
	//}

	listOpt := metav1.ListOptions{
		LabelSelector: GetRsLabelByDeployment(dep),
	}

	list, err := initClient.K8sClient.CoreV1().Pods(namespace).List(ctx, listOpt)
	if err != nil {
		panic(err.Error())
	}
	pods := make([]*Pod, len(list.Items))
	for i, item := range list.Items {
		pods[i] = &Pod{
			Name: item.Name,
			Images: GetImagesByPod(item.Spec.Containers),
			NodeName: item.Spec.NodeName,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:03:04"),
		}
	}

	return pods

}

func GetDeployment(namespace string, name string) *Deployment {
	ctx := context.Background()

	getOpt := metav1.GetOptions{}
	depDetail, err := initClient.K8sClient.AppsV1().Deployments(namespace).Get(ctx, name, getOpt)
	util.CheckError(err)

	return &Deployment{
		Name: depDetail.Name,
		NameSpace: depDetail.Namespace,
		Images: GetImages(*depDetail),
		CreateTime: depDetail.CreationTimestamp.Format("2006-01-02 15:03:04"),
		Pods: GetPodsByDep(namespace, depDetail),
	}
}
