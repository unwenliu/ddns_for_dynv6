package main

import (
	"flag"
	"fmt"
	"github.com/asmcos/requests"
	"github.com/goinggo/mapstructure"
	"log"
	"net"
	"os"
	"time"
)

//dynv6的api地址
const Dynv6Url string = "https://dynv6.com/api/update"

var (
	help      bool   //帮助
	inter     string //网卡
	hostname  string //域名
	token     string //token
	ipv4      bool
	ipv6      bool
	show_ipv4 bool
	show_ipv6 bool
	timer     int
)

func init() {
	flag.BoolVar(&help, "h", false, "帮助")
	flag.StringVar(&inter, "i", "eth0", "要获取ip的网卡")
	flag.StringVar(&hostname, "hostname", "", "要更新的域名")
	flag.StringVar(&token, "token", "", "你的dynv6里的域名所绑定的token")
	flag.BoolVar(&ipv4, "4", false, "更新ipv4地址")
	flag.BoolVar(&ipv6, "6", false, "更新ipv6地址")
	flag.BoolVar(&show_ipv4, "show_ipv4", false, "显示指定网卡的ipv4地址")
	flag.BoolVar(&show_ipv6, "show_ipv6", false, "显示指定网卡的ipv6地址")
	flag.IntVar(&timer, "t", 300, "检查周期（秒），默认300秒")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `ddnsfordynv6 version: 0.1
使用说明：ddnsfordynv6 [-i 网卡名] [-hostname 域名] [-token token] [-4] [-6] [-t]
选项：
    -i 网卡名                  ip所绑定的网卡
    -show_ipv4                 显示指定网卡的ipv4地址
    -show_ipv6                 显示指定网卡的ipv6地址
    -hostname 域名             你的域名
    -token token               你的token
    -4                         更新ipv4地址
    -6                         更新ipv6地址
    -t                         检查周期（秒）,默认300秒
`)
}

// 获取给定网卡的ipv4地址和ipv6地址
func GetIP(interfaceName string) map[string]string {

	inter, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Fatalf("无法获取网卡信息，原因是:%v\n", err)
	}
	addrs, err := inter.Addrs()
	if err != nil {
		log.Fatalf("无法获取ip地址，原因是:%v\n", err)
	}
	ipaddr := make(map[string]string)
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//log.Printf("获取到网卡:%v ipv4地址:%v", inter.Name, ipnet.IP)
				ipaddr["ipv4"] = ipnet.IP.String()
			} else {
				//log.Printf("获取到网卡:%v ipv6地址:%v", inter.Name, ipnet.IP)
				ipaddr["ipv6"] = ipnet.IP.String()
			}
		}
	}
	return ipaddr
}

func res(config map[string]string) {
	p := requests.Params{}
	if err := mapstructure.Decode(config, &p); err != nil {
		log.Fatalf("转换struct失败,原因:%v\n", err)
	}
	resp, err := requests.Get(Dynv6Url, p)
	if err != nil {
		log.Fatalf("向dynv6发送请求失败,原因:%v\n", resp.Text())
	}
	if resp.R.StatusCode != 200 {
		if ipv4 && !ipv6 {
			log.Fatalf("更新ipv4失败,dynv6返回%v\n", resp.Text())
		} else if ipv6 && !ipv4 {
			log.Fatalf("更新ipv6失败,dynv6返回%v\n", resp.Text())
		} else if ipv4 && ipv6 {
			log.Fatalf("更新ipv4和ipv6失败,dynv6返回%v\n", resp.Text())
		}
	} else if ipv4 && !ipv6 {
		if config["ipv4"] != "" {
			log.Printf("更新ipv4地址成功,dyn6返回%v\n", resp.Text())
		} else {
			log.Fatalf("更新ipv4地址失败,原因未取到ipv4地址\n")
		}
	} else if ipv6 && !ipv4 {
		if config["ipv6"] != "" {
			log.Printf("更新ipv6地址成功,当前ipv6地址为%v,dynv6返回%v\n", config["ipv6"], resp.Text())
		} else {
			log.Fatalf("更新ipv6地址失败,原因未取到ipv6地址\n")
		}
	} else if ipv4 && ipv6 {
		if config["ipv4"] != "" && config["ipv6"] != "" {
			log.Printf("更新ipv4地址和ipv6地址成功,dynv6返回%v\n", resp.Text())
		} else {
			log.Fatalf("更新ipv4地址或ipv6地址失败,原因是ipv4地址或ipv6地址未取到\n")
		}
	}
}

func main() {
	// 命令行参数
	flag.Parse()
	if help {
		flag.Usage()
	} else if hostname != "" && token != "" {
		ipaddr := GetIP(inter)
		config := make(map[string]string)
		config["hostname"] = hostname
		config["token"] = token
		if ipv4 && !ipv6 {
			config["ipv4"] = ipaddr["ipv4"]
			res(config)
			for range time.Tick(time.Second * time.Duration(timer)) {
				res(config)
			}
		} else if ipv6 && !ipv4 {
			config["ipv6"] = ipaddr["ipv6"]
			res(config)
			for range time.Tick(time.Second * time.Duration(timer)) {
				res(config)
			}
		} else if ipv4 && ipv6 {
			config["ipv4"] = ipaddr["ipv4"]
			config["ipv6"] = ipaddr["ipv6"]
			res(config)
			for range time.Tick(time.Second * time.Duration(timer)) {
				res(config)
			}
		}
	} else if show_ipv4 {
		ipaddr := GetIP(inter)
		if ipaddr["ipv4"] != "" {
			fmt.Printf("获得网卡%v的ipv4地址为%v\n", inter, ipaddr["ipv4"])
		} else {
			fmt.Printf("获得网卡%v的ipv4地址失败\n", inter)
		}
	} else if show_ipv6 {
		ipaddr := GetIP(inter)
		if ipaddr["ipv6"] != "" {
			fmt.Printf("获得网卡%v的ipv6地址为%v\n", inter, ipaddr["ipv6"])
		} else {
			fmt.Printf("获得网卡%v的ipv6地址失败\n", inter)
		}
	} else {
		flag.Usage()
	}
}
