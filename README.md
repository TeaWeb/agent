# Agent
配置*configs/agent.conf*：
~~~yaml
master: http://192.168.1.100:7777       # 主服务器地址
id: KwOe0dkxtyKHzMRC                    # 当前主机ID
key: ZtJUYdjO6enUwHwx9TyczhrqAHoO2FBv   # 当前主机密钥
~~~

测试和主服务器连接：
~~~bash
bin/teaweb-agent test
~~~

启动：
~~~bash
bin/teaweb-agent start
~~~

停止：
~~~bash
bin/teaweb-agent stop
~~~

重启：
~~~bash
bin/teaweb-agent restart
~~~

手动执行某个任务：
~~~bash
bin/teaweb-agent run [TASK ID]
~~~

手动执行某个监控项数据源程序：
~~~bash
bin/teaweb-agent run [ITEM ID]
~~~
