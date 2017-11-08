package lib

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseMultipartForm(1024 * 1024)
		handErr(err)
		v := r.MultipartForm.Value["data"]
		v2 := r.Form.Get("data")
		log.Println("v:", v)
		log.Println("v2:", v2)

		file, header, err := r.FormFile("file")
		handErr(err)
		defer file.Close()

		fmt.Fprintf(w, "%#v", header.Header)

		f, err := os.OpenFile("./upload/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		handErr(err)
		defer f.Close()
		_, err = io.Copy(f, file)
		handErr(err)

	} else {
		io.WriteString(w, "err")
	}
}
