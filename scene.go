package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	t               int
	clock           *clock
	weather         *weather
	weather_known   bool
	current_weather string
}

func newScene(r *sdl.Renderer) (*scene, error) {
	r.Clear()
	clock, err := newClock(r)
	if err != nil {
		return nil, fmt.Errorf("no running new lock: %v", err)
	}
	return &scene{t: 0, clock: clock}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		tick := time.Tick(1 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	default:
		return false
	}
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()
	clock, err := newClock(r)
	s.clock = clock
	if err != nil {
		return fmt.Errorf("no running new clock: %v", err)
	}
	if err := s.clock.paint(r); err != nil {
		return err
	}
	if time.Now().Format("00:00") == "00:01" {
		s.weather_known = false
	}
	if s.weather_known == false {
		weather, err := newWeather(r)
		if err != nil {
			return fmt.Errorf("cannot do new weather: %v", err)
		}
		s.weather = weather
		s.weather_known = true
		s.current_weather = s.weather.todaysWeather
	}
	if err := s.weather.weatherPaint(r); err != nil {
		return err
	}
	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.clock.destroy()
	s.weather.destroy()
}
