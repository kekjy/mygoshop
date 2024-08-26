package main

import (
	"fmt"
	"mygoshop/common"
	"mygoshop/config"
	"net/http"
	"strconv"
	"time"
)

// 分布式验证
func (m *AccessControl) CheckPermission(r *http.Request) bool {
	// 获取用用户 uid
	uid, err := r.Cookie("uid")
	if err != nil {
		return false
	}

	// 采用一致性 hash 算法，根据用户 ID，判断具体机器
	Server := hashRing.Get(uid.Value)
	if Server == "EMPTY" {
		return false
	}

	// 是否本机
	if Server == localHost {
		// 执行本机数据读取和校验
		return m.FetchData(uid.Value)
	} else {
		// 代理访问结果
		return m.FetchRemoteData(Server, r)
	}
}

func (m *AccessControl) FetchData(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}

	dataRecord := m.Get(uidInt)
	if !dataRecord.IsZero() {
		if dataRecord.Add(time.Duration(config.Interval) * time.Second).After(time.Now()) {
			return false
		}
	}

	m.Put(uidInt)
	return true
}

func (m *AccessControl) FetchRemoteData(host string, request *http.Request) bool {
	hostUrl := "http://" + host + ":" + config.ValidateSet.Port + "/checkRight"
	response, body, err := common.GetCurl(hostUrl, request)
	if err != nil {
		fmt.Println("get curl error")
		return false
	}

	// 判断状态
	if response.StatusCode == 200 {
		return string(body) == "true"
	}
	return false
}
