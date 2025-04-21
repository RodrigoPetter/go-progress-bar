# Go Progress Bar

A simple progress bar implementation for Go applications with support for nested progress bars, stages, and logging integration.

## Features

- ⏱️ Built-in timer showing elapsed time
- 📊 Percentage-based progress display
- 🌲 Support for nested progress bars (sub-tasks)
- 📝 Integrated logging support
- 🎯 Stage-based progress tracking
- 📐 Auto-adjusting bar width based on terminal size

## Installation

```bash
go get github.com/RodrigoPetter/go-progress-bar
```

## Usage

### Basic Progress Bar

```go
package main

import (
    "time"
    progressbar "github.com/RodrigoPetter/go-progress-bar"
)

func main() {
    // Create a new progress bar
    bar := progressbar.NewProgressBar("Main Task")
    
    // Start the render loop
    progressbar.StartRenderLoop(bar)

    // Update progress
    for i := 0; i <= 100; i++ {
        bar.Increment(1)
        time.Sleep(100 * time.Millisecond)
    }

    // Mark as complete
    bar.Finish()
}
```

This code will output the following bar:
```log
⠦[00:00:03] - Main Task - [35% #####          ]
```


### Nested Progress Bars

```go
func main() {
    mainBar := progressbar.NewProgressBar("Main Task")
    subBar1 := mainBar.NewSubBar("Subtask 1")
    subBar2 := mainBar.NewSubBar("Subtask 2")

    progressbar.StartRenderLoop(mainBar)

    // Update progress for different bars
    for i := 0; i <= 100; i++ {
        mainBar.Increment(1)
        subBar1.Increment(2)
        subBar2.Increment(1)
        time.Sleep(100 * time.Millisecond)
    }
}
```

This code will output the following bar:
```log
⠴[00:00:03] - Main Task - [30% ####           ]
    ⠴[00:00:03] - Subtask 1 - [60% ######     ]
    ⠴[00:00:03] - Subtask 2 - [30% ###        ]
```

### Stage-Based Progress

```go
func main() {
    // Create a progress bar with 5 stages
    bar := progressbar.NewProgressBar("Stage Progress", 
        progressbar.WithStages(5))
    
    progressbar.StartRenderLoop(bar)

    // Progress through stages
    for i := 0; i < 5; i++ {
        time.Sleep(1 * time.Second)
        bar.Increment(1)
    }
}
```

This code will output the following bar:
```log
⠦[00:00:03] Stage 3/5 - Stage Progress - [60% ##############          ]
```

### Logging Integration

This feature prevents log messages from breaking the progress bar's visual display. Without this integration, log messages would interrupt and corrupt the progress bar's appearance in the terminal.

```go
func main() {
    bar := progressbar.NewProgressBar("Task with Logs")
    
    // Configure the log writer
    // Ensures logs work harmoniously with the progress bar
    log.SetOutput(bar.GetLogWriter())  
    
    progressbar.StartRenderLoop(bar)

    // Logs will be displayed alongside the progress bar
    log.Println("Starting process...")
    for i := 0; i <= 100; i++ {
        bar.Increment(1)
        if i%20 == 0 {
            log.Printf("Completed %d%%", i)
        }
        time.Sleep(100 * time.Millisecond)
    }
    log.Println("Process completed!")
}
```

This code will output the following bar:
```log
2025/04/21 17:05:56 Starting process...
2025/04/21 17:05:56 Completed 0%
2025/04/21 17:05:58 Completed 20%
2025/04/21 17:06:00 Completed 40%
⠋[00:00:05] - Task with Logs - [55% ##################                ]
```

## Options

The progress bar can be customized using various options:

```go
bar := goprogressbar.NewProgressBar("Custom Bar",
    goprogressbar.WithMaxProgress(200),    // Set custom max progress
    goprogressbar.WithStages(5),           // Use stage-based progress
)
```