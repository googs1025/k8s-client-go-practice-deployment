package util

import "reflect"

// 用于前端渲染用

func DataBuilder() *CommonDataStruct{
	return NewCommonDataStruct()
}

type CommonDataStruct struct {
	Title string
	Data map[string]interface{}
}


func NewCommonDataStruct() *CommonDataStruct {
	return &CommonDataStruct{Data: make(map[string]interface{})}
}


func(c *CommonDataStruct) SetTitle(title string) *CommonDataStruct{
	c.Title = title
	return c
}


func(c *CommonDataStruct) SetData(key string,value interface{}) *CommonDataStruct{
	c.Data[key] = value
	return c
}


func(c *CommonDataStruct) ToMap() (m map[string]interface{})  {
	m = make(map[string]interface{})
	elem := reflect.ValueOf(c).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		m[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	return
}
