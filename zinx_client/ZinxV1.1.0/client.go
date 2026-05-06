package main

import (
	"fmt"
	"net"
	"time"
)

// 客户端模拟
func main() {

	fmt.Println("client start...")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	for {
		_, err := conn.Write([]byte("hello zinx v1.1.0..."))
		if err != nil {
			fmt.Println(err)
		}

		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("server call back: %s, cnt: %d\n", buf[:n], n)

		// cpu 阻塞
		time.Sleep(1 * time.Second)
	}
}
