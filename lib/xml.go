package lib

import (
	"encoding/xml"
	"fmt"
)

const xml1 = `<?xml version="1.0" encoding="utf-8"?>
<servers version="1">
	<server>
		<serverName>Shanghai_VPN</serverName>
		<serverIP>127.0.0.1</serverIP>
	</server>
	<server>
		<serverName>Beijing_VPN</serverName>
		<serverIP>127.0.0.2</serverIP>
	</server>
</servers>`

type Recurlyservers struct {
	XMLName     xml.Name `xml:"servers"`
	Version     string   `xml:"version,attr"`
	Svs         []server `xml:"server"`
	Description string   `xml:",innerxml"`
}

type server struct {
	XMLName    xml.Name `xml:"server"`
	ServerName string   `xml:"serverName"`
	ServerIP   string   `xml:"serverIP"`
}

type user struct {
	XMLName xml.Name `xml:"xml"`
	Name    string   `xml:"user>name"`
	Age     int8     `xml:"user>age"`
}

func init() {
	test1()
}

func test1() {
	fmt.Println(xml1)
	v := &Recurlyservers{}
	err := xml.Unmarshal([]byte(xml1), v)
	handErr(err)
	fmt.Println(v.Svs[0].ServerName)

	user := user{Name: "zengming", Age: 24}
	data, err := xml.MarshalIndent(user, "", "  ")
	handErr(err)
	fmt.Println("user xml: ")
	fmt.Println(string(data))
}
