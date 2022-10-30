package deployment

import (
	"k8s-client-go-api-practice/initClient"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"context"
)

//判断POD是否就绪
func GetPodIsReady(pod v1.Pod)  bool {
	for _,condition:=range pod.Status.Conditions{
		if condition.Type=="ContainersReady" && condition.Status!="True"{
			return false
		}
	}
	for _,rg:=range pod.Spec.ReadinessGates{
		for _,condition:=range pod.Status.Conditions{
			if condition.Type==rg.ConditionType && condition.Status!="True"{
				return false
			}
		}
	}
	return true
}

func DeletePod(ns string,podName string) error {
	return initClient.K8sClient.CoreV1().Pods(ns).
		Delete(context.Background(),podName,v12.DeleteOptions{})
}
