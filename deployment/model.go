package deployment

type Deployment struct {
	Name string
	NameSpace string
	Replicas [3]int32
	Images string
	CreateTime string
	Pods []*Pod
}

type Pod struct {
	Name string
	NameSpace string  //新增一个命名空间
	Images string
	NodeName string
	IP []string   // 第一个是 POD IP 第二个是 node ip
	Phase string  // pod 当前所处的阶段
	IsReady bool  //判断pod 是否就绪
	Message string
	CreateTime string
}
