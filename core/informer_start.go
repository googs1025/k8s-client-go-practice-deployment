package core

import (
	"k8s-client-go-api-practice/initClient"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
)

// InitResourceInformer 初始化调用监听所有资源事件
func InitResourceInformer(){
	// 创建SharedInformerFactory
	fact := informers.NewSharedInformerFactory(initClient.K8sClient, 0)
	// 加入所有资源的informer

	// 加入deployment 资源informer
	depInformer := fact.Apps().V1().Deployments()
	depInformer.Informer().AddEventHandler(&DeploymentHandler{})

	// 加入pod 资源informer
	podInformer := fact.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(&PodHandler{})

	// 加入ReplicaSets 资源informer
	rsInformer := fact.Apps().V1().ReplicaSets()
	rsInformer.Informer().AddEventHandler(&RsHandler{})

	// 加入event 资源informer
	eventInformer := fact.Core().V1().Events()
	eventInformer.Informer().AddEventHandler(&EventHandler{})

	// 启动informer
	fact.Start(wait.NeverStop)

}
