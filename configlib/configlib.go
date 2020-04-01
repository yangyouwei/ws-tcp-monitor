package configlib

import (
	"github.com/Unknwon/goconfig"
	"log"
	"path/filepath"
	"time"
)
type MainConf struct {
	Interval time.Duration
	Log bool
	LogFile string
}

type Service struct {
	TcpAddr string
	WsAddr string
	ServiceCmd string
}

var MainConfig MainConf
var ServicesConf []Service

func Initconfig(c *string)  {
	cstr, err := filepath.Abs(*c)
	if err != nil {
		log.Panicln(err)
	}
	//初始化配置文件
	cfg, err := goconfig.LoadConfigFile(cstr)
	if err != nil {
		log.Println("读取配置文件失败[config.ini]")
		log.Panic(err)
	}
	//解析配置文件
	MainConfig.ReadConf(cfg)
	ReadServicesConf(cfg)
}

func (this *MainConf)ReadConf(c *goconfig.ConfigFile)  {
	var err error
	interval_num, err := c.GetValue("main","interval")
	if err != nil {
		log.Panic(err)
	}
	this.Interval, err = time.ParseDuration(interval_num)
	if err != nil {
		log.Panic(err)
	}

	logbool,err := c.GetValue("main","log")
	if err != nil {
		log.Panicln(err)
	}
	this.Log = ParseBool(logbool)

	this.LogFile,err = c.GetValue("main","log_file")
}

func ReadServicesConf(c *goconfig.ConfigFile)  {
	sections := c.GetSectionList()
	for _,section := range sections {
		if section == "main" {
			continue
		}else {
			var ServiceConfig Service
			ServiceConfig.ReadServiceConf(c,section)
			ServicesConf = append(ServicesConf,ServiceConfig)
		}
	}
}

func (this *Service)ReadServiceConf(c *goconfig.ConfigFile,n string)  {
	var err error
	seviceconf := c.GetKeyList(n)
	for _,keyname :=  range seviceconf {
		switch keyname {
		case "tcp_ip":
			this.TcpAddr,err = c.GetValue(n,keyname)
		case "ws_ip":
			this.WsAddr,err = c.GetValue(n,keyname)
		case "cmd":
			this.ServiceCmd, err = c.GetValue(n,keyname)
		default:
			continue
		}
	}
	if err !=  nil {
		log.Panicln(err)
	}
}

func ParseBool(str string) bool {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True":
		return true
	case "0", "f", "F", "false", "FALSE", "False":
		return false
	}
	return false
}
