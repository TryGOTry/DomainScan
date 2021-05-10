package save

import (
	"fmt"
	"os"
	"time"
)

func init() {
	pwd, _ := os.Getwd()
	timeStr := time.Now().Format("2006-01-02")
	b, _ := PathExists(pwd + "/" + timeStr + "/")
	if b {
		//fmt.Println("ok")
	} else {
		err := os.Mkdir("./"+timeStr, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
}
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
func Savefile(url string,msg string) {
	pwd, _ := os.Getwd()
	timeStr := time.Now().Format("2006-01-02")
	okname := pwd + "\\" + timeStr + "\\" + url + ".csv"
	f, err := os.OpenFile(okname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	_, err = fmt.Fprintln(f, msg)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}