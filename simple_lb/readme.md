## 简易负载均衡器

轮询  
Go 语言标准库里的 ReverseProxy  
mutex  
原子操作  
闭包  
回调  
select  

这个简单的负载均衡器还有很多可以改进的地方：

使用堆来维护后端的状态，以此来降低搜索成本  
收集统计信息  
实现加权轮询或最少连接策略  
支持文件配置  

