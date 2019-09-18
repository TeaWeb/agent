package teaagent

import "fmt"

// 打印帮助
func printHelp() {
	fmt.Print(`Usage:
~~~
bin/teaweb-agent						
   run in foreground

bin/teaweb-agent help 					
   show this help

bin/teaweb-agent start 					
   start agent in background

bin/teaweb-agent stop 					
   stop running agent

bin/teaweb-agent restart				
   restart the agent

bin/teaweb-agent status
   lookup agent status

bin/teaweb-agent run [TASK ID]		
   run task

bin/teaweb-agent run [ITEM ID]		
   run monitor item

bin/teaweb-agent init -master=[MASTER SERVER] -group=[GROUP KEY]
   register agent to master server and specified group

bin/teaweb-agent [-v|version]
   show agent version
~~~
`)
}
