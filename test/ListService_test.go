package test

import (
	"context"
	"fmt"
	."k8s-client-go-api-practice/initClient"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"testing"
)

func TestListService(t *testing.T) {

	// 查看版本。
	fmt.Println(K8sClient.ServerVersion())

	fmt.Println("list default下的svc")
	ctx := context.Background()
	list, err := K8sClient.CoreV1().Services("default").List(ctx, v1.ListOptions{})
	if err != nil {
		log.Println(err)
	}
	for _, item := range list.Items {
		fmt.Println(item.Name)
	}


}
