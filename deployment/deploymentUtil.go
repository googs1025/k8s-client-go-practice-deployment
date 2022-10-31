package deployment

import (
	"context"
	"k8s-client-go-api-practice/initClient"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

)

type DeploymentRequest struct {
	Name string `form:"name" binding:"required,min=2"`
	Image string `form:"image" binding:"required,min=5"`
}

// genContainers 生成容器配置
func genContainers(req *DeploymentRequest) []corev1.Container{
	ret := make([]corev1.Container, 1)
	ret[0] = corev1.Container{
		Name: req.Name,
		Image: req.Image,
	}
	return ret
}
//生成标签配置
func genLabels(req *DeploymentRequest) map[string]string {
	return map[string]string{
		"app": req.Name,
	}
}

func CreateDeployment(req *DeploymentRequest) error {
	ns := "default"
	_, err := initClient.K8sClient.AppsV1().Deployments(ns).
		Create(context.Background(), &v1.Deployment {
			ObjectMeta: metav1.ObjectMeta{Name: req.Name,Namespace: ns},
			Spec: v1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: genLabels(req),
				},
				Template:corev1.PodTemplateSpec{
					ObjectMeta:metav1.ObjectMeta{
						Labels:genLabels(req),
					},
					Spec:corev1.PodSpec{
						Containers:genContainers(req),
					},
				},
			},

		}, metav1.CreateOptions{})
	return err
}

