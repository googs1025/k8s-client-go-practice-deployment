package core

import (
	"fmt"
	appv1 "k8s.io/api/apps/v1"
	"sync"
)

// ReplicaSet 集合
type ReplicaSetMap struct {
	data sync.Map   // [key string] []*appv1.ReplicaSet    key=>namespace
}

// Add 加入RsMap
func(rsMap *ReplicaSetMap) Add(rs *appv1.ReplicaSet){
	if list, ok := rsMap.data.Load(rs.Namespace); ok {
		list = append(list.([]*appv1.ReplicaSet), rs)
		rsMap.data.Store(rs.Namespace, list)
	} else {
		rsMap.data.Store(rs.Namespace, []*appv1.ReplicaSet{rs})
	}
}

// Update 更新RsMap
func(rsMap *ReplicaSetMap) Update(rs *appv1.ReplicaSet) error {
	if list, ok := rsMap.data.Load(rs.Namespace); ok {
		for i, rangeRs := range list.([]*appv1.ReplicaSet){
			if rangeRs.Name == rs.Name {
				list.([]*appv1.ReplicaSet)[i] = rs
			}
		}
		return nil
	}
	return fmt.Errorf("rs-%s not found",rs.Name)
}

// Delete 删除RsMap
func(rsMap *ReplicaSetMap) Delete(rs *appv1.ReplicaSet){
	if list, ok := rsMap.data.Load(rs.Namespace); ok {
		for i, rangeRs := range list.([]*appv1.ReplicaSet) {
			if rangeRs.Name == rs.Name {
				newList := append(list.([]*appv1.ReplicaSet)[:i], list.([]*appv1.ReplicaSet)[i+1:]...)
				rsMap.data.Store(rs.Namespace,newList)
				break
			}
		}
	}
}

// ListByNameSpace 根据namespace获取对应的ReplicaSet
func(rsMap *ReplicaSetMap) ListByNameSpace(ns string) ([]*appv1.ReplicaSet,error){
	if list, ok := rsMap.data.Load(ns); ok {
		return list.([]*appv1.ReplicaSet),nil
	}
	return nil,fmt.Errorf("pods not found ")
}

var RsMap *ReplicaSetMap

func init() {
	RsMap = &ReplicaSetMap{}
}

