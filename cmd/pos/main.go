package main

import (
	"os"

	"pos/internal/infrastructure/console"
	"pos/internal/infrastructure/container"
	"pos/internal/infrastructure/export"
)

func main() {
	app := container.New()

	exporter := export.NewCSVExporter(app.App)

	ui := console.New(
		app.App,
		exporter,
		os.Stdin,
		os.Stdout,
	)

	ui.Run()
}
