package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"strings"
	"os"
	"io"
	"path/filepath"
	"regexp"
	"time"
	"github.com/jakecoffman/cron"
	"github.com/kardianos/service"
)

type Err struct {
	code int
	msg  string
}

type program struct{}

var (
	ErrRequest = Err{100, "Occur Error When Request"}
	ErrIO = Err{101, "Occur Error When Using IO"}
	ErrDone = Err{0, "Done Success"}
	ErrService = Err{102, "Occur Error When Start Service"}
)

//var LOGPATH = getCurrentDirectory() + "/" + "_.log"
var LOGPATH = "_.log"

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	// fetchLogo_Op()

	c := cron.New()
	c.AddFunc("@every 24h", fetchLogo_Op, "op.gg")
	c.Start()
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "GoServer",
		DisplayName: "GoServer",
		Description: "GoServer",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log(ErrService.msg, "op.gg")
	}
	err = s.Run()
	if err != nil {
		log(ErrService.msg, "op.gg")
	}
}

func fetchLogo_Op() {
	// fop := getCurrentDirectory() + "/image"
	fop := "image"
	if !checkFileIsExist(fop) {
		os.MkdirAll(fop, os.ModePerm)
	}

	u := "http://www.op.gg/"
	err, res := get(u)
	if err == ErrDone {
		//yt
		//fn := "_.tmp"
		//if checkFileIsExist(fn) {
		//	f, _ := os.OpenFile(fn,  os.O_CREATE, 0666)
		//	io.WriteString(f, string(res))
		//} else {
		//	f, _ := os.Create(fn)
		//	io.WriteString(f, string(res))
		//}

		pt := "https://attach.s.op.gg/logo/(.*?).PNG"
		reg := regexp.MustCompile(pt)
		up := reg.FindString(res)
		// fp := getCurrentDirectory() + "/image/" + time.Now().Format("20060102") + ".png"
		fp := fop + "/" + time.Now().Format("20060102") + ".png"
		err := download(up, fp)

		log(err.msg, "op.gg")
	} else {
		log(err.msg, "op.gg")
	}
}

func download(u string, p string) (e Err){
	res, err := http.Get(u)
	if err != nil {
		return ErrRequest
	}

	defer res.Body.Close()
	f, err := os.Create(p)
	if err != nil {
		return ErrIO
	}

	_, error := io.Copy(f, res.Body)
	if error != nil {
		return ErrIO
	}

	return ErrDone
}

func get(u string) (e Err, r string) {
	res, err := http.Get(u)
	if err != nil {
		return ErrRequest, ""
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ErrIO, ""
	}

	return ErrDone, string(body)
}

func post(u string) (e Err){
	res, err := http.Post(u,
		"application/x-www-form-urlencodeed",
		strings.NewReader("name=cjb"))
	if err != nil {
		return ErrRequest
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ErrIO
	}

	fmt.Println(string(body))

	return ErrDone
}

func checkFileIsExist(filename string) (bool) {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func getCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	return strings.Replace(dir, "\\", "/", -1)
}

func log(s string, t string) {
	if !checkFileIsExist(LOGPATH) {
		os.Create(LOGPATH)
	}

	f, _ := os.OpenFile(LOGPATH, os.O_APPEND, 0666)
	s = time.Now().Format("[2006/01/02 15:04:05]") + "[" + t + "]" + " " + s + "\r\n"
	io.WriteString(f, s)
}