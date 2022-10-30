package core

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sync"
)

// 保存Pod集合
type PodMapStruct struct {
	data sync.Map  // [key string] []*v1.Pod    key=>namespace
}
func(podMap *PodMapStruct) Get(ns string,podName string) *corev1.Pod{
	if list,ok:= podMap.data.Load(ns);ok{
		for _,pod:=range list.([]*corev1.Pod){
			if pod.Name==podName{
				return pod
			}
		}
	}
	return nil
}
func(podMap *PodMapStruct) Add(pod *corev1.Pod){
	if list,ok:= podMap.data.Load(pod.Namespace);ok{
		list=append(list.([]*corev1.Pod),pod)
		podMap.data.Store(pod.Namespace,list)
	}else{
		podMap.data.Store(pod.Namespace,[]*corev1.Pod{pod})
	}
}
func(podMap *PodMapStruct) Update(pod *corev1.Pod) error {
	if list,ok:=podMap.data.Load(pod.Namespace);ok{
		for i,range_pod:=range list.([]*corev1.Pod){
			if range_pod.Name==pod.Name{
				list.([]*corev1.Pod)[i]=pod
			}
		}
		return nil
	}
	return fmt.Errorf("Pod-%s not found",pod.Name)
}
func(podMap *PodMapStruct) Delete(pod *corev1.Pod){
	if list,ok:= podMap.data.Load(pod.Namespace);ok{
		for i,range_pod:=range list.([]*corev1.Pod){
			if range_pod.Name==pod.Name{
				newList:= append(list.([]*corev1.Pod)[:i], list.([]*corev1.Pod)[i+1:]...)
				podMap.data.Store(pod.Namespace,newList)
				break
			}
		}
	}
}
//func(this *PodMapStruct) ListByLabels(ns string,labels []map[string]string) ([]*corev1.Pod,error){
//	ret:=make([]*corev1.Pod,0)
//	if list,ok:=this.data.Load(ns);ok {
//		for _,pod:=range list.([]*corev1.Pod){
//			for _,label:=range labels{
//				if reflect.DeepEqual(pod.Labels,label){  //标签完全匹配
//					ret=append(ret,pod)
//				}
//			}
//		}
//		return ret,nil
//	}
//	return nil,fmt.Errorf("pods not found ")
//}
//根据标签获取 POD列表
func(podMap *PodMapStruct) ListByLabels(ns string,labels []map[string]string) ([]*corev1.Pod,error){
	ret:=make([]*corev1.Pod,0)
	if list,ok:= podMap.data.Load(ns);ok {
		for _,pod:=range list.([]*corev1.Pod){
			for _,label:=range labels{
				if reflect.DeepEqual(pod.Labels,label){  //标签完全匹配
					ret=append(ret,pod)
				}
			}
		}
		return ret,nil
	}
	return nil,fmt.Errorf("pods not found ")
}
func(podMap *PodMapStruct) DEBUG_ListByNS(ns string) ([]*corev1.Pod){
	ret:=make([]*corev1.Pod,0)
	if list,ok:= podMap.data.Load(ns);ok {
		for _,pod:=range list.([]*corev1.Pod){
			ret=append(ret,pod)
		}

	}
	return ret
}


var PodMap *PodMapStruct  //作为全局对象

func init() {
	PodMap = &PodMapStruct{}
}

