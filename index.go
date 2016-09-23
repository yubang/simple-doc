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
	"strconv"
	"strings"
)

var contentTypeInfoMap map[string]interface{} // 全局变量
var settingConf map[string]interface{}        // 程序配置

// 初始化函数
func startInit() {
	text, _ := readFile("./config/contentType.db")
	contentTypeInfoMap = make(map[string]interface{})
	json.Unmarshal(text, &contentTypeInfoMap)
}

// 获取ContentType
func getContentType(filePath string) string {

	arrs := strings.Split(filePath, ".")
	urlSuffix := arrs[len(arrs)-1]
	for k, v := range contentTypeInfoMap {
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
		w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(int(settingConf["static_timeout"].(float64))))
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

// 读取markdown文件，并且渲染
func readme(w http.ResponseWriter, r *http.Request) {

	// 授权检验
	if settingConf["auth"].(bool) {
		cookie, err := r.Cookie(settingConf["auth_cookie_key"].(string))
		if err != nil || cookie.Value != settingConf["auth_token"] {
			w.Header().Set("Location", "/static/login.html")
			w.WriteHeader(302)
			return
		}
	}

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

// 处理登录问题
func handleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == settingConf["username"].(string) && password == settingConf["password"].(string) {
		cookie := http.Cookie{Name: settingConf["auth_cookie_key"].(string), Value: settingConf["auth_token"].(string), Path: "/", MaxAge: int(settingConf["auth_cookie_timeout"].(float64))}
		http.SetCookie(w, &cookie)
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte("fail"))
	}
}

func main() {
	startInit()
	http.HandleFunc("/static/", handleStatic)
	http.HandleFunc("/login.go", handleLogin)
	http.HandleFunc("/", readme)

	settingConf = make(map[string]interface{})
	t, _ := readFile("./config/setting.conf")
	json.Unmarshal(t, &settingConf)

	hostAndPort := settingConf["host"].(string) + ":" + settingConf["port"].(string)
	fmt.Print("服务器监听 " + hostAndPort)
	http.ListenAndServe(hostAndPort, nil)
}
