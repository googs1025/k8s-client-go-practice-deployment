package test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-client-go-api-practice/initClient"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"log"
	"sync"
	"testing"
	"time"
)

func init() {
	DepMap = &DeploymentMap{}
}

var DepMap *DeploymentMap
type DeploymentMap struct {
	data sync.Map	// key: namespace value: []*v1.Deployment
}

func (depMap *DeploymentMap) Add(deployment *v1.Deployment) {

	if value, ok := depMap.data.Load(deployment.Name); ok {
		v := value.([]*v1.Deployment)
		v = append(v, deployment)
		depMap.data.Store(deployment.Name, v)
	} else {
		depMap.data.Store(deployment.Name, []*v1.Deployment{deployment})
	}
}


// DepHandler 如果要在informer使用，需要实现OnAdd,OnUpdate,OnDelete
type DeploymentHandler struct {

}


func(d *DeploymentHandler) OnAdd(obj interface{}){
	//if deployment, ok := obj.(*v1.Deployment); ok {
	//	fmt.Println("新增的deployment:",deployment.Name)
	//}
	deploymentObj := obj.(*v1.Deployment)
	DepMap.Add(deploymentObj)

}

func(d *DeploymentHandler) OnUpdate(oldObj, newObj interface{}){
	// 拿到新对象数据。
	if deployment,ok := newObj.(*v1.Deployment); ok {
		fmt.Println("deployment名称", deployment.Name)
		fmt.Println("更新后deployment数量", *deployment.Spec.Replicas)
	}
}

func(d *DeploymentHandler) OnDelete(obj interface{}){
	if deployment, ok := obj.(*v1.Deployment); ok {
		fmt.Println("删除的deployment:", deployment.Name)
	}
}


func TestInformerListWatch(t *testing.T) {

	//UseInformerWatchDeployment()

	// 常用方式：利用工厂建立sharedInformer，可以共享informer监听的资源
	fact := informers.NewSharedInformerFactory(initClient.K8sClient, 0)
	deploymentInformer := fact.Apps().V1().Deployments()
	// 可以增加多个资源的handler
	deploymentInformer.Informer().AddEventHandler(&DeploymentHandler{})

	fact.Start(wait.NeverStop)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		log.Fatal("请求超时！")
	default:
		r := gin.New()
		defer func() {
			r.Run(":8081")
		}()

		r.GET("/", func(c *gin.Context) {
			res := make([]string,0)
			DepMap.data.Range(func(key, value any) bool {
				if key == "default" {
					for _, deployment := range value.([]*v1.Deployment) {
						res = append(res, deployment.Name)
					}
				}
				return true
			})
			c.JSON(200, res)


		})

	}


}

func UseInformerWatchDeployment() {
	// 一种方式：建立一个watch的informer
	s, c := cache.NewInformer(
		cache.NewListWatchFromClient(initClient.K8sClient.AppsV1().RESTClient(),
		"deployments", "default", fields.Everything()),
		&v1.Deployment{},
		0,
		&DeploymentHandler{},
		)

	c.Run(wait.NeverStop)
	s.List()
}