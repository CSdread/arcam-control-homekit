package main

import (
	"fmt"
	"github.com/brutella/hc/accessory"
)

func main() {
	fmt.Println("hello")
	info := accessory.Info{
		Name: "AVR11",
	}

	fmt.Println(info)
}
