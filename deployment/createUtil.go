package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s-client-go-api-practice/initClient"
	"k8s-client-go-api-practice/util"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"log"
)

func Create(){

	ctx := context.Background()

	listOpt := metav1.ListOptions{}
	depList,err := initClient.K8sClient.AppsV1().Deployments("default").List(ctx, listOpt)

	for _, item := range depList.Items{ //遍历所有deployment
		fmt.Println(item.Name)
	}
	ngxDep := &v1.Deployment{} //我们要创建的deployment
	b,_ := ioutil.ReadFile("yamls/nginx.yaml")
	ngxJson, _ := yaml.ToJSON(b)
	util.CheckError(json.Unmarshal(ngxJson, ngxDep))

	createopt:=metav1.CreateOptions{}

	_,err = initClient.K8sClient.AppsV1().Deployments("default").
		Create(ctx,ngxDep,createopt)

	if err != nil {
		log.Fatal(err)
	}
}
