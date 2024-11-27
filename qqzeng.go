package iplocation

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func (m *Result) ToString() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s", m.Continents, m.Country, m.Province, m.City, m.Zone, m.ISP, m.Code, m.EnName, m.ShortName, m.Lng, m.Lat)
}

func (m *Result) ToJSON() string {
	b, _ := json.Marshal(&m)
	return string(b)
}

type ipIndex struct {
	startIp     uint32
	endIp       uint32
	localOffset uint32
	localLength uint32
}

type prefixIndex struct {
	startIndex uint32
	endIndex   uint32
}

type IPSearch struct {
	data               []byte
	prefixMap          map[uint32]prefixIndex
	firstStartIpOffset uint32
	prefixStartOffset  uint32
	prefixEndOffset    uint32
	prefixCount        uint32
}

var ipSearch *IPSearch

// NewIPLocation returns *IPSearch
//
//	dataFile the address of the ip data file
func NewIPLocation(dataFile string) (*IPSearch, error) {
	if dataFile == "" {
		return nil, errors.New("need ip data filepath")
	}
	if ipSearch == nil {
		var err error
		ipSearch, err = loadDataFile(dataFile)
		if err != nil {
			return nil, err
		}
	}
	return ipSearch, nil
}

// Get ipv4地址的查询
// result 为空时，返回错误.
func (p *IPSearch) Get(ip string) (result *Result) {
	result = new(Result)
	if ip == "" {
		return
	}
	//ipv4合法性
	if strings.Count(ip, ".") != 3 {
		return
	}
	ip = strings.TrimSpace(ip)
	s := p.get(ip)
	ss := strings.Split(s, "|")
	if len(ss) == 11 {
		result = &Result{
			Continents: ss[0],
			Country:    ss[1],
			Province:   ss[2],
			City:       ss[3],
			Zone:       ss[4],
			ISP:        ss[5],
			Code:       ss[6],
			EnName:     ss[7],
			ShortName:  ss[8],
			Lng:        ss[9],
			Lat:        ss[10],
		}
		return
	}
	return
}

func (p *IPSearch) get(ip string) (ret string) {
	ss := strings.Split(ip, ".")
	x, _ := strconv.Atoi(ss[0])
	prefix := uint32(x)
	intIp := ipToLong(ip)

	var low, high uint32
	if _, ok := p.prefixMap[prefix]; ok {
		low = p.prefixMap[prefix].startIndex
		high = p.prefixMap[prefix].endIndex
	} else {
		return
	}

	var myIndex uint32
	if low == high {
		myIndex = low
	} else {
		myIndex = p.binarySearch(low, high, intIp)
	}
	index := &ipIndex{}
	index.getIndex(myIndex, p)
	if index.startIp <= intIp && index.endIp >= intIp {
		return index.getLocal(p)
	}
	return
}

func loadDataFile(dataFile string) (*IPSearch, error) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	p := &IPSearch{}
	p.data = data
	p.prefixMap = make(map[uint32]prefixIndex)

	p.firstStartIpOffset = bytesToLong(data[0], data[1], data[2], data[3])
	p.prefixStartOffset = bytesToLong(data[8], data[9], data[10], data[11])
	p.prefixEndOffset = bytesToLong(data[12], data[13], data[14], data[15])
	p.prefixCount = (p.prefixEndOffset-p.prefixStartOffset)/9 + 1

	indexBuffer := p.data[p.prefixStartOffset:(p.prefixEndOffset + 9)]
	for k := uint32(0); k < p.prefixCount; k++ {
		i := k * 9
		prefix := uint32(indexBuffer[i] & 0xFF)
		pf := prefixIndex{}
		pf.startIndex = bytesToLong(indexBuffer[i+1], indexBuffer[i+2], indexBuffer[i+3], indexBuffer[i+4])
		pf.endIndex = bytesToLong(indexBuffer[i+5], indexBuffer[i+6], indexBuffer[i+7], indexBuffer[i+8])
		p.prefixMap[prefix] = pf
	}
	return p, nil
}

// binarySearch
// 二分逼近算法
func (p *IPSearch) binarySearch(low uint32, high uint32, k uint32) uint32 {
	var M uint32 = 0
	for low <= high {
		mid := (low + high) / 2

		endIpNum := p.getEndIp(mid)
		if endIpNum >= k {
			M = mid
			if mid == 0 {
				break
			}
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return M
}

func (p *IPSearch) getEndIp(left uint32) uint32 {
	leftOffset := p.firstStartIpOffset + left*12
	return bytesToLong(p.data[4+leftOffset], p.data[5+leftOffset], p.data[6+leftOffset], p.data[7+leftOffset])

}

func (p *ipIndex) getIndex(left uint32, ips *IPSearch) {
	leftOffset := ips.firstStartIpOffset + left*12
	p.startIp = bytesToLong(ips.data[leftOffset], ips.data[1+leftOffset], ips.data[2+leftOffset], ips.data[3+leftOffset])
	p.endIp = bytesToLong(ips.data[4+leftOffset], ips.data[5+leftOffset], ips.data[6+leftOffset], ips.data[7+leftOffset])
	p.localOffset = bytesToLong3(ips.data[8+leftOffset], ips.data[9+leftOffset], ips.data[10+leftOffset])
	p.localLength = uint32(ips.data[11+leftOffset])
}

func (p *ipIndex) getLocal(ips *IPSearch) string {
	bytes := ips.data[p.localOffset : p.localOffset+p.localLength]
	return string(bytes)

}

// ipToLong ip to long
//
//	returns uint32
func ipToLong(ip string) uint32 {
	quads := strings.Split(ip, ".")
	var result uint32 = 0
	a, _ := strconv.Atoi(quads[3])
	result += uint32(a)
	b, _ := strconv.Atoi(quads[2])
	result += uint32(b) << 8
	c, _ := strconv.Atoi(quads[1])
	result += uint32(c) << 16
	d, _ := strconv.Atoi(quads[0])
	result += uint32(d) << 24
	return result
}

// bytesToLong byte to long
//
//	returns uint32
func bytesToLong(a, b, c, d byte) uint32 {
	a1 := uint32(a)
	b1 := uint32(b)
	c1 := uint32(c)
	d1 := uint32(d)
	return (a1 & 0xFF) | ((b1 << 8) & 0xFF00) | ((c1 << 16) & 0xFF0000) | ((d1 << 24) & 0xFF000000)
}

func bytesToLong3(a, b, c byte) uint32 {
	a1 := uint32(a)
	b1 := uint32(b)
	c1 := uint32(c)
	return (a1 & 0xFF) | ((b1 << 8) & 0xFF00) | ((c1 << 16) & 0xFF0000)
}
