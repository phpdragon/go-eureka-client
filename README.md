# go-eurake-client

#### 介绍
Golang 实现的 Eurake Client

#### 软件架构

支持的Eureka server Rest api，参见下面的列表:

| Operation | HTTP action | Support |
|-----------|-------------|-------------|
| Register new application instance | POST /eureka/v2/apps/**appID** | √ |
| De-register application instance | DELETE /eureka/v2/apps/**appID**/**instanceID** | √ |
| Send application instance heartbeat | PUT /eureka/v2/apps/**appID**/**instanceID** | √ |
| Query for all instances | GET /eureka/v2/apps | √ |
| Query for all **appID** instances | GET /eureka/v2/apps/**appID** | √ |
| Query for a specific **appID**/**instanceID** | GET /eureka/v2/apps/**appID**/**instanceID** | √ |
| Query for a specific **instanceID** | GET /eureka/v2/instances/**instanceID** | √ |
| Take instance out of service | PUT /eureka/v2/apps/**appID**/**instanceID**/status?value=OUT_OF_SERVICE| √ |
| Move instance back into service (remove override) | DELETE /eureka/v2/apps/**appID**/**instanceID**/status?value=UP  (The value=UP is optional, it is used as a suggestion for the fallback status due to removal of the override)| √ |
| Update metadata | PUT /eureka/v2/apps/**appID**/**instanceID**/metadata?key=value| √ |
| Query for all instances under a particular **vip address** | GET /eureka/v2/vips/**vipAddress** | √  |
| Query for all instances under a particular **secure vip address** | GET /eureka/v2/svips/**svipAddress** | √  |


#### 安装教程

1.  xxxx
2.  xxxx
3.  xxxx

#### 使用说明

```java
clientConfig, _ := eureka.LoadConfig("etc/app.yaml", false)

// create eureka client
eurekaClient = eureka.NewClientWithLog(clientConfig, logger.GetLogger())
eurekaClient.Run()
//eurekaClient.Shutdown()

// http server
http.HandleFunc("/actuator/info", func(writer http.ResponseWriter, request *http.Request) {
	writeJsonResponse(writer, request, eureka.ActuatorStatus(8080, "go-example"), true)
})
http.HandleFunc("/actuator/health", func(writer http.ResponseWriter, request *http.Request) {
	writeJsonResponse(writer, request, eureka.ActuatorHealth(), true)
})
http.HandleFunc("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write(gFaviconIco)
	if err != nil {
		logger.Info(err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
})
http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	indexHandler(writer, request, eurekaClient)
})

// start http server
if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
	log.Fatal(err)
}
```

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 感谢

1. [https://github.com/xuanbo/eureka-client](https://github.com/xuanbo/eureka-client)
2. [https://github.com/HikoQiu/go-eureka-client](https://github.com/HikoQiu/go-eureka-client)

#### 用例

1. [gateway_proxy(转发代理)](https://gitee.com/phpdragon/gateway_proxy)

#### 码云特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  码云官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解码云上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是码云最有价值开源项目，是码云综合评定出的优秀开源项目
5.  码云官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  码云封面人物是一档用来展示码云会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
