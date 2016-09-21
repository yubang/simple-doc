package main

/**
 * 一个简易的文档展示系统
 * @author: yubang
 **/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

var ContentTypeInfoMap map[string]interface{} // 全局变量

// 初始化函数
func startInit() {
	text, _ := readFile("./config/contentType.db")
	ContentTypeInfoMap = make(map[string]interface{})
	json.Unmarshal(text, &ContentTypeInfoMap)
}

// 获取ContentType
func getContentType(filePath string) string {

	arrs := strings.Split(filePath, ".")
	urlSuffix := arrs[len(arrs)-1]
	for k, v := range ContentTypeInfoMap {
		if k == "."+urlSuffix {
			return v.(string)
		}
	}
	return "application/octet-stream"
}

// 获取不带后缀的文件名
func getBaseFileName(filePath string) string {
	fullFileName := path.Base(filePath)
	suffix := path.Ext(filePath)
	s := strings.TrimSuffix(fullFileName, suffix)
	return s
}

// 读取文件内容
func readFile(path string) ([]byte, error) {

	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

// 处理静态资源文件
func handleStatic(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path

	text, err := readFile("./" + filePath)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("not found!"))
	} else {
		w.Header().Set("Content-Type", getContentType(filePath))
		w.Write(text)
	}

}

func subString(str string, begin, length int) string {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func readme(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path

	if hasSuffix("/", filePath) {
		filePath += "index.md"
	}
	text, err := readFile("./db/" + filePath)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("not found!"))
	} else {
		w.Header().Set("Content-Type", "text/html")
		html, _ := readFile("./config/template.html")
		t := strings.Replace(string(html), "{{code}}", string(text), -1)
		t = strings.Replace(t, "{{title}}", getBaseFileName(filePath), -1)
		w.Write([]byte(t))
	}
}

func main() {
	startInit()
	http.HandleFunc("/static/", handleStatic)
	http.HandleFunc("/", readme)

	settingConf := make(map[string]interface{})
	t, _ := readFile("./config/setting.conf")
	json.Unmarshal(t, &settingConf)

	hostAndPort := settingConf["host"].(string) + ":" + settingConf["port"].(string)
	fmt.Print("服务器监听 " + hostAndPort)
	http.ListenAndServe(hostAndPort, nil)
}
