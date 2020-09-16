package main

import "fmt"

const prefix = "Hello, "

func Hello(name string) string {
	if name == "" {
		return prefix + "world"
	} 

	return prefix + name
}

func main() {
    fmt.Println(Hello("world"))
}
