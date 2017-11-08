package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"webtest/foo"
	"webtest/lib"
)

func login(w http.ResponseWriter, r *http.Request) {

	s, err := filepath.Abs("form.gtpl")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(s)

	if r.Method == "GET" {
		t, err := template.ParseFiles("./form.gtpl")
		if err != nil {
			log.Panicln(err)
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Panicln(err)
		}
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println(r.Form.Get("username"))
		fmt.Println(r.FormValue("username"))

		matched, _ := regexp.MatchString("^\\d+$", r.Form.Get("age"))
		log.Println("matched: ", matched)
	} else {
		err := r.ParseForm()
		handErr(err)
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		fmt.Println("username:", r.Form["username"])

		agestr := r.PostFormValue("age")
		if len(agestr) > 0 {
			age, err := strconv.Atoi(agestr)
			handErr(err)
			log.Println(age)
		}

		io.WriteString(w, username)
		io.WriteString(w, ":")
		io.WriteString(w, password)
	}
}

func handErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func han(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("text")
	// 是否是汉字
	isMatch, _ := regexp.MatchString("^\\p{Han}+$", text)
	fmt.Fprintf(w, "'%s' : %v", text, isMatch)
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	var cookie = &http.Cookie{}
	cookie.Expires = time.Now().AddDate(0, 1, 0)

	cookie.Name = "cookieTest"
	cookie.Value = "cookieValue:" + time.Now().Format(time.RFC3339Nano)
	http.SetCookie(w, cookie)

	b := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, b)
	handErr(err)
	sid := base64.URLEncoding.EncodeToString(b)

	cookie.HttpOnly = true
	cookie.Name = "sessionId"
	cookie.Value = sid
	http.SetCookie(w, cookie)

	fmt.Fprintln(w, cookie.Expires.Format(time.RFC3339Nano))
}

func getCookie(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Powered-By", "golang")

	j, _ := json.Marshal(r.Cookies())
	fmt.Fprintln(w, string(j))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://baidu.com", http.StatusSeeOther)
}

func main() {
	foo.Foo()
	foo2()
	timer()

	http.HandleFunc("/login", login)
	http.HandleFunc("/han", han)
	http.HandleFunc("/upload", lib.Upload)
	http.HandleFunc("/setCookie", setCookie)
	http.HandleFunc("/getCookie", getCookie)
	http.HandleFunc("/redirect", redirect)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Println("stoped.")

}

func foo2() {
	d := md5.Sum([]byte("helloworld"))
	str := fmt.Sprintf("%X", d)
	fmt.Println(str)

	str = strconv.FormatInt(0x12, 2)
	fmt.Println(str)

	b := make([]byte, 16)
	n, err := rand.Read(b)
	handErr(err)
	fmt.Printf("n:%d, b:%X\n", n, b)
}

func timer() {
	log.Println("timer")
	time.AfterFunc(10*time.Second, func() {
		timer()
	})
}
