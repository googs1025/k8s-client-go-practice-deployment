package core

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"sync"
)

/*
	可以先去test/ListWatch_test.go 查看test方法。
 */

type DeploymentMap struct {
	data sync.Map  // [key string] []*v1.Deployment    key=>namespace
}

/*
	使用一个map string -> []*v1.Deployment
	来存储每个namespace下的deployment
 */

// Add 添加
func(depMap *DeploymentMap) Add(dep *v1.Deployment) {
	if list, ok := depMap.data.Load(dep.Namespace); ok {
		list = append(list.([]*v1.Deployment), dep)
		depMap.data.Store(dep.Namespace, list)
	} else {
		depMap.data.Store(dep.Namespace, []*v1.Deployment{dep})
	}
}
// Update 更新
func(depMap *DeploymentMap) Update(dep *v1.Deployment) error {
	if list,ok := depMap.data.Load(dep.Namespace); ok {
		for i, rangeDep := range list.([]*v1.Deployment){
			if rangeDep.Name == dep.Name {
				list.([]*v1.Deployment)[i] = dep 	// 替换
			}
		}
		return nil
	}
	return fmt.Errorf("deployment-%s not found",dep.Name)
}
// Delete 删除
func(depMap *DeploymentMap) Delete(dep *v1.Deployment){
	if list, ok := depMap.data.Load(dep.Namespace); ok {
		for i,rangeDep := range list.([]*v1.Deployment) {
			if rangeDep.Name == dep.Name {
				// 切片操作
				newList := append(list.([]*v1.Deployment)[:i], list.([]*v1.Deployment)[i+1:]...)
				depMap.data.Store(dep.Namespace,newList)	// 需要重新存入缓存
				break
			}
		}
	}
}

// ListDeploymentByNamespace 根据namespace取出deploymentList
func(depMap *DeploymentMap) ListDeploymentByNamespace(ns string) ([]*v1.Deployment,error){
	// 从map中拿出对应namespace的 deploymentList
	if list, ok := depMap.data.Load(ns); ok {
		return  list.([]*v1.Deployment), nil
	}
	return nil, fmt.Errorf("record not found")
}

// GetDeployment 由namespace与deploymentName取到deployment
func(depMap *DeploymentMap) GetDeployment(ns string, deploymentName string) (*v1.Deployment,error){
	if list, ok := depMap.data.Load(ns); ok {
		for _, item := range list.([]*v1.Deployment) {
			if item.Name == deploymentName {
				return item, nil
			}
		}
	}
	return nil, fmt.Errorf("record not found")
}

var DepMap *DeploymentMap  //作为全局对象

func init() {
	DepMap = &DeploymentMap{}
}



