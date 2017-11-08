package lib

import (
	"encoding/json"
	"fmt"
)

type Serverslice struct {
	Servers []struct {
		ServerName string
		ServerIP   string
	}
}

/*
bool, for JSON booleans
float64, for JSON numbers
string, for JSON strings
[]interface{}, for JSON arrays
map[string]interface{}, for JSON objects
nil for JSON null
*/
func init() {
	var s Serverslice
	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(str), &s)
	fmt.Println("json.go:")
	fmt.Println(s.Servers[0].ServerName)

	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":[null,123,3.14, true, [1,3], {"a":true},"Gomez","Morticia"]}`)
	var f map[string]interface{}
	json.Unmarshal(b, &f)
	fmt.Printf("1. %v\n", f)
	if v, ok := f["Parents"].([]interface{}); ok {
		for _, tmp := range v {
			switch vv := tmp.(type) {
			case string:
				fmt.Println("string:", vv)
			case float64:
				fmt.Println("float64:", vv)
			case map[string]interface{}:
				type obj struct {
					a bool
				}
				// o, ok := vv.(obj)
				// fmt.Println("oooooo:", o.a)
				fmt.Printf("%T: %#[1]v\n", vv)
			default:
				fmt.Printf("%T: %#[1]v\n", vv)
			}
		}
	}
	fmt.Println("==================================================")
}
