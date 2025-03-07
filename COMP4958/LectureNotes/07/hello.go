package main

import (
	"fmt"
	"os"
)

func greet(lang string) {
		switch lang {
		case "en":
			fmt.Println("Hello!")
		case "ja":
			fmt.Println("こんにちは!")
		case "fr":
			fmt.Println("Bonjour!")
		default:
			fmt.Println("Hello!")
		}
}

func main() {
	if len(os.Args) == 1 {
		greet("en")
	} else {
		greet(os.Args[1])
	}
}