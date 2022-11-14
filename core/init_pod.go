package core

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sync"
)

// 保存 Pod集合
type PodMapStruct struct {
	data sync.Map  // [key string] []*v1.Pod    key=>namespace
}

// Get 由namespace与podName取到对应的pod
func(podMap *PodMapStruct) Get(ns string,podName string) *corev1.Pod {
	if list, ok := podMap.data.Load(ns); ok {
		for _, pod := range list.([]*corev1.Pod){
			if pod.Name == podName {
				return pod
			}
		}
	}
	return nil
}

// Add 加入PodMapStruct
func(podMap *PodMapStruct) Add(pod *corev1.Pod) {
	if list, ok := podMap.data.Load(pod.Namespace); ok {
		list = append(list.([]*corev1.Pod), pod)
		podMap.data.Store(pod.Namespace, list)
	} else {
		podMap.data.Store(pod.Namespace, []*corev1.Pod{pod})
	}
}

// Update 更新PodMapStruct
func(podMap *PodMapStruct) Update(pod *corev1.Pod) error {
	if list, ok := podMap.data.Load(pod.Namespace); ok {
		for i, rangePod := range list.([]*corev1.Pod) {
			if rangePod.Name == pod.Name {
				list.([]*corev1.Pod)[i] = pod
			}
		}
		return nil
	}
	return fmt.Errorf("Pod-%s not found",pod.Name)
}

// Delete 删除PodMapStruct
func(podMap *PodMapStruct) Delete(pod *corev1.Pod){
	if list, ok := podMap.data.Load(pod.Namespace); ok {
		for i, rangePod := range list.([]*corev1.Pod) {
			if rangePod.Name == pod.Name {
				newList := append(list.([]*corev1.Pod)[:i], list.([]*corev1.Pod)[i+1:]...)
				podMap.data.Store(pod.Namespace, newList)
				break
			}
		}
	}
}


// ListByLabels 根据标签获取 POD列表
func(podMap *PodMapStruct) ListByLabels(ns string,labels []map[string]string) ([]*corev1.Pod,error){
	ret := make([]*corev1.Pod,0)
	if list, ok := podMap.data.Load(ns); ok {
		for _, pod := range list.([]*corev1.Pod){
			for _, label := range labels {
				if reflect.DeepEqual(pod.Labels, label) {  //标签完全匹配
					ret = append(ret, pod)
				}
			}
		}
		return ret,nil
	}
	return nil, fmt.Errorf("pods not found ")
}

// DebugListbyns 根据namespace取出podList
func(podMap *PodMapStruct) DebugListbyns(ns string) ([]*corev1.Pod){
	ret := make([]*corev1.Pod, 0)
	if list, ok := podMap.data.Load(ns); ok {
		for _, pod := range list.([]*corev1.Pod) {
			ret = append(ret, pod)
		}

	}
	return ret
}


var PodMap *PodMapStruct  //作为全局对象

func init() {
	PodMap = &PodMapStruct{}
}

