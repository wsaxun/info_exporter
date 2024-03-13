package configs

import (
	"crypto/md5"
	"encoding/hex"
	"info_exporter/pkg/tools"
	"os"
	"strconv"
	"time"
)

// LISTENSERSER 本服务监听地址及端口
const LISTENSERSER string = ":8080"

// DEVDOMAIN dev环境
const DEVDOMAIN string = "http://test-dev.com"

// PRODDOMAIN prod环境
const PRODDOMAIN string = "https://test-prod.com"

// HOSTPATH 主机path
const HOSTPATH string = "/api/servers/monitors"

// SWITCH 交换机path
const SWITCH string = "/api/switches/list"

var ECDNKEY string
var DOMAIN string
var mfyEnv string
var ECDNCYCLE time.Duration
var cycle int = 300

// ConfInit 初始化相关配置变量
func ConfInit() {
	mfyEnv = os.Getenv("MFYENV")
	if mfyEnv == "" {
		mfyEnv = "dev"
	}
	DOMAIN = DEVDOMAIN
	if mfyEnv == "prod" {
		DOMAIN = PRODDOMAIN
	}
	ECDNKEY = os.Getenv("MFYKEY")
	num, err := strconv.Atoi(os.Getenv("MFYCYCLE"))
	if err == nil {
		cycle = num
	}
	ECDNCYCLE = time.Second * time.Duration(cycle)
	tools.LogInfo(map[string]interface{}{"env": mfyEnv, "cycle": cycle})
}

// GenToken 生成token
func GenToken(t int64) string {
	tokenData := []byte(ECDNKEY + strconv.FormatInt(t, 10))
	token := md5.Sum(tokenData)
	tokenStr := hex.EncodeToString(token[:])
	return tokenStr
}
