package main

import (
	"cimer-vms/models"
	"fmt"
	"github.com/hprose/hprose-go/hprose"
	"flag"
)

type VmsTemplateList struct {
	Vmstl func() []string
}

type GetVmsIp struct {
	GetVmsIp func (string)string
}

func main() {
	client := hprose.NewClient("tcp4://1.1.1.1:4321/")
	var vmstl *VmsTemplateList
	client.UseService(&vmstl)
	var getvmsip *GetVmsIp
	client.UseService(&getvmsip)
	
	act := flag.String("act", "", "Please tell me what will you do!")
	instance := flag.String("instance", "", "Please tell me the vms instance!")

	flag.Parse()

	if *act != "" {
		switch *act {
		case "addtpl":
			tpls := vmstl.Vmstl()
			for _, v := range tpls {
				fmt.Println(v)	//	For debug
				models.AddVmsTpl(v, "")
			}
		case "getip":
			for _ , v := range models.GetVmsUuidByInstance(*instance) {
				models.UpdateVmsIp(v,getvmsip.GetVmsIp(v))
			}
		}
	}
}
