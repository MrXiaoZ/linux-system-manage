package main

import (
	"github.com/hprose/hprose-go/hprose"
	"flag"
	"fmt"
)

type DhcpSet struct {
    
	DhcpSet func(string,string,string,string,string)(string,string,string,string,string)

}


func main() {

//	act := flag.String("act", "", "please give me the command")
	subnet := flag.String("subnet","","Please enter the subnet address!")
	mask := flag.String("mask","","Please enter the netmask!")
	start := flag.String("start","","Please enter the start ip address!")
	end := flag.String("end","","Please enter the end ip address!")
	routers := flag.String("rt","","Please enter the routers!")
    
	flag.Parse()

	client := hprose.NewClient("http://172.16.1.231:8080/")

	var dhcpset *DhcpSet
        client.UseService(&dhcpset)
	fmt.Println(dhcpset.DhcpSet("1.1.1.1","24","1.1.1.10","1.1.1.100","1.1.1.254"))
	fmt.Println(dhcpset.DhcpSet(*subnet,*mask,*start,*end,*routers))
}
