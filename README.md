# StariverOCR-BestNodeChecker
星河最佳节点检测（就是个DDNS工具）

# 初始化并部署
## 配置
### HandleDomain(操作域名)
操作域名指的是程序将在运行过程中改变解析记录的域名，它将与本地解析的ResolveDomain同步。
### ResolveDomain(解析域名)
解析域名指的是程序在运行过程中解析的域名，程序将会查询它的DNS结果并同步至HandleDomain。
### Token
DNSPod API的Token，可前往[DNSPod控制台](https://console.dnspod.cn/)获取。
### 样例
```
Token: 114514,1919810
HandleDomain:
  RootDomain: daoxiangcity.com
  SubDomain: test-0
ResolveDomain:
  RootDomain: daoxiangcity.com
  SubDomain: test-1
```