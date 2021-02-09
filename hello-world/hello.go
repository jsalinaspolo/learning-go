package hello_world

import "fmt"

func Hello(name string, language string) string {
	prefix := "Hello, "

	if name == "" {
		return prefix + "world"
	}

	switch language {
	case "russian":
		prefix = "Privet, "
	case "spanish":
		prefix = "Hola, "
	}

	return prefix + name
}

func main() {
	fmt.Println(Hello("world", ""))
}
