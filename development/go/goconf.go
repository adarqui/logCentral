package main

import (
	"github.com/postfix/goconf"
	"fmt"
)

func main() {

	conf, err := goconf.ReadConfigFile("go.conf")
	host, err := conf.GetString("default", "host")
	port, err := conf.GetInt("default", "port")

	misc, err := conf.GetString("misc", "hello")

	fmt.Println(host, port, misc, err)

	// http://blog.golang.org/go-maps-in-action
	// &{map[default:map[host:0.0.0.0 port:8000] misc:map[hello:world]]}
	fmt.Println(conf)
}
