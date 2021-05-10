/*
* @Author: Try
* @Date:   2021/5/5 18:21
 */
package root

import (
	"domainScan/getdomain"
	"domainScan/golimit"
	"domainScan/save"
	"domainScan/scan"
	"fmt"
	"github.com/gookit/color"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

func setTitle(title string) {
	kernel32, _ := syscall.LoadLibrary(`kernel32.dll`)
	sct, _ := syscall.GetProcAddress(kernel32, `SetConsoleTitleW`)
	syscall.Syscall(sct, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	syscall.FreeLibrary(kernel32)
}
func DomainScan(urls string, filename string, flServerAddr string, num int, timeout int64) {
	dicall := getdomain.Getdomain(filename, urls)
	//fmt.Println(dicall)
	color.Red.Println("---------------------------------------")
	color.Red.Println("[Info] domain-Scan|Try| By T00ls.Net; ")
	color.Red.Println("[Info] 一款子域名爆破工具 当前dns服务器:", flServerAddr)
	color.Red.Println("[Info] 爆破中.当前线程:", num)
	color.Red.Println("[Info] 成功加载字典数量:", len(dicall))
	color.Red.Println("---------------------------------------")
	//fmt.Println(dicall)
	g := golimit.NewG(num) //设置线程数量
	wg := &sync.WaitGroup{}
	beg := time.Now()
	var okurl int
	okurl = 0
	for i := 0; i < len(dicall); i++ {
		wg.Add(1)
		task := dicall[i]
		g.Run(func() {
			setTitle("当前目标:" + task + " 任务:" + fmt.Sprintf("%d", i) + "/" + fmt.Sprintf("%d", len(dicall)))
			respBody, err := scan.Goscan(task, flServerAddr, timeout)
			if err != nil {
				//color.Warn.Println("目标访问错误，可能被ban了！")
				wg.Done()
				return
			}
			if respBody.StatusCode == 200 || strings.Contains(respBody.Ipadd, " ") {
				color.Info.Println("[200] ", respBody.Res+"   [len]", respBody.Bodylen, "   [title]", respBody.Title, "   [server]", respBody.Server, "   [ip]", respBody.Ipadd)
				//writefile.Write(url, "[200] "+respBody.Res+"\n")
				save.Savefile(urls, respBody.Res+","+respBody.Title+","+fmt.Sprintf("%d",respBody.Bodylen)+","+respBody.Server+","+respBody.Ipadd+","+respBody.Okdomain)
				okurl = okurl +1
				color.Green.Println("--------------------------------------------------------------------")
			} else if respBody.StatusCode == 403 || strings.Contains(respBody.Ipadd, " ") {
				color.Warn.Println("[403] ", respBody.Res+"   [len]", respBody.Bodylen, "   [title]", respBody.Title, "   [server]", respBody.Server, "   [ip]", respBody.Ipadd)
				//writefile.Write(url, "[403] "+respBody.Res+"\n")
				save.Savefile(urls, respBody.Res+","+respBody.Title+","+fmt.Sprintf("%d",respBody.Bodylen)+","+respBody.Server+","+respBody.Ipadd+","+respBody.Okdomain)
				okurl = okurl +1
				color.Green.Println("--------------------------------------------------------------------")
			} else if respBody.StatusCode == 302 || strings.Contains(respBody.Ipadd, " ") {
				color.Warn.Println("[302] ", respBody.Res+"   [len]", respBody.Bodylen, "   [title]", respBody.Title, "   [server]", respBody.Server, "   [ip]", respBody.Ipadd)
				//writefile.Write(url, "[302] "+respBody.Res+"\n")
				save.Savefile(urls, respBody.Res+","+respBody.Title+","+fmt.Sprintf("%d",respBody.Bodylen)+","+respBody.Server+","+respBody.Ipadd+","+respBody.Okdomain)
				okurl = okurl +1
				color.Green.Println("--------------------------------------------------------------------")
			}
			wg.Done()
		})
	}
	wg.Wait()
	color.Red.Printf("[info] 爆破完成！当前用时: %fs", time.Now().Sub(beg).Seconds(),"成功发现域名数量:",okurl)
}
