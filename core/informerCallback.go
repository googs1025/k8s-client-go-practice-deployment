package core

import (
	v1 "k8s.io/api/apps/v1"
	"log"
)

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
