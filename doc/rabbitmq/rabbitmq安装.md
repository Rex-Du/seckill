apt-get install rabbitmq-server  #安装成功自动启动

service rabbitmq-server start    # 启动
service rabbitmq-server stop     # 停止
service rabbitmq-server restart  # 重启 

rabbitmq-plugins list # 查看有哪些可用的插件和哪些已经启用的插件
rabbitmq-plugins enable rabbitmq_management   # 启用插件

此时，应该可以通过 http://localhost:15672 查看，使用默认账户guest/guest 登录。

如果有报错，看看是不是执行命令时少了sudo

远程不能访问15672时：
第一步：添加 mq 用户并设置密码
rabbitmqctl add_user rexdu rootroot
第二步：添加 mq 用户为administrator角色
rabbitmqctl set_user_tags rexdu administrator
rabbitmqctl list_users
第三步：设置 mq 用户的权限，指定允许访问的vhost以及write/read
rabbitmqctl set_permissions -p "/" rexdu ".*" ".*" ".*"
第四步：查看vhost（/）允许哪些用户访问
rabbitmqctl list_permissions -p /
第五步：配置允许远程访问的用户，rabbitmq的guest用户默认不允许远程主机访问。

在windows 下的 rabbitmq安装文件下的etc文件下的配置文件添加以下

 [
    {rabbit, [{tcp_listeners, [5672]}, {loopback_users, ["账户名"]}]}
    ].

 

管理账户命令如下：

复制代码
# 在rabbitmq的内部数据库添加用户；
add_user <username> <password>  
 
# 删除一个用户；
delete_user <username>  
 
# 改变用户密码（也是改变web管理登陆密码）；
change_password <username> <newpassword>  
 
# 清除用户的密码，该用户将不能使用密码登陆，但是可以通过SASL登陆如果配置了SASL认证；
clear_password <username>
 
# 设置用户tags；
set_user_tags <username> <tag> ...
 
# 列出用户；
list_users  
 
# 创建一个vhosts；
add_vhost <vhostpath>  
 
# 删除一个vhosts；
delete_vhost <vhostpath>  
 
# 列出vhosts；
list_vhosts [<vhostinfoitem> ...]  
 
# 针对一个vhosts给用户赋予相关权限；
set_permissions [-p <vhostpath>] <user> <conf> <write> <read>  
 
# 清除一个用户对vhosts的权限；
clear_permissions [-p <vhostpath>] <username>  
 
# 列出哪些用户可以访问该vhosts；
list_permissions [-p <vhostpath>]  
 
# 列出该用户的访问权限；
list_user_permissions <username>  
 
set_parameter [-p <vhostpath>] <component_name> <name> <value>
clear_parameter [-p <vhostpath>] <component_name> <key>
list_parameters [-p <vhostpath>]