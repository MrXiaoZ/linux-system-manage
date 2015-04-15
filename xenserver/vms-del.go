package main

import (
	"cimer-vms/models"
	"fmt"
	"github.com/hprose/hprose-go/hprose"
	"flag"
)

type DelVms struct {
	DelVms func([]string) bool
}

func main() {
	client := hprose.NewClient("tcp4://1.1.1.1:4321/")
	var delvms *DelVms
	client.UseService(&delvms)

	method := flag.String("method", "tpl", "Please tell me the method!")
	tpl := flag.String("tpl", "", "Please tell me the vms template!")
	instance := flag.String("instance", "", "Please tell me the vms instance!")

	flag.Parse()

	if *method != "" {
		switch *method {
		case "tpl":
			vms_uuid := models.GetVmsUuidByTpl(*tpl)
			fmt.Println("Del vms_uuid <---> ",vms_uuid)
			if delvms.DelVms(vms_uuid) {
				models.DelVmsInstanceByTpl(*tpl)
			}
		case "instance":
			vms_uuid := models.GetVmsUuidByInstance(*instance)
			fmt.Println("Del vms_uuid <---> ",vms_uuid)
			if delvms.DelVms(vms_uuid) {
				models.DelVmsInstanceByInstance(*instance)
			}
		}
	} else {
		fmt.Println("Error......")
	}
}
