package main

import (
	"bufio"
	"fmt" //for debug
	"github.com/hprose/hprose-go/hprose"
	"io"
	//	"net"
	"cimer-vms/config"
	"cimer-vms/models"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type VmsInstance struct {
	Vms_instance string
	Vms_uuid     string
	Vms_ip       string
	Vms_creator  string
}
func CreateVmsFromTpl(tpl string, num int, creator string) (rst string, uuid []string, vmsinfo map[int]VmsInstance) {

	var vms_uuid []string
	var vms_ip string = "Null"
	var rst_c string = "虚拟机实例创建成功!"
	ip_chan := make(chan string, 10)
	sr_ssd_uuid := config.GetSsdUuid()

	//	var vmsInfo map[int] VmsInstance
	//	vmsInfo := make(map[int]*VmsInstance)
	vmsInfo := make(map[int]VmsInstance)
	for i := 0; i < num; i++ {
		ts := strconv.FormatInt(time.Now().UnixNano(), 10)
		//构建执行的shell
		c_create := "xe vm-install template=\"" + tpl + "\" new-name-label=\"" + tpl + "_" + ts[9:] + "\"" + " sr-uuid=" + sr_ssd_uuid //strconv.Itoa(i): int to string
		c_start := "xe vm-start name-label=\"" + tpl + "_" + ts[9:] + "\"" + " --multiple"
		//
		vms_instance_name := tpl + "_" + ts[9:]

		cmd_create := exec.Command("/bin/bash", "-c", c_create)

		//指针变量
		p_vms_uuid := &vms_uuid
		p_rst_c := &rst_c
		p_vms_ip := &vms_ip
		//		p_vms_info := &vmsInfo
		//		p_vms_info = make(map[int]*VmsInstance)

		vms_uuid_inner, err := cmd_create.CombinedOutput() // 创建虚拟机并获得其UUID
		vms_uuid_s := strings.TrimSpace(string(vms_uuid_inner))
		//	vms_instance_name := tpl + "_" + strings.Split(vms_uuid_s ,"-")[0]
		if err != nil {
			*p_rst_c = "虚拟机创建失败!" //需要改进，返回具体的错误信息
			*p_vms_uuid = append(*p_vms_uuid, "No Vms_uuid!")
		} else {
			*p_vms_uuid = append(*p_vms_uuid, vms_uuid_s)

			models.AddVmsInstance(tpl, vms_instance_name, vms_uuid_s, "", creator)

		//	time.Sleep(1 * time.Second)
			cmd_start := exec.Command("/bin/bash", "-c", c_start)
			cmd_start.Run()
			//			time.Sleep(5 * time.Second)
			c_get_ip := "xe vm-list params=networks uuid=" + vms_uuid_s + " | grep \"networks\" | awk -F'[;|:]' '{print $3}'"

			go func() {
				for k := 0; k < 150; k++ {
					cmd_get_ip := exec.Command("/bin/bash", "-c", c_get_ip)
					vms_ip_inner, err_get_ip := cmd_get_ip.CombinedOutput()
					if len(strings.TrimSpace(string(vms_ip_inner))) != 0 && err_get_ip == nil {
						//		if strings.TrimSpace(string(vms_ip_inner)) != "" && err_get_ip == nil {
						*p_vms_ip = strings.TrimSpace(string(vms_ip_inner))
						models.UpdateVmsIp(vms_uuid_s, vms_ip)
						break
					}

					time.Sleep(2 * time.Second)
				}
				ip_chan <- vms_ip
			}()
			vmsInfo[i] = VmsInstance{vms_instance_name, vms_uuid_s, vms_ip, creator}
		}
		fmt.Println("Debug <--->", vmsInfo)
	}
	for v := range ip_chan {
		if len(ip_chan) == 0 {
			break
		} else {
			fmt.Println(v)
		}

	}
	//	return "虚拟机实例创建成功!", string(vms_uuid)
//	fmt.Println("Debug:", vms_uuid)
	return rst_c, vms_uuid, vmsInfo
}

func main() {
	server := hprose.NewTcpServer("tcp4://0.0.0.0:4321/")
	server.AddFunction("cvft", CreateVmsFromTpl)

	server.Start()
	for {
		time.Sleep(1 * time.Second )
	}
	server.Stop()
}
