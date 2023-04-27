package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	// 在 cmd `ipconfig /all` 命令查找网卡名字 , 取"描述"那一栏
	intel = "Intel(R) Ethernet Connection (11) I219-V"
)

// 使用该程序之前还需要先安装 rsrc ` go install github.com/akavel/rsrc@latest` 用于打包为可以使用管理员方式运行的程序
// 该程序需要使用管理员方式运行, 如果在 goland 中调试, 需要使用管理员方式打开 goland
//
//go:generate rsrc -manifest ../admin.manifest -o admin.syso
func main() {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Class\{4D36E972-E325-11CE-BFC1-08002BE10318}`, registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	defer key.Close()

	names, err := key.ReadSubKeyNames(0)
	for i, name := range names {
		if name == "Properties" {
			break
		}
		key2, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Class\{4D36E972-E325-11CE-BFC1-08002BE10318}\`+name, registry.ALL_ACCESS)
		if err != nil {
			panic(err)
		}

		val, _, err := key2.GetStringValue("DriverDesc")
		if err != nil && err != registry.ErrNotExist {
			panic(err)
		}

		if val == intel {
			randMac := GenerateRandMAC()
			log.Println("randMac = " + randMac)
			err := key2.SetStringValue("NetworkAddress", randMac)
			if err != nil {
				panic(err)
			}

			command := exec.Command("cmd", "/c", fmt.Sprintf("wmic path win32_networkadapter where index=%d call disable", i))
			output, err := command.CombinedOutput()
			if err != nil {
				panic(err)
			}
			reader := transform.NewReader(bytes.NewReader(output), simplifiedchinese.GBK.NewDecoder())
			all, _ := io.ReadAll(reader)
			log.Println(string(all))

			command2 := exec.Command("cmd", "/c", fmt.Sprintf("wmic path win32_networkadapter where index=%d call enable", i))
			output, err = command2.CombinedOutput()
			if err != nil {
				panic(err)
			}
			reader = transform.NewReader(bytes.NewReader(output), simplifiedchinese.GBK.NewDecoder())
			all, _ = io.ReadAll(reader)
			log.Println(string(all))
			key2.Close()

			return
		}
		key2.Close()
	}
}

// GenerateRandMAC generates a random unicast and locally administered MAC address.
func GenerateRandMAC() string {
	randBytes := make([]byte, 6)
	_, _ = rand.Read(randBytes)
	// Set locally administered addresses bit and reset multicast bit
	randBytes[0] = (randBytes[0] | 0x02) & 0xfe
	mac := net.HardwareAddr(randBytes).String()
	return strings.ReplaceAll(mac, ":", "")
}
