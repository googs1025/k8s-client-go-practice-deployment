package core

import (
	"fmt"
	"k8s.io/api/core/v1"
	"sync"
)

// EventSet 集合 用来保存事件, 只保存最新的一条
var EventMap *EventMapStruct

type EventMapStruct struct {
	data sync.Map   // [key string] *v1.Event
	// key=>namespace+"_"+kind+"_"+name 这里的name 不一定是pod ,这样确保唯一
}

// GetMessage 由namespace、kind、name找到特定的event
func(eventMap *EventMapStruct) GetMessage(ns string, kind string, name string) string {
	key := fmt.Sprintf("%s-%s-%s", ns, kind, name)
	if v, ok := eventMap.data.Load(key); ok {
		return v.(*v1.Event).Message
	}

	return ""
}

func init() {
	EventMap = &EventMapStruct{}
}
