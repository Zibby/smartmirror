package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

const (
	windowheight = 1080
	windowwidth  = 1920
)

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("Could not init sdl %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not init TTF: %v", err)
	}
	w, err := sdl.CreateWindow("", 0, 0, windowwidth, windowheight, sdl.WINDOW_FULLSCREEN_DESKTOP)
	if err != nil {
		return fmt.Errorf("could not draw window %v", err)
	}
	r, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return fmt.Errorf("could not create rendorer %v", err)
	}
	defer w.Destroy()
	r.SetDrawColor(0, 0, 0, 255)
	scene, err := newScene(r)
	if err != nil {
		return fmt.Errorf("could not newScene: %v", err)
	}
	if err := scene.paint(r); err != nil {
		return fmt.Errorf("could not paint scene: %v", err)
	}

	defer scene.destroy()

	events := make(chan sdl.Event)
	errc := scene.run(events, r)
	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}
