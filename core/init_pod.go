package core

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"log"
	"reflect"
	"sync"
)

// 保存Pod集合
type PodMapStruct struct {
	data sync.Map  // [key string] []*v1.Pod    key=>namespace
}
func(this *PodMapStruct) Get(ns string,podName string) *corev1.Pod{
	if list,ok:=this.data.Load(ns);ok{
		for _,pod:=range list.([]*corev1.Pod){
			if pod.Name==podName{
				return pod
			}
		}
	}
	return nil
}
func(this *PodMapStruct) Add(pod *corev1.Pod){
	if list,ok:=this.data.Load(pod.Namespace);ok{
		list=append(list.([]*corev1.Pod),pod)
		this.data.Store(pod.Namespace,list)
	}else{
		this.data.Store(pod.Namespace,[]*corev1.Pod{pod})
	}
}
func(this *PodMapStruct) Update(pod *corev1.Pod) error {
	if list,ok:=this.data.Load(pod.Namespace);ok{
		for i,range_pod:=range list.([]*corev1.Pod){
			if range_pod.Name==pod.Name{
				list.([]*corev1.Pod)[i]=pod
			}
		}
		return nil
	}
	return fmt.Errorf("Pod-%s not found",pod.Name)
}
func(this *PodMapStruct) Delete(pod *corev1.Pod){
	if list,ok:=this.data.Load(pod.Namespace);ok{
		for i,range_pod:=range list.([]*corev1.Pod){
			if range_pod.Name==pod.Name{
				newList:= append(list.([]*corev1.Pod)[:i], list.([]*corev1.Pod)[i+1:]...)
				this.data.Store(pod.Namespace,newList)
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
func(this *PodMapStruct) ListByLabels(ns string,labels []map[string]string) ([]*corev1.Pod,error){
	ret:=make([]*corev1.Pod,0)
	if list,ok:=this.data.Load(ns);ok {
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
func(this *PodMapStruct) DEBUG_ListByNS(ns string) ([]*corev1.Pod){
	ret:=make([]*corev1.Pod,0)
	if list,ok:=this.data.Load(ns);ok {
		for _,pod:=range list.([]*corev1.Pod){
			ret=append(ret,pod)
		}

	}
	return ret
}
type PodHandler struct {}
func(this *PodHandler) OnAdd(obj interface{}){
	PodMap.Add(obj.(*corev1.Pod))
}
func(this *PodHandler) OnUpdate(oldObj, newObj interface{}){
	err:=PodMap.Update(newObj.(*corev1.Pod))
	if err!=nil{
		log.Println(err)
	}
}
func(this *PodHandler)	OnDelete(obj interface{}){
	if d,ok:=obj.(*corev1.Pod);ok{
		PodMap.Delete(d)
	}
}


var PodMap *PodMapStruct  //作为全局对象

func init() {
	PodMap=&PodMapStruct{}
}

