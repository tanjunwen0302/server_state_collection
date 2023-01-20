package main

import "fmt"

func main() {
	config := configInfo()
	if config["mode"] == "client" {
		client(config)
	} else if config["mode"] == "server" {
		server(config)
	} else {
		fmt.Errorf("错误")
	}
}
