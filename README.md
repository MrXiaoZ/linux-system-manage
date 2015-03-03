# linux-system-manage
Golang + Hprose ，简单地实现了几个功能
关机、重启、IP地址设置、dhcp server 配置

使用：

假设编译好的client文件是：t1024，并已添加到系统可执行路径

关机：t1024 -act=”shutdown”

重启：t1024 -act=”reboot”

ip设置：t1024 -ip=”192.168.1.10″ -mask=”255.255.255.0″ -gw=”192.168.1.254″

dhcp配置：t1024 -subnet=”1.1.1.0″ -mask=”255.255.255.0″ -rt=”1.1.1.1″ -start=”1.1.1.110″ -end=”1.1.1.200″

 

TODO：

    Client请求验证（权限），限制可执行的client
    Server执行反馈
    代码优化
    加入更多功能

