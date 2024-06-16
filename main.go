package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const paddleHeight int = 6
const paddleAsset rune = 0x2588
const ballAsset rune = 0x25CF

var cnt = '0'

var screenH, screenW int

type Paddle struct {
	x, y int
}

type Ball struct {
	x, y int
}

var paddleLeft, paddleRight *Paddle = &Paddle{}, &Paddle{}
var ball *Ball = &Ball{}

func renderHorizontal(screen tcell.Screen, x, y int, s string, style tcell.Style) {
	for _, ch := range s {
		screen.SetContent(x, y, ch, nil, style)
		x++
	}
}

func renderChar(screen tcell.Screen, x, y, w, h int, ch rune, style tcell.Style) {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			screen.SetContent(x+i, y+j, ch, nil, style)
		}
	}
}

func initScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize screen: %v\n", err)
	}

	err = screen.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize screen: %v\n", err)
	}
	return screen
}

func main() {
	// x-axis right, y-axis down
	screen := initScreen()
	screenW, screenH = screen.Size()
	initState(screen, paddleLeft, paddleRight, ball)

	// start the background process for getting the user input
	// if user presses any key, its detail is sent through a channel to the
	// foreground process
	keyPressedChan := initPlayerInputProcess(screen)
	for {
		renderState(screen, paddleLeft, paddleRight, ball)
		time.Sleep(16 * time.Millisecond)

		keyPressed := getKeyPressed(keyPressedChan)
		if keyPressed == "q" {
			screen.Fini()
			os.Exit(0)
			return
		} else if keyPressed == "up" {
			if paddleRight.y > 0 {
				paddleRight.y--
			}
		} else if keyPressed == "down" {
			if paddleRight.y < screenH-paddleHeight {
				paddleRight.y++
			}
		} else if keyPressed == "w" {
			if paddleLeft.y > 0 {
				paddleLeft.y--
			}
		} else if keyPressed == "s" {
			if paddleLeft.y < screenH-paddleHeight {
				paddleLeft.y++
			}
		}
	}
}

func getKeyPressed(keyPressedChan chan string) string {
	keyPressed := ""
	select {
	case keyPressed = <-keyPressedChan:
	default:
		keyPressed = ""
	}
	return keyPressed
}

func initState(screen tcell.Screen, paddleLeft, paddleRight *Paddle, ball *Ball) {
	paddleLeft.y = screenH/2 - paddleHeight/2
	paddleRight.y = screenH/2 - paddleHeight/2
	paddleLeft.x = 0
	paddleRight.x = screenW - 1
	ball.x = screenW / 2
	ball.y = screenH / 2
	renderState(screen, paddleLeft, paddleRight, ball)
}

func renderState(screen tcell.Screen, paddleLeft, paddleRight *Paddle, ball *Ball) {
	screen.Clear()
	screen.SetContent(0, 0, cnt, nil, tcell.StyleDefault)
	cnt++
	renderChar(screen, paddleLeft.x, paddleLeft.y, 1, paddleHeight, paddleAsset, tcell.StyleDefault)
	renderChar(screen, paddleRight.x, paddleRight.y, 1, paddleHeight, paddleAsset, tcell.StyleDefault)
	renderChar(screen, ball.x, ball.y, 1, 1, ballAsset, tcell.StyleDefault)
	screen.Show()
}

// starts the process for taking user inputs in the background
// and returns the channel
func initPlayerInputProcess(screen tcell.Screen) chan string {
	input := make(chan string)
	keyPressed := ""
	go func() {
		for {
			eventPoll := screen.PollEvent()
			switch event := eventPoll.(type) {
			case *tcell.EventKey:
				if event.Key() == tcell.KeyUp {
					keyPressed = "up"
				} else if event.Key() == tcell.KeyDown {
					keyPressed = "down"
				} else if event.Rune() == 'w' {
					keyPressed = "w"
				} else if event.Rune() == 's' {
					keyPressed = "s"
				} else if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC || event.Rune() == 'q' {
					keyPressed = "q"
				}
				input <- keyPressed
			}
		}
	}()
	return input
}
