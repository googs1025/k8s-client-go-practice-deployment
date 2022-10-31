package core

import (
	"fmt"
	"k8s-client-go-api-practice/initClient"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"sync"
)

/*
	可以先去test/ListWatch_test.go 查看test方法。
 */

type DeploymentMap struct {
	data sync.Map  // [key string] []*v1.Deployment    key=>namespace
}
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
			if rangeDep.Name == dep.Name{
				newList := append(list.([]*v1.Deployment)[:i], list.([]*v1.Deployment)[i+1:]...)
				depMap.data.Store(dep.Namespace,newList)
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
func(depMap *DeploymentMap) GetDeployment(ns string,depName string) (*v1.Deployment,error){
	if list, ok := depMap.data.Load(ns); ok {
		for _, item := range list.([]*v1.Deployment) {
			if item.Name == depName{
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


// InitDeployment 初始化调用监听deployment事件
func InitDeployment(){
	// 创建SharedInformerFactory
	fact:=informers.NewSharedInformerFactory(initClient.K8sClient, 0)
	// 加入所有资源的informer
	depInformer := fact.Apps().V1().Deployments()
	depInformer.Informer().AddEventHandler(&DeploymentHandler{})

	podInformer := fact.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(&PodHandler{})

	rsInformer := fact.Apps().V1().ReplicaSets()
	rsInformer.Informer().AddEventHandler(&RsHandler{})

	eventInformer := fact.Core().V1().Events()
	eventInformer.Informer().AddEventHandler(&EventHandler{})

	fact.Start(wait.NeverStop)

}
