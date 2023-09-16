package main

import (
	"time"
)

// ServiceUsage tracks the number of API requests made for each service.
var serviceUsage = make(map[string]int)

func main() {

	core := NewCore()
	ui := NewUI(core)
	ticker := time.NewTicker(time.Millisecond * 500)
	running := true

	// Start the UI in a separate goroutine
	go func() {
		if err := ui.App.Run(); err != nil {
			panic(err)
		}
		running = false
	}()

	// Main loop
	for range ticker.C {
		if !running {
			break
		}
		// Update the core
		core.Update()
		// Update the UI
		ui.Draw(core)
	}

}
