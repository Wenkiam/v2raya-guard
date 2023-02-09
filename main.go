package main

import (
	"encoding/json"
	"flag"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	cronExp string
	logFile string
)

func init() {
	conf := flag.String("c", "config.json", "配置文件")
	flag.Parse()
	file, err := os.Open(*conf)
	if err != nil {
		panic("read config file failed:" + err.Error())
	}
	buf, err := io.ReadAll(file)
	if err != nil {
		panic("read config file failed:" + err.Error())
	}
	config := make(map[string]string)
	json.Unmarshal(buf, &config)
	username = config["username"]
	password = config["password"]
	logFile = config["logFile"]
	initServerBase(config["serverBase"])
	cronExp = config["cron"]
	initLogSetting()
	log.Println("load config success:" + string(buf))
}
func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	cron := cron.New()
	log.Println("add task.name:restart,cron:" + cronExp)
	cron.AddFunc(cronExp, restart)
	log.Println("add task.name:running check,cron:0 * * * *")
	cron.AddFunc("*/1 * * * *", func() {
		if !isRunning() {
			log.Println("v2ray is not running. try start")
			restart()
		}
	})
	cron.Start()
	wg.Wait()
}
func initServerBase(s string) {
	if strings.HasSuffix(s, "/api/") {
		serverBase = s
		return
	}
	if strings.HasSuffix(s, "/api") {
		serverBase = s + "/"
		return
	}
	if strings.HasSuffix(s, "/") {
		serverBase = s + "api/"
		return
	}
	serverBase = s + "/api/"
}
func initLogSetting() {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err == nil {
		log.SetOutput(file)
	}
}
