package main

import (
	"domainScan/root"
	"flag"
)

func main() {
	//a := getdomain.Getdomain("t.txt", "baidu.com")
	//b := strings.Replace(a[1], ":443/", "", -1)
	//fmt.Println(b)
	//addr, err := net.ResolveIPAddr("ip", "www.baidu.com")
	//if err != nil {
	//}
	//if addr.String() != "<nil>" {
	//	fmt.Println(1)
	//} else {
	//	fmt.Println(2)
	//}
	//fmt.Println(addr.String())
	domain := flag.String("domain", "", "一级域名,如:baidu.com")
	filename := flag.String("f", "", "加载字典")
	flServerAddr  := flag.String("server", "8.8.8.8", "dns服务器")
	num := flag.Int("s", 5, "线程")
	flag.Parse()
	if *domain != "" && *filename != "" {
		root.DomainScan(*domain, *filename, *flServerAddr, *num,3)
	} else {
		flag.Usage()
	}

}
