package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"time"
)

func TcpScanner(x string, ports chan int, opens chan int, closes chan int) {
	for p := range ports {
		address := fmt.Sprintf("%v:%d", x, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			closes <- p
		} else {
			conn.Close()
			opens <- p
		}
	}
}
func main() {
	start := time.Now()
	ports := make(chan int, 100)
	opens := make(chan int)
	closes := make(chan int)
	var openpoats []int
	var closeports []int
	var ipname string
	var portquantity int
	flag.StringVar(&ipname, "ip", "localhost", "ip地址")
	flag.IntVar(&portquantity, "portq", 120, "扫描的最大端口号")
	flag.Parse()
	for i := 0; i < cap(ports); i++ {
		go TcpScanner(ipname, ports, opens, closes)
	}
	go func(n int) {
		for i := 1; i < n; i++ {
			ports <- i
		}
	}(portquantity)

	for i := 1; i < portquantity; i++ {
		select {
		case port1 := <-opens:
			openpoats = append(openpoats, port1)
		case port2 := <-closes:
			closeports = append(closeports, port2)
		}
	}
	close(ports)
	close(opens)
	close(closes)
	sort.Ints(openpoats)
	sort.Ints(closeports)
	end := time.Since(start) / 1e9

	for _, p := range openpoats {
		fmt.Printf("open:%d\n", p)
	}
	fmt.Println("|||||||||分开了|||||||||")
	for _, p := range closeports {
		fmt.Printf("close:%d\n", p)
	}
	fmt.Printf("%ds #####\n", end)
}
