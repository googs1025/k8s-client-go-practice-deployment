package core

import (
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	"log"
	"sync"
)

// ReplicaSet 集合
type RsMapStruct struct {
	data sync.Map   // [key string] []*appv1.ReplicaSet    key=>namespace
}
func(this *RsMapStruct) Add(rs *appv1.ReplicaSet){
	if list,ok:=this.data.Load(rs.Namespace);ok{
		list=append(list.([]*appv1.ReplicaSet),rs)
		this.data.Store(rs.Namespace,list)
	}else{
		this.data.Store(rs.Namespace,[]*appv1.ReplicaSet{rs})
	}
}
func(this *RsMapStruct) Update(rs *appv1.ReplicaSet) error {
	if list,ok:=this.data.Load(rs.Namespace);ok{
		for i,range_rs:=range list.([]*appv1.ReplicaSet){
			if range_rs.Name==rs.Name{
				list.([]*appv1.ReplicaSet)[i]=rs
			}
		}
		return nil
	}
	return fmt.Errorf("rs-%s not found",rs.Name)
}
func(this *RsMapStruct) Delete(rs *appv1.ReplicaSet){
	if list,ok:=this.data.Load(rs.Namespace);ok{
		for i,range_rs:=range list.([]*appv1.ReplicaSet){
			if range_rs.Name==rs.Name{
				newList:= append(list.([]*appv1.ReplicaSet)[:i], list.([]*appv1.ReplicaSet)[i+1:]...)
				this.data.Store(rs.Namespace,newList)
				break
			}
		}
	}
}
//普普通通的函数， 就是根据ns获取 对应的rs列表
func(this *RsMapStruct) ListByNameSpace(ns string) ([]*appv1.ReplicaSet,error){
	if list,ok:=this.data.Load(ns);ok {
		return list.([]*appv1.ReplicaSet),nil
	}
	return nil,fmt.Errorf("pods not found ")
}
var RsMap *RsMapStruct
type RsHandler struct {}
func(this *RsHandler) OnAdd(obj interface{}){
	RsMap.Add(obj.(*appv1.ReplicaSet))
}
func(this *RsHandler) OnUpdate(oldObj, newObj interface{}){
	err:=RsMap.Update(newObj.(*appv1.ReplicaSet))
	if err!=nil{
		log.Println(err)
	}
}
func(this *RsHandler)	OnDelete(obj interface{}){
	if d,ok:=obj.(*appv1.ReplicaSet);ok{
		RsMap.Delete(d)
	}
}

func init() {
	RsMap=&RsMapStruct{}
}

