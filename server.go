package main

import (
	"fmt" //for debug
	"github.com/hprose/hprose-go/hprose"
	"net/http"
	"os"
	"os/exec"
)

func Shutdown() {
	cmd := exec.Command("shutdown", "-h", "0")
	cmd.Run()
}

func Reboot() {
	cmd := exec.Command("reboot")
	cmd.Run()
}

func IpSet(ip string, mask string, gw string) (string, string, string) {
	userFile := "/etc/sysconfig/network-scripts/ifcfg-eth1"
	fout, err := os.Create(userFile)
	if err != nil {
		fmt.Println(userFile, err)
		//        return err
	}
	defer fout.Close()
	fout.WriteString("DEVICE=eth1\nHWADDR=0E:70:B1:11:18:D8\nTYPE=Ethernet\nONBOOT=yes\nNM_CONTROLLED=yes\n")
	fout.WriteString("BOOTPROTO=static\n")
	fout.WriteString("IPADDR=" + ip + "\n")
	fout.WriteString("NETMASK=" + mask + "\n")
	fout.WriteString("GATEWAY=" + gw + "\n")
	cmd := exec.Command("service", "network", "restart")
	cmd.Run()
	//	fmt.Println(ip,mask,gw)		//debug
	return ip, mask, gw
}

func DhcpSet(subnet string, mask string, start string, end string, routers string) (string, string, string, string, string) {
	userFile := "/etc/dhcp/dhcpd.conf"
	fout, err := os.Create(userFile)
	if err != nil {
		fmt.Println(userFile, err)
		//        return err
	}
	defer fout.Close()
	fout.WriteString("ddns-update-style none;\n")
	fout.WriteString("subnet " + subnet + " netmask " + mask + " {\n")
	fout.WriteString("range " + start + " " + end + ";\n")
	fout.WriteString("option routers " + routers + ";\n")
	fout.WriteString("default-lease-time 600;\n")
	fout.WriteString("max-lease-time 7200;\n}")
	cmd := exec.Command("service", "dhcpd", "restart")
	cmd.Run()
	//	fmt.Println(subnet,mask,start,end,routers)		//debug
	return subnet, mask, start, end, routers
}

func main() {
	service := hprose.NewHttpService()
	service.DebugEnabled = true
	service.AddFunction("reboot", Reboot)
	service.AddFunction("shutdown", Shutdown)
	service.AddFunction("ipset", IpSet)
	service.AddFunction("dhcpset", DhcpSet)
	http.ListenAndServe(":8080", service)
}
