# k8s-client-go-practice-deployment
![](https://github.com/googs1025/k8s-client-go-practice-deployment/blob/main/image/client-go-kubernetes.png?ram=true)
### 项目思路逻辑
结合**golang**+**k8s**+**client-go**实现restful-http-api 简易管理系统
### 项目目录
```bigquery
.
├── README.md
├── deployment  // deployment 增刪改查的主逻辑
│   ├── DetailUtil.go   // get deployment
│   ├── commonUtil.go   // 通用
│   ├── handler.go      // post请求: 增减副本数
│   ├── list.go         // list操作
│   └── model.go        // 结构
├── example.go  // 暂时无用，请求反向代理用的脚本
├── go.mod
├── go.sum
├── html    // 简易前端页面
│   ├── common
│   │   ├── footer.html
│   │   └── header.html
│   └── deployment
│       ├── deployment_detail.html
│       ├── deployment_list.html
│       └── list.html
├── initClient  // k8s客户端
│   └── clientSet.go
├── main.go
├── myproxy // 反向代理，功能暂时不能用，有k8s token问题
│   ├── Readme.md
│   ├── myproxy
│   └── myproxy.go
├── proxy_request_try.go
├── static  // 前端
│   ├── css
│   │   ├── bulma.min.css
│   │   └── common.css
│   └── jq3.js
└── util    // gin通用函数
    ├── CommonData.go
    └── error_handle.go
```

### 项目依赖
```bigquery
目前暂无，预计提供mysql建立表存储。
```
### 项目启动
1. 启动
```bigquery
cd 进入项目根目录
[root@vm-0-12-centos k8s-api-practice]# go build
[root@vm-0-12-centos k8s-api-practice]# ./k8s-client-go-api-practice
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /update/deployment/scale  --> k8s-client-go-api-practice/deployment.incrReplicas (2 handlers)
                        [GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (2 handlers)
[GIN-debug] Loaded HTML Templates (8):
	- header
	- deployment_detail.html
	- deployment_list.html
	- list.html
	-
	- footer.html
	- footer
	- header.html

[GIN-debug] GET    /deployments              --> main.main.func2 (2 handlers)
[GIN-debug] GET    /deployments/:name        --> main.main.func3 (2 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.

```
2. 浏览器+使用postman调用
### 预计提供方法
```bigquery
1. list deployment:获取deployment列表 (完成)
2. get deployment:获取特定deployment (完成)
3. delete deployment:删除特定deployment
4. 增减deployment副本 (完成)
5. update deployment:自定义更新deployment
6. create deployment:创建deployment
```

### list-watch
```
list:http短链接调用资源的Api，获取列表 => 或有大量请求的缺点
watch:http长连接持续监听资源，有变化则返回一个watchEvent
使用client-go中的informer可以对list-watch机制进行封装。
解释：
1.刚开始初始化，调用list api获取全量list，并缓存起来。
2.调用watch api去watch资源，发生变更后会通过机制维护缓存(不需要每次都请求api server)
```