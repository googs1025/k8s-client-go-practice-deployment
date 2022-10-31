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

// 普通用的工具

// 解耦一下
//func GetImages(dep v1.Deployment) string   {
//	images := dep.Spec.Template.Spec.Containers[0].Image
//	if imgLen := len(dep.Spec.Template.Spec.Containers); imgLen>1{
//		images += fmt.Sprintf("+其他%d个镜像",imgLen-1)
//	}
//	return images
//}

// GetImagesByPod 拼凑出images
func GetImagesByPod(containers []core.Container) string {
	images := containers[0].Image
	if imgLen := len(containers); imgLen>1 {
		images += fmt.Sprintf("+其他%d个镜像", imgLen-1)
	}
	return images
}

// GetImages 取得image
func GetImages(dep v1.Deployment) string {
	return GetImagesByPod(dep.Spec.Template.Spec.Containers)
}

// GetLabels 取得标签
func GetLabels(m map[string]string) string {
	labels := ""
	// aa=xxx,bb=xxxx
	for k, v := range m {
		if labels != "" {
			labels += ","
		}
		labels += fmt.Sprintf("%s=%s", k, v)
	}
	return labels
}

// GetRsLabelByDeployment 根据deployment获取ReplicaSet的label
func GetRsLabelByDeployment(deployment *v1.Deployment) string {
	// 取得selector
	selector, err := v12.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		log.Fatal(err)
	}
	// 过滤list
	listOpt := v12.ListOptions{
		LabelSelector: selector.String(),
	}

	ctx := context.Background()
	rs, err := initClient.K8sClient.AppsV1().ReplicaSets(deployment.Namespace).List(ctx, listOpt)
	fmt.Println(deployment.ObjectMeta.Annotations["deployment.kubernetes.io/revision"])	// 需要比对到最新的。

	// 遍历rs
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

// ListDeploymentBySelector 用selector 过滤deploymentList，找出对应的rs
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


//根据dep 获取 所有rs的标签  --- 给listwatch使用
func GetRsLableByDeployment_ListWatch(dep *v1.Deployment,rslist []*v1.ReplicaSet) ([]map[string]string,error){

	ret := make([]map[string]string, 0)
	for _, item := range rslist{
		if IsRsFromDep(dep, *item) {
			s,err := v12.LabelSelectorAsMap(item.Spec.Selector)
			if err != nil{
				return nil,err
			}
			ret = append(ret, s)
		}
	}
	return ret, nil

}

// GetRsLableByDeployment 根据dep 获取 当前rs的标签 ---普通调用API
func GetRsLableByDeployment(dep *v1.Deployment) string{
	selector, _ :=  v12.LabelSelectorAsSelector(dep.Spec.Selector)
	listOpt := v12.ListOptions{
		LabelSelector:selector.String(),
	}
	rs, _ := initClient.K8sClient.AppsV1().ReplicaSets(dep.Namespace).
		List(context.Background(),listOpt)

	for _, item := range rs.Items{
		if IsCurrentRsByDep(dep,item) {
			s, err := v12.LabelSelectorAsSelector(item.Spec.Selector)
			if err != nil {
				return ""
			}
			return s.String()
		}
	}
	return ""
}

// IsRsFromDep 判断 rs 是否属于 某个 dep
func IsRsFromDep(dep *v1.Deployment,set v1.ReplicaSet) bool{
	for _, ref := range set.OwnerReferences {
		if ref.Kind == "Deployment" && ref.Name == dep.Name {
			return true
		}
	}
	return false
}

// IsCurrentRsByDep 判断当前 的rs 是否是最新的
func IsCurrentRsByDep(dep *v1.Deployment,set v1.ReplicaSet) bool{
	if set.ObjectMeta.Annotations["deployment.kubernetes.io/revision"] != dep.ObjectMeta.Annotations["deployment.kubernetes.io/revision"]{
		return false
	}
	return  IsRsFromDep(dep, set)

}





