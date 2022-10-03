package deployment

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
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