# go-eureka-client

#### Description
English | [简体中文](./README.md)

Golang implementation of the unofficial Spring Cloud Eureka client. 

#### Software Architecture

Supported Eureka server Rest API, see the following list:

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

Supported features, see the following list:

| Feature | Support |
|-----------|-------------|
| Regular reconnection | √  |
| Failure to retry | √  |
| Registration redirection | × |

#### Installation

1.  xxxx
2.  xxxx
3.  xxxx

#### Samples

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

#### Contribution

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### Thanks

1. [https://github.com/xuanbo/eureka-client](https://github.com/xuanbo/eureka-client)
2. [https://github.com/HikoQiu/go-eureka-client](https://github.com/HikoQiu/go-eureka-client)

#### Use cases

1. [gateway_proxy(转发代理)](https://github.com/phpdragon/gateway_proxy)

#### Gitee Feature

1.  You can use Readme\_XXX.md to support different languages, such as Readme\_en.md, Readme\_zh.md
2.  Gitee blog [blog.gitee.com](https://blog.gitee.com)
3.  Explore open source project [https://gitee.com/explore](https://gitee.com/explore)
4.  The most valuable open source project [GVP](https://gitee.com/gvp)
5.  The manual of Gitee [https://gitee.com/help](https://gitee.com/help)
6.  The most popular members  [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)