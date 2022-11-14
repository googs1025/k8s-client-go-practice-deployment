package core

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
)

/*
	回调函数：统一操作：
	调用add update delete方法
	同步到map中
 */

// Deployment回调函数
type DeploymentHandler struct {
}

// add回调 加入map中
func(dh *DeploymentHandler) OnAdd(obj interface{}){
	DepMap.Add(obj.(*v1.Deployment))
}

func(dh *DeploymentHandler) OnUpdate(oldObj, newObj interface{}){
	err := DepMap.Update(newObj.(*v1.Deployment))
	if err != nil {
		log.Println(err)
	}
}

// delete回调 删除map
func(dh *DeploymentHandler)	OnDelete(obj interface{}){
	if d,ok := obj.(*v1.Deployment); ok {
		DepMap.Delete(d)
	}
}

// Pod回调函数
type PodHandler struct {

}

func(ph *PodHandler) OnAdd(obj interface{}){
	PodMap.Add(obj.(*corev1.Pod))
}

func(ph *PodHandler) OnUpdate(oldObj, newObj interface{}){
	err := PodMap.Update(newObj.(*corev1.Pod))
	if err != nil {
		log.Println(err)
	}
}

func(ph *PodHandler) OnDelete(obj interface{}){
	if d, ok := obj.(*corev1.Pod); ok {
		PodMap.Delete(d)
	}
}

// ReplicaSet回调函数
type RsHandler struct {}
func(rh *RsHandler) OnAdd(obj interface{}){
	RsMap.Add(obj.(*v1.ReplicaSet))
}
func(rh *RsHandler) OnUpdate(oldObj, newObj interface{}){
	err := RsMap.Update(newObj.(*v1.ReplicaSet))
	if err != nil {
		log.Println(err)
	}
}

func(rh *RsHandler)	OnDelete(obj interface{}){
	if d, ok := obj.(*v1.ReplicaSet); ok {
		RsMap.Delete(d)
	}
}

// Event回调函数
type EventHandler struct {}

func(eh *EventHandler) storeData(obj interface{}, isDelete bool){
	if event, ok := obj.(*corev1.Event); ok {
		// InvolvedObject 可以找到事件对应的源对象
		key := fmt.Sprintf("%s-%s-%s", event.Namespace, event.InvolvedObject.Kind, event.InvolvedObject.Name)
		if !isDelete {
			EventMap.data.Store(key,event)
		} else {
			EventMap.data.Delete(key)
		}
	}
}

func(eh *EventHandler) OnAdd(obj interface{}){
	eh.storeData(obj,false)
}

func(eh *EventHandler) OnUpdate(oldObj, newObj interface{}){
	eh.storeData(newObj,false)
}

func(eh *EventHandler)	OnDelete(obj interface{}){
	eh.storeData(obj,true)
}
