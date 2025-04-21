package progressbar

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

var spinnerFrames = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}

type progressBarLogsWriter struct {
	buffer chan string
}

func (pbl *progressBarLogsWriter) Write(p []byte) (n int, err error) {
	pbl.buffer <- string(p)
	return len(p), nil
}

func (pbl *progressBarLogsWriter) flush() {
	for {
		select {
		case msg := <-pbl.buffer:
			fmt.Print(msg)
		default:
			return
		}
	}
}

func (pb *ProgressBar) GetLogWriter() io.Writer {
	return &pb.logWriter
}

func StartRenderLoop(pb *ProgressBar) {
	ticker := time.NewTicker(500 * time.Millisecond)

	var lastLoopPrintedRows int

	go func() {
		for {
			select {
			case <-ticker.C:

				rowsToPrint := pb.render(len(pb.Name))
				if lastLoopPrintedRows > 0 {
					fmt.Print(fmt.Sprintf("\033[%dA", lastLoopPrintedRows)) // Move up X lines
				}
				lastLoopPrintedRows = len(rowsToPrint)

				fmt.Print("\033[J") // Clear everything below the current cursor position

				// Print any pending logs
				pb.logWriter.flush()

				for _, r := range rowsToPrint {
					fmt.Println(r)
				}
			}
		}
	}()
}

type ProgressBar struct {
	Name         string
	startedAt    time.Time
	finishedAt   time.Time
	isStageMode  bool
	maxProgress  int
	progress     int
	subLevel     int
	subBars      []*ProgressBar
	spinnerFrame int
	logWriter    progressBarLogsWriter
}

func (pb *ProgressBar) render(leftPadding int) []string {
	pb.cleanupOldBars()
	prefix := pb.renderPrefix()

	if pb.isStageMode {
		prefix += fmt.Sprintf(" Stage %d/%d", pb.progress, pb.maxProgress)
	}

	prefix += fmt.Sprintf(" - %s - ", fmt.Sprintf("%-*s", leftPadding, pb.Name))
	rows := []string{prefix + pb.renderBar(len(prefix))}

	for _, subBar := range pb.subBars {
		rows = append(rows, subBar.render(pb.getMaxNameSize())...)
	}
	pb.cleanupOldBars()

	return rows
}

func (pb *ProgressBar) cleanupOldBars() {
	const secondsToKeepBarAfterFinish = 10.0
	for i := len(pb.subBars) - 1; i >= 0; i-- {
		subBar := pb.subBars[i]

		// Check if should cleanup the bar
		if !subBar.finishedAt.IsZero() && time.Since(subBar.finishedAt).Seconds() > secondsToKeepBarAfterFinish {
			pb.subBars = append(pb.subBars[:i], pb.subBars[i+1:]...)
		}
	}
}

func (pb *ProgressBar) renderBar(leftPadding int) string {
	terminalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	barLength := terminalWidth - leftPadding - 6
	progressRatio := float64(pb.progress) / float64(pb.maxProgress)
	safeProgressRatio := min(progressRatio, 1)
	filledLength := int(float64(barLength) * safeProgressRatio)
	return "[" + fmt.Sprintf("%-4s", fmt.Sprintf("%d%%", int(progressRatio*100))) + strings.Repeat("#", filledLength) + strings.Repeat(" ", barLength-filledLength) + "]"
}

func (pb *ProgressBar) renderPrefix() string {
	var spinnerFrame string
	if pb.finishedAt.IsZero() {
		spinnerFrame = string(spinnerFrames[pb.spinnerFrame])
		pb.spinnerFrame = (pb.spinnerFrame + 1) % len(spinnerFrames)
	} else {
		spinnerFrame = string('✔')
	}

	var duration time.Duration
	if pb.finishedAt.IsZero() {
		duration = time.Since(pb.startedAt)
	} else {
		duration = pb.finishedAt.Sub(pb.startedAt)
	}
	h := int(duration.Hours())
	m := int(duration.Minutes()) % 60
	s := int(duration.Seconds()) % 60
	return fmt.Sprintf("%s%s%s", strings.Repeat("    ", pb.subLevel), spinnerFrame, fmt.Sprintf("[%02d:%02d:%02d]", h, m, s))
}

func (pb *ProgressBar) getMaxNameSize() int {
	maxLength := 0
	for _, sub := range pb.subBars {
		if len(sub.Name) > maxLength {
			maxLength = len(sub.Name)
		}
	}

	return maxLength
}

func (pb *ProgressBar) NewSubBar(name string, options ...Option) *ProgressBar {
	subPb := NewProgressBar(
		name,
		append(options,
			WithSubLevel(pb.subLevel+1),
			withWriter(pb.logWriter))...,
	)
	pb.subBars = append(pb.subBars, subPb)
	return subPb
}

func (pb *ProgressBar) Increment(amount int) {
	pb.progress += amount
}

func (pb *ProgressBar) Finish() {
	if pb.finishedAt.IsZero() {
		pb.finishedAt = time.Now()
	}
	pb.progress = pb.maxProgress
}

type Option func(*ProgressBar)

func WithMaxProgress(maxProgress int) Option {
	return func(s *ProgressBar) {
		s.maxProgress = maxProgress
	}
}

func WithStages(stages int) Option {
	return func(s *ProgressBar) {
		s.maxProgress = stages
		s.isStageMode = true
	}
}

func WithSubLevel(subLevel int) Option {
	return func(s *ProgressBar) {
		s.subLevel = subLevel
	}
}

func withWriter(writer progressBarLogsWriter) Option {
	return func(s *ProgressBar) {
		s.logWriter = writer
	}
}

func NewProgressBar(name string, options ...Option) *ProgressBar {
	pb := &ProgressBar{
		Name:        name,
		startedAt:   time.Now(),
		maxProgress: 100,
		logWriter: progressBarLogsWriter{
			buffer: make(chan string, 1000),
		},
	}

	for _, opt := range options {
		opt(pb)
	}

	return pb
}

func (pb *ProgressBar) UpdateProgressBar(options ...Option) {
	for _, opt := range options {
		opt(pb)
	}
}
