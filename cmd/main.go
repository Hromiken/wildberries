package main

import (
	"order-notification/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
