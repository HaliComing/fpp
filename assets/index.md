[TOC]



## fpp代理池接口文档


### 1. 服务测试接口

#### 接口功能

> 测试服务是否在线

#### URL

> [/api/v1/site/ping](/api/v1/site/ping)

#### 支持格式

> JSON

#### HTTP请求方式

> GET

#### 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|token |是 |string|服务启动时配置的令牌 |

#### 返回字段

|返回字段|字段类型|说明 |
|:----- |:------|:----------------------------- |
|Code | int |状态码 0正常 非0异常 |
|Msg | string | 状态信息 |
|Data | string |版本号 |

#### 接口示例

> 地址：[http://www.example.com/api/v1/site/ping](http://www.example.com/api/v1/site/ping)
```json
{
    "Code":0,
    "Data":"0.0.1",
    "Msg":""
}
```


### 2. 统计接口

#### 接口功能

> 获取当前服务总IP数量

#### URL

> [/api/v1/site/count](/api/v1/site/count)

#### 支持格式

> JSON

#### HTTP请求方式

> GET

#### 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|token |是 |string|服务启动时配置的令牌 |

#### 返回字段

|返回字段|字段类型|说明 |
|:----- |:------|:----------------------------- |
|Code | int |状态码 0正常 非0异常 |
|Msg | string | 状态信息 |
|Data | string |版本号 |
|Count | int |IP总数 |

#### 接口示例

> 地址：[http://www.example.com/api/v1/site/count](http://www.example.com/api/v1/site/count)
```json
{
    "Code":0,
    "Data":{
        "Count":19
    },
    "Msg":""
}
```


### 3. 随机获取一个IP接口

#### 接口功能

> 随机获取一个IP接口

#### URL

> [/api/v1/proxy/random](/api/v1/proxy/random)

#### 支持格式

> JSON

#### HTTP请求方式

> GET

#### 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|token |是 |string|服务启动时配置的令牌 |

#### 返回字段

|返回字段|字段类型|说明 |
|:----- |:------|:----------------------------- |
|Code | int |状态码 0正常 非0异常 |
|Msg | string | 状态信息 |
|Data | string |版本号 |
|IP | string |IP地址 |
|Port | string |端口 |
|Protocol | int |支持协议 |
|Anonymous | int |匿名程度 |
|Country | string |国家 |
|Province | string |省份 |
|Attribution | string |归属地 |
|ISP | string |运营商 |
|Score | int |当前分数 满分10分 |
|Source | string |来源地址 |
|Speed | int |响应速度 单位ms |
|CreateTime | string |创建时间 |
|LastTime | string |最后检测时间 |
|FailedCount | int |连续失败次数 |

#### 接口示例

> 地址：[http://www.example.com/api/v1/proxy/random](http://www.example.com/api/v1/proxy/random)
```json
{
    "Code":0,
    "Data":{
        "IP":"113.214.13.1",
        "Port":"1080",
        "Protocol":1,
        "Anonymous":4,
        "Country":"中国",
        "Province":"浙江省",
        "Attribution":"宁波市",
        "ISP":"华数",
        "Score":10,
        "Source":"ip.jiangxianli.com",
        "Speed":1574,
        "CreateTime":"2021-03-16 09:21:22",
        "LastTime":"2021-03-16 09:31:36",
        "FailedCount":0
    },
    "Msg":""
}
```


### 4. 查询全部IP接口

#### 接口功能

> 分页查询全部IP的接口

#### URL

> [/api/v1/proxy/all](/api/v1/proxy/all)

#### 支持格式

> JSON

#### HTTP请求方式

> GET

#### 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|page |否 |int|页数，默认1开始 |
|size |否 |int|每页数量，默认10 |
|token |是 |string|服务启动时配置的令牌 |

#### 返回字段

|返回字段|字段类型|说明 |
|:----- |:------|:----------------------------- |
|Code | int |状态码 0正常 非0异常 |
|Msg | string | 状态信息 |
|Data | string |版本号 |
|IP | string |IP地址 |
|Port | string |端口 |
|Protocol | int |支持协议 |
|Anonymous | int |匿名程度 |
|Country | string |国家 |
|Province | string |省份 |
|Attribution | string |归属地 |
|ISP | string |运营商 |
|Score | int |当前分数 满分10分 |
|Source | string |来源地址 |
|Speed | int |响应速度 单位ms |
|CreateTime | string |创建时间 |
|LastTime | string |最后检测时间 |
|FailedCount | int |连续失败次数 |

#### 接口示例

> 地址：[http://www.example.com/api/v1/proxy/all](http://www.example.com/api/v1/proxy/all)
```json
{
    "Code":0,
    "Data":[
        {
            "IP":"113.214.13.1",
            "Port":"1080",
            "Protocol":1,
            "Anonymous":4,
            "Country":"中国",
            "Province":"浙江省",
            "Attribution":"宁波市",
            "ISP":"华数",
            "Score":10,
            "Source":"ip.jiangxianli.com",
            "Speed":1574,
            "CreateTime":"2021-03-16 09:21:22",
            "LastTime":"2021-03-16 09:31:36",
            "FailedCount":0
        }
    ],
    "Msg":""
}
```


### 5. 删除IP接口

#### 接口功能

> 删除指定IP，当获取IP使用无效时可使用此接口进行删除

#### URL

> [/api/v1/proxy/delete](/api/v1/proxy/delete)

#### 支持格式

> JSON

#### HTTP请求方式

> GET

#### 请求参数

|参数|必选|类型|说明|
|:----- |:-------|:-----|----- |
|ip |是 |string|欲删除IP地址 |
|port |是 |string|欲删除端口号 |
|proxy |否 |string|欲删除的代理，格式：ip:port，proxy参数存在时ip和port参数失效 |
|token |是 |string|服务启动时配置的令牌 |

#### 返回字段

|返回字段|字段类型|说明 |
|:----- |:------|:----------------------------- |
|Code | int |状态码 0正常 非0异常 |
|Msg | string | 状态信息 |
|Data | string |无 |

#### 接口示例

> 地址：[http://www.example.com/api/v1/proxy/delete](http://www.example.com/api/v1/proxy/delete)
```json
{
    "Code":0,
    "Data":"",
    "Msg":""
}
```

