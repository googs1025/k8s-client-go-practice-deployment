package core

import (
	"fmt"
	"k8s-client-go-api-practice/initClient"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"log"
	"sync"
)

type DeploymentMap struct {
	data sync.Map  // [key string] []*v1.Deployment    key=>namespace
}
// Add 添加
func(depMap *DeploymentMap) Add(dep *v1.Deployment) {
	if list,ok := depMap.data.Load(dep.Namespace); ok {
		list=append(list.([]*v1.Deployment), dep)
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
// 删除
func(depMap *DeploymentMap) Delete(dep *v1.Deployment){
	if list, ok := depMap.data.Load(dep.Namespace); ok {
		for i,rangeDep := range list.([]*v1.Deployment) {
			if rangeDep.Name == dep.Name{
				newList := append(list.([]*v1.Deployment)[:i], list.([]*v1.Deployment)[i+1:]...)
				depMap.data.Store(dep.Namespace,newList)
				break
			}
		}
	}
}


func(depMap *DeploymentMap) ListDeploymentByNamespace(ns string) ([]*v1.Deployment,error){
	if list,ok := depMap.data.Load(ns); ok {
		return  list.([]*v1.Deployment), nil
	}
	return nil, fmt.Errorf("record not found")
}

var DepMap *DeploymentMap  //作为全局对象
func init() {
	DepMap = &DeploymentMap{}
}

// 回调函数
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

// InitDeployment 初始化调用监听deployment事件
func InitDeployment(){

	fact:=informers.NewSharedInformerFactory(initClient.K8sClient, 0)

	depInformer:=fact.Apps().V1().Deployments()

	depInformer.Informer().AddEventHandler(&DeploymentHandler{})

	fact.Start(wait.NeverStop)


}
