# Code runner
## 功能
本模块实现特定编程语言的Code runner factory  
期望状态为部署在集群中的一个Factory+多Wroker实例
## 特点
支持横向、动态拓展，沙箱运行，错误自动重启(K8s)  

Factory
* 定义Worker容器的编译环境
* 负责充当微服务中的消息中台，提高应对大流量场景的上限
* Factory对worker进行管理，维护多个Worker实例
* 分发文本形式代码给Worker运行，使用通过grpc形式与多个worker交互
* 将Worker返回值进行转发，并将资源等参数转发给Data Analysis Model
* 限流+负载均衡  

Worker
* 运行代码的基本单位，绑定唯一容器
* 维护Factory提交的多个代码，按提交顺序返回
# Data Analysis Model

## 功能
本模块实现对基本的代码性能、格式等进行分析，维护后台个人代码数据  
有时间的话
## 特点

# User
## 功能
感觉不如用Ruoyi，权限管理吃力没效果
## 特点
没啥特点，重心在Factory+Runner