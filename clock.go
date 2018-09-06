package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

type clock struct {
	timeNow    string
	color      sdl.Color
	font       *ttf.Font
	x, y, h, w int32
}

func newClock(r *sdl.Renderer) (*clock, error) {
	clocktime := time.Now().Format("15:04:05")
	clockcolor := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	clockfont, err := ttf.OpenFont("fonts/LemonMilk.ttf", 250)
	if err != nil {
		return nil, fmt.Errorf("could not open font %v", err)
	}
	return &clock{timeNow: clocktime,
		color: clockcolor,
		font:  clockfont,
		x:     20, y: 20, w: 250, h: 110}, nil
}

func (clock *clock) paint(r *sdl.Renderer) error {
	rect := &sdl.Rect{X: clock.x, Y: clock.y, W: clock.w, H: clock.h}
	s, err := clock.font.RenderUTF8Solid(string(clock.timeNow), clock.color)
	if err != nil {
		return fmt.Errorf("could not render clock: %v", err)
	}
	texture, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create surface: %v", err)
	}

	if err := r.Copy(texture, nil, rect); err != nil {
		return fmt.Errorf("could not copy clock texture: %v", err)
	}
	r.Present()
	return nil
}

func (clock *clock) destroy() {
	clock.font.Close()
}
