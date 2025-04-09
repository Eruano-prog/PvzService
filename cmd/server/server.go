package main

import (
	"AvitoPvz/internal/app"
	"fmt"
	"os"
)

func main() {
	err := app.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
