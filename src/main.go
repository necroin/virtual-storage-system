package main

import (
	"vss/src/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
