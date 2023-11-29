# Worker model
## Overview
该模块是最小的OJ单元，每个单元关联一个Container  
该单元接收提交的题目，返回输出+运行信息  
## Action
``` mermaid
---
title: Workflow
---
flowchart LR
    Factory--发送代码-->Worker
    Worker--发送运行性能参数+代码-->Analyzer
    Worker--返回运行结果-->Factory
```
## Functional Decomposition
* 容器操作
    * 构建容器
    * 执行代码
    * 外部接口
## Package
* docker 实现docker container的所有操作
``` go
type Container interface {
    Build()
    Create()
    Start()
}
```
* worker 实现worker相关业务组合  
``` go
type Runner interface {
    Run()
}
```
* server 实现对外接口提供服务  
``` go
type Server interface {
    Judge()
    PostPerformance()
}
```
## Structure
```
worker/  
├── Dockerfile worker容器的Dockerfile
├── doc 项目设计文档等
│   └── doc.md
├── docker  负责链接Docker沙箱等容器操作
│   ├── code  代码暂存区
│   │   └── main.go  
│   ├── container 容器内操作
│   │   └── run.sh
│   ├── docker.go
│   └── dockerfile 沙箱容器dockerfile
├── proto 负责外部接口/service
├── go.mod  
├── go.sum  
└── main.go 入口文件  
```
