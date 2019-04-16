// +build !windows

package teaagent

import (
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/maps"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// 获取系统发行版本信息
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
			osFile := files.NewFile("/etc/os-release")
			if osFile.Exists() {
				s, err := osFile.ReadAllString()
				if err == nil && len(s) > 0 {
					m := maps.Map{}
					for _, field := range strings.Split(s, "\n") {
						pieces := strings.SplitN(field, "=", 2)
						if len(pieces) != 2 {
							continue
						}
						m[pieces[0]] = strings.Trim(pieces[1], "\"")
					}
					name := m.GetString("NAME")
					version := m.GetString("VERSION_ID")
					if len(name) > 0 {
						return name + " " + version
					}
				}
			}
		}

		{
			etcFile := files.NewFile("/etc/redhat-release")
			if etcFile.Exists() {
				s, err := etcFile.ReadAllString()
				if err == nil && len(s) > 0 {
					return strings.TrimSpace(shortenOsName(s))
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

					return strings.TrimSpace(shortenOsName(s))
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
					return strings.TrimSpace(shortenOsName(s))
				}
			}
		}
	} else if runtime.GOOS == "freebsd" {
		cmd := exec.Command("uname", "-a")
		data, err := cmd.CombinedOutput()
		if err != nil {
			return "FreeBSD"
		}

		result := regexp.MustCompile("FreeBSD\\s+([0-9.]+)").FindStringSubmatch(string(data))
		if len(result) == 0 {
			return "FreeBSD"
		}
		return "FreeBSD " + result[1]
	}

	return ""
}

func shortenOsName(osName string) string {
	osName = strings.Replace(osName, "CentOS Linux", "CentOS", -1)
	osName = strings.Replace(osName, "Red Hat Enterprise Linux Server", "RHEL", -1)
	osName = strings.Replace(osName, "SUSE Linux Enterprise Server", "SLES", -1)
	return osName
}
