package main

import (
	"bytes"
	"crypto/tls"
	"github.com/tidwall/sjson"
	"net/http"
	"time"
)

func client(config map[string]string) {
	ticker := time.NewTicker(4 * time.Second)
	for range ticker.C {
		dataUpload(config)
	}
}

// dataUpload 客户端数据上传
func dataUpload(config map[string]string) {
	dataJson := dataAcquisition()
	dataJson, _ = sjson.Set(dataJson, "client.key", config["clientKey"])
	dataJson, _ = sjson.Set(dataJson, "client.username", config["username"])
	data := bytes.NewBuffer([]byte(dataJson))

	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// http cookie接口
	// cookieJar, _ := cookiejar.New(nil)

	// c := &http.Client{
	// 	Jar:       cookieJar,
	// 	Transport: tr,
	// }
	// c.Post("https://"+config["serverIp"]+"/dataReception", "application/json", data)

	client := &http.Client{Transport: tr}
	request, _ := http.NewRequest("POST", "https://"+config["serverIp"]+"/dataReception", data)
	request.Header.Set("Connection", "close")
	client.Do(request)
	defer request.Body.Close()
	defer client.CloseIdleConnections()
}
