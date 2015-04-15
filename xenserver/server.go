package main

import (
	"bufio"
	"fmt" //for debug
	"github.com/hprose/hprose-go/hprose"
	"io"
//	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	//	"cimer-vms/models"
)

type VmsInstance struct {
	Vms_instance string
	Vms_uuid     string
	Vms_ip       string
}

func ReadFileToSlice() (tpls []string) {
	fl, err := os.Open("/tmp/vms-templates")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fl.Close()

	vmstpl := []string{}

	br := bufio.NewReader(fl)
	for {
		a, _, c := br.ReadLine()
		//      a = strings.TrimSpace(a)
		b := strings.TrimSpace(string(a))
		if c == io.EOF {
			break
		}
		//      fmt.Println(strings.TrimSpace(string(a)))
		//		fmt.Println(b)
		vmstpl = append(vmstpl, b)
		//      fmt.Println(a)
	}
	return vmstpl
}

func VmsTemplateList() (tpls []string) {
	cmd := exec.Command("/bin/bash", "-c", `xe template-list | grep "_tpl" | awk -F': ' '{print $2}' > /tmp/vms-templates`)
	cmd.Run()
	vmstpl := ReadFileToSlice()
	for _, v := range vmstpl {
		fmt.Println(v)
		//		models.AddVmsTpl(v,"")
	}
	return vmstpl
}

func CreateVmsFromTpl(tpl string, num int) (rst string, uuid []string, vmsinfo map[int]VmsInstance) {
	var vms_uuid []string
	var vms_ip string
	var rst_c string = "虚拟机实例创建成功!"
	//	var vmsInfo map[int] VmsInstance
	//	vmsInfo := make(map[int]*VmsInstance)
	vmsInfo := make(map[int]VmsInstance)
	for i := 0; i < num; i++ {
		//构建执行的shell
		c_create := "xe vm-install template=\"" + tpl + "\" new-name-label=\"" + tpl + "_" + strconv.Itoa(i) + "\"" //strconv.Itoa(i): int to string
		c_start := "xe vm-start name-label=\"" + tpl + "_" + strconv.Itoa(i) + "\"" + " --multiple"
		//
		vms_instance_name := tpl + "_" + strconv.Itoa(i)

		cmd_create := exec.Command("/bin/bash", "-c", c_create)

		//指针变量
		p_vms_uuid := &vms_uuid
		p_rst_c := &rst_c
		p_vms_ip := &vms_ip
		//		p_vms_info := &vmsInfo
		//		p_vms_info = make(map[int]*VmsInstance)

		vms_uuid_inner, err := cmd_create.CombinedOutput()
		vms_uuid_s := strings.TrimSpace(string(vms_uuid_inner))
		if err != nil {
			*p_rst_c = "虚拟机创建失败!"
			*p_vms_uuid = append(*p_vms_uuid, "No Vms_uuid!")
		} else {
			*p_vms_uuid = append(*p_vms_uuid, vms_uuid_s)

			time.Sleep(5 * time.Second)
			cmd_start := exec.Command("/bin/bash", "-c", c_start)
			cmd_start.Run()
			time.Sleep(15 * time.Second)
			c_get_ip := "xe vm-list params=networks uuid=" + vms_uuid_s + " | grep \"networks\" | awk -F'[;|:]' '{print $3}'"

			for {
				cmd_get_ip := exec.Command("/bin/bash", "-c", c_get_ip)
				vms_ip_inner, err_get_ip := cmd_get_ip.CombinedOutput()
				if len(strings.TrimSpace(string(vms_ip_inner))) != 0 && err_get_ip == nil {
					//		if strings.TrimSpace(string(vms_ip_inner)) != "" && err_get_ip == nil {
					*p_vms_ip = strings.TrimSpace(string(vms_ip_inner))
					break
				}

				time.Sleep(2 * time.Second)
			}
			vmsInfo[i] = VmsInstance{vms_instance_name, vms_uuid_s, vms_ip}
		}
		fmt.Println("Debug <--->", vmsInfo)

		//		cmd_create.Run()
		/*
			time.Sleep(5 * time.Second)
			cmd_start := exec.Command("/bin/bash", "-c", c_start)
			cmd_start.Run()
			time.Sleep(15 * time.Second)
		*/
	}
	//	return "虚拟机实例创建成功!", string(vms_uuid)
	fmt.Println("Debug:", vms_uuid)
	return rst_c, vms_uuid, vmsInfo
}

func DelVmsInstanceByUuid(uuid []string) (result bool) {
	rst := true
	for _, v := range uuid {
		vms_uuid := v
		del := "xe vm-uninstall uuid=" + vms_uuid + " --force"
		cmd_del := exec.Command("/bin/bash", "-c", del)
		err := cmd_del.Run()
		if err != nil {
			rst = false
		}
	}
	return rst
}

func GetVmsIp(vms_uuid string) (rst string) {
	var vms_ip string
	c_get_ip := "xe vm-list params=networks uuid=" + vms_uuid + " | grep \"networks\" | awk -F'[;|:]' '{print $3}'"

	for i := 0; i < 15; i++{
		cmd_get_ip := exec.Command("/bin/bash", "-c", c_get_ip)
		vms_ip_inner, err_get_ip := cmd_get_ip.CombinedOutput()
		if len(strings.TrimSpace(string(vms_ip_inner))) != 0 && err_get_ip == nil {
			//      if strings.TrimSpace(string(vms_ip_inner)) != "" && err_get_ip == nil {
			vms_ip = strings.TrimSpace(string(vms_ip_inner))
			break
		}
		time.Sleep(2 * time.Second)
	}
	return vms_ip
}

func main() {
	server := hprose.NewTcpServer("tcp4://0.0.0.0:4321/")
	server.AddFunction("vmstl", VmsTemplateList)
	server.AddFunction("cvft", CreateVmsFromTpl)
	server.AddFunction("delvms", DelVmsInstanceByUuid)
	server.AddFunction("getvmsip", GetVmsIp)
	server.Start()
	for {
		time.Sleep(1)
	}
	server.Stop()
}
