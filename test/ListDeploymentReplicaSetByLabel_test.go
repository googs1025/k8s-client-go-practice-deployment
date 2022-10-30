package test

import (
	"context"
	"fmt"
	"k8s-client-go-api-practice/initClient"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	v12 "k8s.io/api/apps/v1"
	"testing"
)

func TestListDeploymentBySelector(t *testing.T) {
	res := ListDeploymentBySelector("default", "webapp")
	fmt.Println(res)
}

func ListDeploymentBySelector(namespace string, deploymentName string) string {

	ctx := context.Background()
	// 先找到特定label
	deployment, err := initClient.K8sClient.AppsV1().Deployments(namespace).
		Get(ctx, deploymentName, v1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	selector, err := v1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		log.Fatal(err)
	}

	listOpt := v1.ListOptions{
		LabelSelector: selector.String(),
	}
	// 再用ReplicaSets 查一变
	rs, err := initClient.K8sClient.AppsV1().ReplicaSets(namespace).List(ctx, listOpt)
	fmt.Println("deployment.kubernetes.io/revision:",deployment.ObjectMeta.Annotations["deployment.kubernetes.io/revision"])	// 需要比对到最新的。

	for _, item := range rs.Items {
		// 打印name
		log.Println("rs name: ", item.Name)
		// 打印上层引用对象(数组)
		log.Println("OwnerReferences:", item.OwnerReferences)
		fmt.Println(IsCurrentRs(deployment, item))

		// 打印selector标签
		s, err := v1.LabelSelectorAsSelector(item.Spec.Selector)
		if err != nil {
			log.Fatal(err)
			return ""
		}
		log.Println(s.String())
	}
	return ""

}

// IsCurrentRs 根据deployment 获取 正确的replicaset，因为有版本 标签的不同，需要进行过滤
func IsCurrentRs(deployment *v12.Deployment, replicaset v12.ReplicaSet) bool {
	if replicaset.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] != deployment.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] {
		return false
	}
	for _, ref := range replicaset.OwnerReferences {
		if ref.Kind == "Deployment" && ref.Name == deployment.Name {
			return true
		}
	}
	return false
}