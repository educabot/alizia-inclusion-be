package main

import "github.com/educabot/alizia-inclusion-be/config"

func main() {
	app := NewApp(config.Load())
	defer app.Close()
	app.Run()
}
