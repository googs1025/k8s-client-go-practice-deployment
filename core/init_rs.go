package core

import (
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	"sync"
)

// ReplicaSet 集合
type RsMapStruct struct {
	data sync.Map   // [key string] []*appv1.ReplicaSet    key=>namespace
}
func(rsMap *RsMapStruct) Add(rs *appv1.ReplicaSet){
	if list,ok:=rsMap.data.Load(rs.Namespace);ok{
		list=append(list.([]*appv1.ReplicaSet),rs)
		rsMap.data.Store(rs.Namespace,list)
	}else{
		rsMap.data.Store(rs.Namespace,[]*appv1.ReplicaSet{rs})
	}
}
func(rsMap *RsMapStruct) Update(rs *appv1.ReplicaSet) error {
	if list,ok:=rsMap.data.Load(rs.Namespace);ok{
		for i,range_rs:=range list.([]*appv1.ReplicaSet){
			if range_rs.Name==rs.Name{
				list.([]*appv1.ReplicaSet)[i]=rs
			}
		}
		return nil
	}
	return fmt.Errorf("rs-%s not found",rs.Name)
}
func(rsMap *RsMapStruct) Delete(rs *appv1.ReplicaSet){
	if list,ok:=rsMap.data.Load(rs.Namespace);ok{
		for i,range_rs:=range list.([]*appv1.ReplicaSet){
			if range_rs.Name==rs.Name{
				newList:= append(list.([]*appv1.ReplicaSet)[:i], list.([]*appv1.ReplicaSet)[i+1:]...)
				rsMap.data.Store(rs.Namespace,newList)
				break
			}
		}
	}
}
//普普通通的函数， 就是根据ns获取 对应的rs列表
func(rsMap *RsMapStruct) ListByNameSpace(ns string) ([]*appv1.ReplicaSet,error){
	if list,ok:=rsMap.data.Load(ns);ok {
		return list.([]*appv1.ReplicaSet),nil
	}
	return nil,fmt.Errorf("pods not found ")
}
var RsMap *RsMapStruct

func init() {
	RsMap=&RsMapStruct{}
}

