package main

import (
	"flag"
	"log"
	"net"
)

var (
	mac string
	i   string
)

func init() {
	flag.StringVar(&mac, "mac", "00:aa:bb:cc:dd:ee", "请输入MAC地址")
	flag.StringVar(&i, "i", "eth0", "请输入网卡名")
}

func main() {
	flag.Parse()

	//var errBuffer bytes.Buffer
	//var err error

	//down := exec.Command("ifconfig", i, "down")
	//down.Stderr = &errBuffer
	//down.Stdout = os.Stdout
	//err = down.Run()
	//if err != nil {
	//	log.Print(errBuffer.String())
	//	log.Println(err.Error())
	//	return
	//}

	//change := exec.Command("ifconfig", i, "hw", "ether", mac)
	//change.Stderr = &errBuffer
	//change.Stdout = os.Stdout
	//err = change.Run()
	//if err != nil {
	//	log.Print(errBuffer.String())
	//	log.Println(err.Error())
	//	return
	//}

	//up := exec.Command("ifconfig", i, "up")
	//up.Stderr = os.Stderr
	//up.Stdout = os.Stdout
	//err = up.Run()
	//if err != nil {
	//	log.Print(errBuffer.String())
	//	log.Println(err.Error())
	//	return
	//}

	interfaces, _ := net.Interfaces()
	for _, v := range interfaces {
		log.Printf("网卡 [%s], MAC [%s]\n", v.Name, v.HardwareAddr.String())

	}
}
