package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	printtitle("Starting server on port 8080")

	http.HandleFunc("/", index)
	http.HandleFunc("/receive", receive)

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	printtitle("New client connected")

	t, _ := template.ParseFiles("www/index.html")
	t.Execute(w, nil)
}

func receive(w http.ResponseWriter, req *http.Request) {
	ip, port, _ := net.SplitHostPort(req.RemoteAddr)
	userIP := net.ParseIP(ip)

	printtitle(userIP.String() + ":" + port)

	req.ParseForm()
	//for k, v := range req.Form {
	//fmt.Println(k, ": ", strings.Join(v, ""))
	//}

	nb := req.Form.Get("nb")
	print("Number of char", nb)

	val := req.Form.Get("ukrast")
	print("raw", val)

	uns := unescape(val)
	print("unescape", uns)

	hex := hex(val)
	print("hex", hex)

	fmt.Print("\n\n")
}

func printtitle(str string) {
	len := len(str) + 4
	i := 0

	for i < len {
		fmt.Print("#")
		i++
	}

	fmt.Print("\n")
	fmt.Println("| " + str + " |")
	fmt.Print("\n")
}

func print(key string, str string) {
	fmt.Print("| ")
	fmt.Print(key, ": ", str)
	fmt.Print("\n")
}

func unescape(str string) string {
	str, _ = url.QueryUnescape(str)
	return str
}

func hex(str string) string {
	var output []string

	for i := 0; i < len(str); i++ {
		output = append(output, fmt.Sprintf("%x", str[i]))
	}

	str = strings.Join(output, "")

	return str
}
