package main

import (
	"log"
	"time"

	progressbar "github.com/RodrigoPetter/go-progress-bar"
)

func main() {
	mainBar := progressbar.NewProgressBar("Main Task")
	subBar1 := mainBar.NewSubBar("Subtask 1", progressbar.WithMaxProgress(10))
	subBar2 := mainBar.NewSubBar("Subtask 2")

	// Configure the log writer
	log.SetOutput(mainBar.GetLogWriter())

	progressbar.StartRenderLoop(mainBar)

	// Logs will be displayed alongside the progress bar
	log.Println("Starting process...")
	for i := 0; i <= 100; i++ {

		mainBar.Increment(1)
		subBar2.Increment(1)
		time.Sleep(100 * time.Millisecond)

		if i == 10 {
			subBar1.Finish()
			log.Printf("Subtask 1 Completed %d%%", i)
		} else if i < 10 {
			subBar1.Increment(1)
		}
		time.Sleep(1 * time.Second)
	}
	log.Println("Process completed!")
}
