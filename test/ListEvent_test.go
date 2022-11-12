package test

import (
	"context"
	"fmt"
	."k8s-client-go-api-practice/initClient"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestListEvent(t *testing.T) {
	ctx := context.Background()
	eventList, err := K8sClient.CoreV1().Events("default").
		List(ctx, v1.ListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	for _ , event := range eventList.Items {
		fmt.Printf("事件名称：%s, 事件类型: %s, 事件原因: %s, 事件消息: %s, 事件对应原因: %s \n", event.Name, event.Type, event.Reason,
			event.Message, event.InvolvedObject)
	}


}
