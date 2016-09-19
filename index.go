package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

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

func showMarkdownJs(w http.ResponseWriter, r *http.Request) {
	text, err := readFile("./static/marked.js")
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("not found!"))
	} else {
		w.Header().Set("Content-Type", "application/x-javascript")
		w.Write(text)
	}

}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func readme(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.RequestURI()

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
		w.Write([]byte(t))
	}
}

func main() {
	http.HandleFunc("/static/marked.js", showMarkdownJs)
	http.HandleFunc("/", readme)
	http.ListenAndServe(":9000", nil)
}
