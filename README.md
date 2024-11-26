# ip-location

离线的IP归属地查询

## 安装及使用

安装

```shell
go get github.com/zituocn/iplocation
```
使用

```go
package iplocation

import (
	"fmt"
	"testing"
)

func TestIPSearch(t *testing.T) {
	ip, err := NewIPLocation("/Users/samsong/mygo/src/github.com/zituocn/iplocation/data/qqzeng-ip-china-utf8.dat")
	if err != nil {
		t.Errorf(err.Error())
	}
	result := ip.Get("218.88.127.69")
	fmt.Printf("struct: %v \n", result)
	fmt.Printf("ToJSON: %v \n", result.ToJSON())
	fmt.Printf("ToString: %v \n", result.ToString())
}
```
输出
```shell
struct: &{亚洲 中国 四川 成都 双流 中国电信 510116 China CN 103.9237 30.5744} 
ToJSON: {"continents":"亚洲","country":"中国","province":"四川","city":"成都","zone":"双流","isp":"中国电信","code":"510116","en_name":"China","short_name":"CN","lng":"103.9237","lat":"30.5744"} 
ToString: 亚洲|中国|四川|成都|双流|中国电信|510116|China|CN|103.9237|30.5744 
```

返回结构说明

```go

// Result the structure returned by the query
type Result struct {

	// Continents 大洲
	Continents string `json:"continents"`

	// Country 国家
	Country string `json:"country"`

	// Province 省份
	Province string `json:"province"`

	// City 市
	City string `json:"city"`

	// Zone 区县
	Zone string `json:"zone"`

	// ISP运营商
	ISP string `json:"isp"`

	// Code  行政代码
	Code string `json:"code"`

	//EnName 英语名
	EnName string `json:"en_name"`

	// 英文简称
	ShortName string `json:"short_name"`

	// Lng 经度
	Lng string `json:"lng"`

	// Lat 纬度
	Lat string `json:"lat"`
}

```

## 获取离线数据文件

* https://www.qqzeng.com/



