// +build !windows

package teaagent

import (
	"github.com/iwind/TeaGo/files"
	"os/exec"
	"runtime"
	"strings"
)

// 获取系统发行版本信息
// 参考 https://whatsmyos.com/
func retrieveOsName() string {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("sw_vers", "-productVersion")
		data, err := cmd.CombinedOutput()
		if err != nil {
			return "Mac OS X"
		}
		return "Mac OS X " + strings.TrimSpace(string(data))
	}

	if runtime.GOOS == "linux" {
		{
			etcFile := files.NewFile("/etc/redhat-release")
			if etcFile.Exists() {
				s, err := etcFile.ReadAllString()
				if err == nil && len(s) > 0 {
					s = strings.Replace(s, " Release ", " ", -1)
					s = strings.Replace(s, " release", " ", -1)
					s = strings.Replace(s, " Linux ", " ", -1)
					s = strings.Replace(s, "(Core)", "", -1)

					if strings.HasPrefix(s, "Red Hat Enterprise Linux") {
						s = strings.Replace(s, "Red Hat Enterprise Linux", "RHEL", -1)
					}

					return strings.TrimSpace(s)
				}
			}
		}

		{
			etcFile := files.NewFile("/etc/issue")
			if etcFile.Exists() {
				s, err := etcFile.ReadAllString()
				if err == nil && len(s) > 0 {
					s = strings.Replace(s, "\\n", "", -1)
					s = strings.Replace(s, "\\l", "", -1)
					return strings.TrimSpace(s)
				}
			}
		}

		{
			etcFile := files.NewFile("/etc/issue.net")
			if etcFile.Exists() {
				s, err := etcFile.ReadAllString()
				if err == nil && len(s) > 0 {
					s = strings.Replace(s, "\\n", "", -1)
					s = strings.Replace(s, "\\l", "", -1)
					s = strings.Replace(s, " GNU/Linux ", " ", -1)
					return strings.TrimSpace(s)
				}
			}
		}
	}

	return ""
}
