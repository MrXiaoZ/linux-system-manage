package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
//	"strings"
)

type VmsInstance struct {
	Id           int64  `orm:"pk"`
	Vms_template string `orm:"size(100)"`
	Vms_instance string `orm:"size(100)"`
	Vms_ip       string `orm:"size(20)"`
	Vms_uuid     string `orm:"size(50)"`
	Public       int
	Vms_creator  string `orm:"size(50)"`
}

type VmsTemplate struct {
	Id           int64  `orm:"pk"`
	Vms_template string `orm:"size(100)"`
	Desc         string `orm:"size(100)"`
}

func init() {
    orm.RegisterDriver("mysql", orm.DR_MySQL)
    orm.RegisterModel(new(VmsInstance), new(VmsTemplate))
	orm.RegisterDataBase("default", "mysql", "root:roottoor@/isdp?charset=utf8", 30)
}

func GetIdUnix() (id int64) {
    t := time.Now().UnixNano()
    return t
}

func AddVmsTpl (tpl , desc string){
	orm.Debug = true
    o := orm.NewOrm()
    o.Using("default")
	var tpls VmsTemplate
	exist := o.QueryTable("vms_template").Filter("Vms_template", tpl).Exist()
	if !exist == true {
		tpls.Id = GetIdUnix()
		tpls.Vms_template = tpl
		tpls.Desc = desc
		fmt.Println(o.Insert(&tpls))
	}
}

func IsExistVmsTpl (tpl string)(bool){
	orm.Debug = true
    o := orm.NewOrm()
    o.Using("default")
//	var tpls VmsTemplate
    exist := o.QueryTable("vms_template").Filter("Vms_template", tpl).Exist()
    if !exist == true {
		return false
	}else{
		return true
	}
}


func AddVmsInstance ( tpl , instance , uuid , vms_ip string){
	orm.Debug = true
    o := orm.NewOrm()
    o.Using("default")
	var vms_instance VmsInstance
	vms_instance.Id = GetIdUnix()
	vms_instance.Vms_template = tpl
	vms_instance.Vms_instance = instance
	vms_instance.Vms_ip = vms_ip
	vms_instance.Vms_uuid = uuid
	fmt.Println(o.Insert(&vms_instance))
}

func DelVmsInstanceByTpl(tpl string){
	orm.Debug = true
    o := orm.NewOrm()
    o.Using("default")
//	var vms_instance VmsInstance
//	o.Delete(&VmsInstance{Vms_template: tpl})
	o.QueryTable("vms_instance").Filter("Vms_template", tpl).Delete()
}

func DelVmsInstanceByInstance(instance string){
	orm.Debug = true
    o := orm.NewOrm()
    o.Using("default")
//	var vms_instance VmsInstance
//	o.Delete(&VmsInstance{Vms_instance: instance})
	o.QueryTable("vms_instance").Filter("Vms_instance", instance).Delete()
}

func GetVmsUuidByTpl (tpl string)( uuid []string) {
	orm.Debug = true
    o := orm.NewOrm()
    o.Using("default")
//	var vms_instance VmsInstance
	var vms_uuid []string
	var maps []orm.Params
    num, err := o.QueryTable("vms_instance").Filter("Vms_template", tpl).Values(&maps)
    if err == nil {
        fmt.Println("Result numbers:", num)
        for _, m := range maps {
            //fmt.Println(m["Mgr_uuid"], m["Mgr_login"])    取特定字段的值
            fmt.Println(m)
			if str , ok := m["Vms_uuid"].(string); ok {
				vms_uuid = append(vms_uuid , str)
			}
        }
    }
	return vms_uuid
}

func GetVmsUuidByInstance(instance string)(uuid []string){
	orm.Debug = true
    o := orm.NewOrm()
    o.Using("default")
//	var vms_instance VmsInstance
	var vms_uuid []string
	var maps []orm.Params
    num, err := o.QueryTable("vms_instance").Filter("Vms_instance", instance).Values(&maps)
    if err == nil {
        fmt.Println("Result numbers:", num)
        for _, m := range maps {
            //fmt.Println(m["Mgr_uuid"], m["Mgr_login"])    取特定字段的值
            fmt.Println(m)
			if str , ok := m["Vms_uuid"].(string); ok {
				vms_uuid = append(vms_uuid , str)
			}
        }
    }
	return vms_uuid
}

func UpdateVmsIp( uuid , vms_ip string){
	orm.Debug = true
    o := orm.NewOrm()
    o.Using("default")
	var vms_instance VmsInstance
	err := o.QueryTable("vms_instance").Filter("Vms_uuid", uuid).One(&vms_instance)
	if err == nil {
		vms_instance.Vms_ip = vms_ip
		o.Update(&vms_instance)
	}
}
