# Worker model
## Overview
该服务应该由Factory创建
该模块是最小的OJ单元，每个单元关联一个Container  
该单元接收提交的题目，返回输出+运行信息  
## Flowchart

## Functional Decomposition
* 容器操作
    * 构建容器
    * 执行代码
    * 外部接口
## Package
* docker 实现worker内 container的所有操作
``` go
type Container interface {
    Build() // 通过dockerfile创建镜像
    Create() // 创建容器
    Run() // 运行容器
}
```
* runner 实现 worker 逻辑
``` go
type Runner interface {
    GetDockerfile()
    Start()
    Run()
}
```
* server 实现对外接口提供服务  
``` go
type Server interface {
    Status()
    RunCode()
}
```
## Structure
```

```
