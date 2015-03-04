package main

import (
	"github.com/hprose/hprose-go/hprose"
	"flag"
	"fmt"
)

type Shutdown struct {
    Shutdown func()
}

type Reboot struct {
	Reboot func()
}

func main() {

	act := flag.String("act", "", "please give me the command")
    	flag.Parse()

	client := hprose.NewClient("http://172.16.1.231:8080/")

	var shutdown *Shutdown
    	client.UseService(&shutdown)
    	var reboot *Reboot
    	client.UseService(&reboot)

    if *act != "" {
        switch *act {
        case "shutdown":
            shutdown.Shutdown()
        case "reboot":
            reboot.Reboot()
        }
    }else{
        fmt.Println("Error......")    //debug
    }
}
