package main

import "fmt"

const prefix = "Hello, "

func Hello(name string, language string) string {
	if name == "" {
		return prefix + "world"
	}

	if language == "Spanish" {
		return "Hola, " + name
	}

	if language == "Russian" {
		return "Privet, " + name
	}

	return prefix + name
}

func main() {
    fmt.Println(Hello("world", ""))
}
