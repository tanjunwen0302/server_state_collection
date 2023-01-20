# server_state_collection

服务器状态监控（采集端）

# 使用说明

## 配置信息

### serverState.json文件

    {
        "mode": "server",
        "username": "windows主机",
        "client": {
        "serverIp": "127.0.0.1:3432",
        "clientKey": "tanjunwen"
        },
        "server": {
        "port": ":3432",
        "serverKey": "tanjunwen"
        }
    }
1、mode：运行模式（client或server）

2、username：主机名（多个服务器不可重复）

#### client模式
1、client.serverIp：将自己的信息推给的中心服务器的ip与端口
2、client.clientKey：将自己的信息推送给中心服务器的时用来认证的key值（一个字符串）

#### server模式
1、server.port：运行的端口
2、server.serverKey：用来认证客户端的key值（一个字符串）

## https通信使用的ssl证书

自行申请或直接使用提供的成品，如果自行申请将/main/serverMode.go文件下的48行

    err := http.ListenAndServeTLS(config["port"], "cert.pem", "private.key", nil)
上面的密钥信息更改你的


# 启动
当只有一台服务器时，直接将配置文件的mode修改为server模式，然后指定运行的端口“:port”（注意端口的数字前有:）

当有多台服务器时，将用来推送给手机app的那台服务器mode修改为server模式，将其他被采集的mode修改为client,将serverIp修改为server模式运行的服务器ip与端口号，将clientKey值修改为与serverkey值相同

然后

    nohup >log.out ./serverStateRun &
即可正常运行