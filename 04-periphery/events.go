package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

var winTitle = "GO-Events"
var winWidth, winHeight int32 = 800, 600

func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var event sdl.Event
	var running bool
	var err error

	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create surface %s\n", err)
		return 2
	}

	surface.FillRect(nil, 0)
	window.UpdateSurface()
	renderer, err = sdl.CreateRenderer(window, 2, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer.Destroy()

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			case *sdl.MouseButtonEvent:
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			case *sdl.MouseWheelEvent:
				fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y)
			}
		}
	}

	return 0
}

func main() {
	os.Exit(run())
}
