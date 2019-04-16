// +build windows

package teaagent

import (
	"golang.org/x/sys/windows/registry"
)

// 获取系统发行版本信息
func retrieveOsName() string {
	key, _ := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	productName, _, _ := key.GetStringValue("ProductName")
	key.Close()

	if len(productName) > 0 {
		return productName
	}
	return "Windows"
}
