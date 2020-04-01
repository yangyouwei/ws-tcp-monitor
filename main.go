package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/yangyouwei/ws-tcp-monitor/configlib"
	"github.com/yangyouwei/ws-tcp-monitor/loglib"
	"github.com/yangyouwei/ws-tcp-monitor/tcplib"
	"github.com/yangyouwei/ws-tcp-monitor/wslib"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

const ShellToUse = "bash"

var MainConfig *configlib.MainConf
var ServiceConfig *[]configlib.Service
var logtofile *log.Logger

func init() {
	s := flag.String("c", "./config.ini", "-c /etc/wstcp_monitor/config.ini")
	flag.Parse()
	//解析参数
	if *s == "" {
		flag.Usage()
		panic("process exit!")
	}

	//初始化config
	configlib.Initconfig(s)
	MainConfig = &configlib.MainConfig
	ServiceConfig = &configlib.ServicesConf
	//初始化log
	loglib.InitLog(MainConfig)
	logtofile = loglib.Logtofile
	logtofile.Println("monitor is starting.")
}

//判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//调用系统命令
func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func main() {
	var wg sync.WaitGroup
	//fmt.Println(len(*ServiceConfig))
	wg.Add(len(*ServiceConfig))
	for _, service := range *ServiceConfig {
		go func(L *log.Logger, wg *sync.WaitGroup, s configlib.Service) {
			ticker := time.NewTicker(MainConfig.Interval)
			defer wg.Done()
			for {
				select {
				case <-ticker.C:
					checkres := checkservice(s)
					if checkres != true {
						fmt.Println("服务不可用！重启服务，正在运行命令：",s.ServiceCmd)
						logtofile.Println("服务不可用！重启服务，正在运行命令：",s.ServiceCmd)
						restart_service(s.ServiceCmd, L)
						time.Sleep(5 * time.Second)
					}
				}
			}
		}(logtofile, &wg, service)
	}
	wg.Wait()
}

func checkservice(serviceconf configlib.Service) bool {
	fullpack := tcplib.MakePackage()
	//send package to server
	tcpres := tcplib.TcpConnect(serviceconf.TcpAddr, fullpack)

	//ws connect
	wsmesages := wslib.WSMessages{2}
	wsres := wslib.WebSocket(wsmesages, "ws://"+serviceconf.WsAddr)

	if tcpres && wsres {
		return true
	} else {
		return false
	}
}

func restart_service(cmd string, l *log.Logger) {
	err, sout, serr := Shellout(cmd)
	if err != nil {
		fmt.Println("1", err)
		l.Println(err)
	}
	if sout != "" {
		fmt.Println("2", err)
		l.Println(sout)
	}
	if serr != "" {
		fmt.Println("3", err)
		l.Println(serr)
	}
}
