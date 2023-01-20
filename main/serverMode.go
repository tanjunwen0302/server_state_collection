package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

var serverKey string
var serverIp string

//数据接收
func dataReception(w http.ResponseWriter, req *http.Request) {
	all, _ := ioutil.ReadAll(req.Body)
	host := req.RemoteAddr
	dataString, _ := sjson.Set(string(all), "client.ip", host)
	username := gjson.Get(dataString, "client.username").Value().(string)
	key := gjson.Get(dataString, "client.key").Value().(string)
	if key != serverKey {
		return
	}
	//存储
	dbSave(username, dataString)
}

//数据发送，与用户端交互
func dataSending(w http.ResponseWriter, req *http.Request) {
	all, _ := ioutil.ReadAll(req.Body)
	if gjson.Get(string(all), "key").Value().(string) == serverKey {
		w.Write([]byte(dbRead()))
	}
}

func server(config map[string]string) {
	getIp()
	go nativeDataFetching(config)
	serverKey = config["serverKey"]
	http.HandleFunc("/dataReception", dataReception)
	http.HandleFunc("/dataSending", dataSending)
	fmt.Println("服务器启动了")
	err := http.ListenAndServeTLS(config["port"], "cert.pem", "private.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func nativeDataFetching(config map[string]string) {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		dataJson := dataAcquisition()
		dataJson, _ = sjson.Set(dataJson, "client.username", config["username"])
		dataJson, _ = sjson.Set(dataJson, "client.ip", serverIp)
		dbSave(config["username"], dataJson)
	}

}

func getIp() {
	res, err := http.Get("https://api.ipify.org")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	serverIp = string(ip)
}

func dbRead() string {
	db, err := sql.Open("sqlite3", "serverState.db")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer db.Close()
	rows, err := db.Query("select server_name,data,time from state")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	dataJson := ""
	len := 0
	for rows.Next() {
		var username string
		var data string
		var createTime string
		err = rows.Scan(&username, &data, &createTime)
		dataJson, _ = sjson.Set(dataJson, strconv.Itoa(len)+".username", username)
		dataJson, _ = sjson.Set(dataJson, strconv.Itoa(len)+".data", data)
		dataJson, _ = sjson.Set(dataJson, strconv.Itoa(len)+".time", createTime)
		len++

	}
	return dataJson
}

func dbSave(username, data string) {
	db, err := sql.Open("sqlite3", "serverState.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()
	localTime := time.Now()
	// 使用 date 函数将时间格式化为 YYYY-MM-DD 的字符串
	formattedDate := fmt.Sprintf("%s", localTime.Format("2006-01-02 15:04:05"))
	stmt, _ := db.Prepare("insert or replace into state(server_name,data,time) values(?,?,?)")
	_, err = stmt.Exec(username, data, formattedDate)

	if err != nil {
		fmt.Println(err)
		return
	}

}
