## fpp (free proxy pool)免费代理池

## 介绍
这是一个基于Golang开发的免费HTTP代理池。定时采集网上发布的免费代理验证并入库，同时你也可以在`extractors`下扩展代码来增加代理IP池的数量。

## 下载与运行

### 下载

前往`https://github.com/HaliComing/fpp/releases`，下载适用于您目标机器操作系统、CPU架构的主程序，解压直接运行即可。

### 运行

```shell
# 解压程序包
tar -zxvf fpp_VERSION_OS_ARCH.tar.gz

# 赋予执行权限
chmod +x ./fpp

# 启动 fpp
./fpp
```

## 配置

首次直接运行会生成`conf.ini`配置文件，首次请勿自己创建。

```ini
[System]
; 是否Debug运行，默认false
; Debug = false
; api的HTTP监听端口，默认9826
Listen = :9826
; 检测线程数，默认50
NumberOfThreads = 70
; Token 默认会随机生成
Token = STO64ysaLOte9J732YF1aw1Gf17xsnTV
; IP提取时间间隔，单位分钟，默认60
ExtractionInterval = 45
; 检测时间间隔，单位分钟，默认15
CheckInterval = 10

[Database]
; 数据库类型，支持memory和sqlite3，默认sqlite3
; 由于memory类型会导致no such table，暂时不推荐，如果您能解决欢迎issue讨论。
Type = sqlite3
; Type为sqlite3时启用此字段，数据库名称
DBFile = fpp.db

[SSL]
; api的HTTPS监听端口
Listen = :443
; 证书位置
CertPath = cert
; key位置
KeyPath = key
```

## 接口

简单的接口文档，详情请运行`fpp`后打开`/`即可查看。例如：`http://localhost:9826/`

| api | method | Description |
| ----| ---- | ---- |
| / | GET | api介绍和文档 |
| /api/v1/site/ping | GET | 服务连通测试 |
| /api/v1/site/count | GET | 统计接口 |
| /api/v1/proxy/random | GET | 随机获取一个IP接口 |
| /api/v1/proxy/all | GET | 分页查询全部IP的接口  |
| /api/v1/proxy/delete | GET | 删除指定IP  |

## 构建

自行构建前需要拥有 `Go >= 1.13` 必要依赖。

#### 克隆代码

```shell
git clone https://github.com/HaliComing/fpp.git
```

#### 嵌入静态资源

```shell
# 回到项目主目录
# 将静态资源copy在assets/build/目录下

# 安装 statik, 用于嵌入静态资源
go get github.com/rakyll/statik

# 开始嵌入
statik -src=assets/build/  -include=*.html,*.ico,*.icon -f
```

#### 编译项目
您可以选择在Releases界面下载已编译好的二进制文件。手动编译如下：
```shell
# 开始编译
# 必须开启CGO_ENABLED
SET CGO_ENABLED=1

go build -a -o fpp -ldflags "-s"
```

## 提取器目录

不分先后顺序，如果您有更好的代理网站欢迎提交issue

| 代理名称               | 代理网址               | 提取器包名    |
| ---------------------- | ---------------------- | ------------- |
| 齐云代理               | www.7yip.cn            | yip7          |
| 66免费代理             | www.66ip.cn            | ip66          |
| 89免费代理             | www.89ip.cn            | ip89          |
| 全网代理IP             | www.goubanjia.com      | goubanjia     |
| IP3366云代理           | www.ip3366.net         | ip3366        |
| 高可用全球免费代理IP库 | ip.jiangxianli.com     | jiangxianli   |
| 快代理                 | www.kuaidaili.com      | kuaidaili     |
| ProxyList+             | list.proxylistplus.com | proxylistplus |
| 方法SEO代理            | seofangfa.com          | seofangfa     |
| 速代理                 | www.sudaili.com        | sudaili       |
| 西拉免费代理IP         | www.xiladaili.com      | xiladaili     |

## 添加自定义提取器

首先download代码后，在`extractors`文件夹中创建需要抓取的代理网站文件夹和go文件，并实现`Extract`接口，以下以`example.com`举例。

```
目录
extractors  #提取器文件夹
|--example  #示例提取器
|--|--example.go  #示例提取器代码
|--|--example_test.go #示例提取器测试代码
|--extractors.go #提取者
```

`example.go`，只需要修改`Key`接口和`Extract`接口即可。

```go
package example

import (
	"github.com/HaliComing/fpp/extractors/types"
	"github.com/HaliComing/fpp/models"
	"github.com/HaliComing/fpp/pkg/request"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strings"
)

type extractor struct{}

// New returns a extractor. 
func New() types.Extractor {
	return &extractor{}
}

// 修改此方法
func (e *extractor) Key() string {
    // 改为代理网站的域名，不带http前缀
	return "www.example.com"
}

// Extract is the main function to extract the data.
func (e *extractor) Extract() ([]*models.ProxyIP, error) {
	var proxyIPs []*models.ProxyIP
	// 抓取数据的代码，可参考其他提取器代码，代码很简陋，有更好的可以提issue或pr
    // ...
    // proxyIP := &models.ProxyIP{
	//		IP:   ip, //只需要填充IP和Port，其他勿填
	//		Port: port,//只需要填充IP和Port，其他勿填
	//	}
    //proxyIPs = append(proxyIPs, proxyIP)
	return proxyIPs, nil
}

```

`example_test.go`别忘了测试类，如无其他情况测试类无需改动。

```go
package example

import (
	"fmt"
	"testing"
)

// 测试提取
func TestExtract(t *testing.T) {
	extract, err := New().Extract()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, ip := range extract {
            // 打印出提取出来的代理
			fmt.Printf("%s:%s\n", ip.IP, ip.Port)
		}
	}
}
```

在`extractors.go`的`init`方法添加示例提取器。

```go
func init() {
	extractors = []types.Extractor{
		yip7.New(),
		ip66.New(),
		ip89.New(),
		goubanjia.New(),
		ip3366.New(),
		jiangxianli.New(),
		kuaidaili.New(),
		proxylistplus.New(),
		seofangfa.New(),
		sudaili.New(),
		xiladaili.New(),
		example.New(),// 添加示例提取器，注意后面的逗号
	}
}
```

## 未来规划(最近在忙其他项目，此暂停维护一段时间，有能力可以clone下来改改。2022/03/28)

- **集成代理自动切换IP**：其他需要代理IP的项目再也不用主动维护代理池啦，所有请求设置这个fpp代理后，fpp会自动为每个请求切换IP。
- **全球IP扫描模块**：自动扫描可用代理并验证入库
- **付费代理集成**：在自动切换IP的基础上集成付费代理，美滋滋，需要使用代理的程序只需要请求时设置fpp的代理后就可以直接使用高可用的代理了。

## 最后

- 首先感谢您的使用，如果喜欢本程序，不妨给个star
- 如果发现bug或者建议，欢迎提交issue
- 如果愿意贡献代码那就更棒啦，欢迎提交Pr
