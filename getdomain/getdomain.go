/*
* @Author: Try
* @Date:   2021/5/5 18:31
 */
package getdomain

import (
	"bufio"
	"fmt"
	"github.com/ozgio/strutil"
	"io"
	"os"
	"strings"
)

func Getdomain(filename string, urls string) []string {
	var s []string
	ipps := strutil.Words("80,443")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("[info] 字典加载失败！")
		return nil
	}
	reader := bufio.NewReader(file)
	for {
		url, err := reader.ReadString('\n') //注意是字符
		str1 := strings.Replace(url, "\n", "", -1)
		str := strings.Replace(str1, "\r", "", -1)
		if err == io.EOF {
			file.Close()
		}
		if err != nil {
			break
		}
		for i := 0; i < len(ipps); i++ {
			//fmt.Printf("a 的值为: %d\n", a)
			ipp := ipps[i]
			s = append(s, str+"."+urls+":"+ipp+"/")
		}
		//fmt.Println(str)
	}
	return s
}
