package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	_ "github.com/zengming00/testgo/gojatest/lib"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
	_ "github.com/go-sql-driver/mysql"
)

func handErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("please input a file name")
		os.Exit(1)
	}
	filename := os.Args[1]
	f, err := os.Open(filename)
	handErr(err)
	datas, err := ioutil.ReadAll(f)
	handErr(err)
	str := string(datas)
	_ = str

	runtime := goja.New()

	registry := new(require.Registry) // this can be shared by multiple runtimes
	req := registry.Enable(runtime)
	console.Enable(runtime)

	time.AfterFunc(60*time.Second, func() {
		runtime.Interrupt("run code timeout, halt")
	})
	// 直接执行时，如果有错误，无法知道是在哪个文件报的错
	// v, err := runtime.RunString(str)
	v, err := req.Require(filename)

	if err != nil {
		if interruptErr, ok := err.(*goja.InterruptedError); ok {
			fmt.Println("InterruptedError:", interruptErr)
			return
		}
		panic(err)
	}
	if val := v.Export(); val != nil {
		// fmt.Println(val)
	}
}
