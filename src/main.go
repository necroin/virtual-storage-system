package main

import (
	"fmt"
	"vss/src/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}
