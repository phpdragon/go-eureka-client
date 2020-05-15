# go-eureka-client

#### 介绍
简体中文 | [English](./README.en.md)

Golang 实现的非官方 Spring Cloud Eureka client.

#### 软件架构

支持的[Eureka server Rest api](https://github.com/Netflix/eureka/wiki/Eureka-REST-operations) ，参见下面的列表:

| 操作 | HTTP动作 | 支持 |
|-----------|-------------|-------------|
| 注册应用实例 | POST /eureka/v2/apps/**appID** | √ |
| 注销应用实例 | DELETE /eureka/v2/apps/**appID**/**instanceID** | √ |
| 发送心跳 | PUT /eureka/v2/apps/**appID**/**instanceID** | √ |
| 查询所有应用实例 | GET /eureka/v2/apps | √ |
| 通过 **appID** 查询所有实例 | GET /eureka/v2/apps/**appID** | √ |
| 通过 **appID**/**instanceID**  查询实例 | GET /eureka/v2/apps/**appID**/**instanceID** | √ |
| 通过 **instanceID** 查询实例 | GET /eureka/v2/instances/**instanceID** | √ |
| 从服务中取出实例 | PUT /eureka/v2/apps/**appID**/**instanceID**/status?value=OUT_OF_SERVICE| √ |
| 将实例移除(移除覆盖) | DELETE /eureka/v2/apps/**appID**/**instanceID**/status?value=UP  (The value=UP is optional, it is used as a suggestion for the fallback status due to removal of the override)| √ |
| 更新元数据 | PUT /eureka/v2/apps/**appID**/**instanceID**/metadata?key=value| √ |
| 通过 **vip address** 查询所有实例 | GET /eureka/v2/vips/**vipAddress** | √  |
| 通过 **secure vip address** 查询所有实例 | GET /eureka/v2/svips/**svipAddress** | √  |

支持的特性，见以下列表:
| 特性 | 支持 |
|-----------|-------------|
| 定期重连 | √  |
| 失败重试 | √  |
| 注册重定向 | × |

#### 安装教程

1.  xxxx
2.  xxxx
3.  xxxx

#### 使用说明

```java
// create eureka client
eurekaClient = eureka.NewClientWithLog("etc/app.yaml", logger.GetLogger())
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

1. [gateway_proxy(转发代理)](https://github.com/phpdragon/gateway_proxy)

#### 码云特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  码云官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解码云上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是码云最有价值开源项目，是码云综合评定出的优秀开源项目
5.  码云官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  码云封面人物是一档用来展示码云会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
