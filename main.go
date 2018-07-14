package main

import (
	"chip8/cpu"
	"fmt"
)

func main() {
	fmt.Println("hello")
	c := cpu.NewCPU()
	fmt.Println(c)
}
