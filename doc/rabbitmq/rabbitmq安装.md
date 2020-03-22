apt-get install rabbitmq-server  #安装成功自动启动

service rabbitmq-server start    # 启动
service rabbitmq-server stop     # 停止
service rabbitmq-server restart  # 重启 

rabbitmq-plugins list # 查看有哪些可用的插件和哪些已经启用的插件
rabbitmq-plugins enable rabbitmq_management   # 启用插件

此时，应该可以通过 http://localhost:15672 查看，使用默认账户guest/guest 登录。

如果有报错，看看是不是执行命令时少了sudo