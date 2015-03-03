package main

import (
	"github.com/hprose/hprose-go/hprose"
	"flag"
	"fmt"
)

type IpSet struct {
    
	IpSet func(string,string,string)(string,string,string)

}


func main() {

//	act := flag.String("act", "", "please give me the command")
	ip := flag.String("ip","","Please enter the ip address!")
	mask := flag.String("mask","","Please enter the netmask!")
	gateway := flag.String("gw","","Please enter the gateway!")
    
	flag.Parse()

	client := hprose.NewClient("http://172.16.1.231:8080/")

	var ipset *IpSet
        client.UseService(&ipset)
	fmt.Println(ipset.IpSet("1.1.1.1","24","1.1.1.254"))
	fmt.Println(ipset.IpSet(*ip,*mask,*gateway))
}
