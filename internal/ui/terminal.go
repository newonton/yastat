package ui

import (
	"fmt"
	"time"
)

func HideCursor()      { fmt.Print("\033[?25l") }
func ShowCursor()      { fmt.Print("\033[?25h") }
func ClearThreeLines() { fmt.Print("\033[3F\033[J") }

var firstRender = true

func Render(now, last time.Time, next time.Time) {
	if !firstRender {
		ClearThreeLines()
	} else {
		firstRender = false
	}

	fmt.Printf("run         %s\nlast record %s\nnext record in %s\n",
		now.Format("2006-01-02 15:04:05"),
		last.Format("2006-01-02 15:04:05"),
		time.Until(next).Truncate(time.Second))
}
