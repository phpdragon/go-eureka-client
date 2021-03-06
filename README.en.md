# go-eureka-client

#### Description
English | [简体中文](./README.md)

Golang implementation of the unofficial Spring Cloud Eureka client. 
> Tips: Non-full-features, only some basic and useful features implemented.

#### Configuration and features

Spring Cloud Eureka Configurations:

| Configuration | Support |
|-----------|-------------|
|availabilityZones| × |
|serviceUrl| √ |
|useDnsForFetchingServiceUrls| × |
|preferSameZoneEureka| × |
|filterOnlyUpInstances| √ |
|registryFetchIntervalSeconds| √ |
|fetchRegistry| √ |
|registerWithEureka| √ |
|shouldUnregisterOnShutdown| √ |
|instanceEnabledOnInit| √ |
|renewalIntervalInSecs| √ |

**go-eureka-client** Supported and extended features, refer to list below:

| Feature | Support |
|-----------|-------------|
| RegisterWithEureka | √ | 
| RegistryFetchIntervalSeconds | √ |
| FetchRegistry | √ |
| Regular reconnection | √ |
| Failure to retry | √ |
| Registration redirection | × |
| HeartbeatIntervals | √ |

[Eureka server Rest api](https://github.com/Netflix/eureka/wiki/Eureka-REST-operations) supported, refer to list below:

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

#### Installation

1.  go get github.com/phpdragon/go-eurake-client
2.  import eureka "github.com/phpdragon/go-eureka-client"

#### Samples

```java
// create eureka client
eurekaClient = eureka.NewClientWithLog("config_sample.yaml", logger.GetLogger())
eurekaClient.Run()
//eurekaClient.Shutdown()

//httpUrl, _ := eurekaClient.GetRealHttpUrl("http://DEMO/action")
//fmt.Println(httpUrl)

// http server
http.HandleFunc("/actuator/info", func(writer http.ResponseWriter, request *http.Request) {
	writeJsonResponse(writer, request, eureka.ActuatorStatus(), true)
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

1.  Fork the repository
2.  Create Feat_xxx branch
3.  Commit your code
4.  Create Pull Request

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