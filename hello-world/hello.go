package main

import "fmt"

const prefix = "Hello, "
func Hello(name string) string {
	return prefix + name
}

func main() {
    fmt.Println(Hello("world"))
}
