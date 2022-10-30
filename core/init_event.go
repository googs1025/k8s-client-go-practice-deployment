package core

import (
	"fmt"
	"k8s.io/api/core/v1"
	"sync"
)

// EventSet 集合 用来保存事件, 只保存最新的一条
var EventMap *EventMapStruct
type EventMapStruct struct {
	data sync.Map   // [key string] *v1.Event
	// key=>namespace+"_"+kind+"_"+name 这里的name 不一定是pod ,这样确保唯一
}
func(this *EventMapStruct) GetMessage(ns string,kind string,name string) string{
	key:=fmt.Sprintf("%s_%s_%s",ns,kind,name)
	if v,ok:=this.data.Load(key);ok{
		return v.(*v1.Event).Message
	}

	return ""
}
type EventHandler struct {}
func(this *EventHandler) storeData(obj interface{},isdelete bool){
	if event,ok:=obj.(*v1.Event);ok{
		key:=fmt.Sprintf("%s_%s_%s",event.Namespace,event.InvolvedObject.Kind,event.InvolvedObject.Name)
		if !isdelete{
			EventMap.data.Store(key,event)
		}else{
			EventMap.data.Delete(key)
		}
	}
}
func(this *EventHandler) OnAdd(obj interface{}){
	this.storeData(obj,false)
}
func(this *EventHandler) OnUpdate(oldObj, newObj interface{}){
	this.storeData(newObj,false)
}
func(this *EventHandler)	OnDelete(obj interface{}){
	this.storeData(obj,true)
}
func init() {
	EventMap=&EventMapStruct{}
}
