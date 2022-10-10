package deployment

import (
	"context"
	"fmt"
	"k8s-client-go-api-practice/initClient"
	v1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

// 解耦一下
//func GetImages(dep v1.Deployment) string   {
//	images := dep.Spec.Template.Spec.Containers[0].Image
//	if imgLen := len(dep.Spec.Template.Spec.Containers); imgLen>1{
//		images += fmt.Sprintf("+其他%d个镜像",imgLen-1)
//	}
//	return images
//}



func GetLabels(m map[string]string) string {
	labels := ""
	for k, v := range m {
		if labels != "" {
			labels += ","
		}
		labels += fmt.Sprintf("%s=%s", k, v)
	}
	return labels
}

func GetImagesByPod(containers []core.Container) string{
	images:=containers[0].Image
	if imgLen:=len(containers);imgLen>1{
		images+=fmt.Sprintf("+其他%d个镜像",imgLen-1)
	}
	return images
}

func GetImages(dep v1.Deployment) string {
	return GetImagesByPod(dep.Spec.Template.Spec.Containers)
}

// GetRsLabelByDeployment 根据deployment获取ReplicaSet的label
func GetRsLabelByDeployment(deployment *v1.Deployment) string {
	selector, err := v12.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		log.Fatal(err)
	}

	listOpt := v12.ListOptions{
		LabelSelector: selector.String(),
	}
	ctx := context.Background()
	rs, err := initClient.K8sClient.AppsV1().ReplicaSets(deployment.Namespace).List(ctx, listOpt)
	fmt.Println(deployment.ObjectMeta.Annotations["deployment.kubernetes.io/revision"])	// 需要比对到最新的。

	for _, item := range rs.Items {
		// 打印name
		log.Println(item.Name)
		// 打印上层引用对象
		log.Println(item.OwnerReferences)




		if IsCurrentRs(deployment, item) {
			if err != nil {
				log.Println(err)
				return ""
			}
			// 打印selector标签
			s, _ := v12.LabelSelectorAsSelector(item.Spec.Selector)
			log.Println(s.String())
			return s.String()
		}

	}
	return ""


}

// IsCurrentRs 根据deployment 获取 正确的replicaset，因为有版本 标签的不同，需要进行过滤
func IsCurrentRs(deployment *v1.Deployment, replicaset v1.ReplicaSet) bool {
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

func ListDeploymentBySelector(namespace string, deploymentName string) {
	ctx := context.Background()
	deployment, err := initClient.K8sClient.AppsV1().Deployments(namespace).Get(ctx, deploymentName, v12.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	selector, err := v12.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		log.Fatal(err)
	}

	listOpt := v12.ListOptions{
		LabelSelector: selector.String(),
	}
	rs, err := initClient.K8sClient.AppsV1().ReplicaSets(namespace).List(ctx, listOpt)
	fmt.Println(deployment.ObjectMeta.Annotations["deployment.kubernetes.io/revision"])	// 需要比对到最新的。

	for _, item := range rs.Items {
		// 打印name
		log.Println(item.Name)
		// 打印上层引用对象
		log.Println(item.OwnerReferences)
		fmt.Println(IsCurrentRs(deployment, item))

		// 打印selector标签
		s, err := v12.LabelSelectorAsSelector(item.Spec.Selector)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(s.String())
	}

}




