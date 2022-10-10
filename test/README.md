### Informer机制
```bigquery
1. 使用client-go创建SharedInformer
2. 加入需要监听的资源handler ex:DeploymentHandler
3. 注意：资源Handler 必须实现 OnAdd OnUpdate OnDelete 方法
4. 启动 fact.Start(wait.NeverStop)
5. 如果有超时控制，可以使用
```

![]()