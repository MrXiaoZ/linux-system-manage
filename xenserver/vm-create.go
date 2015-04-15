package main

import (
	"cimer-vms/models"
	"fmt"
	"flag"
	"github.com/hprose/hprose-go/hprose"
)

type VmsInstance struct {
    Vms_instance string
    Vms_uuid     string
    Vms_ip       string
}

type CreateVmsFromTpl struct {
//	Cvft func(string, int) (string, []string, map[int] interface{})
	Cvft func(string, int) (string, []string, map[int] VmsInstance)
}

func main() {
	client := hprose.NewClient("tcp4://1.1.1.1:4321/")
	var cvft *CreateVmsFromTpl
	client.UseService(&cvft)

	tpl := flag.String("tpl", "", "Please give me the vms_template!")
	num := flag.Int("num", 1 , "Please tell me how many vms do you want to create!")
//	creator := flag.String("creator", "root", "Who create the vms?")
//	pub := flag.Int("public", 1 , "Is the vms opened?")

	flag.Parse()

	if models.IsExistVmsTpl(*tpl) {
		rst, uuids, vmsinfo := cvft.Cvft(*tpl, *num)
		fmt.Println(rst)
		for _, v := range uuids {
			fmt.Println(v)
		}
		fmt.Println(vmsinfo)
		for _ , vms := range vmsinfo {
			p_v := &vms
			fmt.Println(p_v.Vms_ip)
			models.AddVmsInstance(*tpl , p_v.Vms_instance , p_v.Vms_uuid , p_v.Vms_ip)
		}
	} else {
		fmt.Println("虚拟机模板不存在!")
	}
}
