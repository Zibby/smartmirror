package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	t              int
	clock          *clock
	weather        *weather
	weatherKnown   bool
	currentWeather string
	quote          *quoteBlock
	quoteKnown     bool
	currentQuote   string
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
		tick := time.Tick(1 * time.Second)
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
		s.weatherKnown = false
		s.quoteKnown = false
	}
	if s.weatherKnown == false {
		weather, err := newWeather(r)
		if err != nil {
			return fmt.Errorf("cannot do new weather: %v", err)
		}
		s.weather = weather
		s.weatherKnown = true
		s.currentWeather = s.weather.todaysWeather
	}
	if err := s.weather.weatherPaint(r); err != nil {
		return err
	}
	if s.quoteKnown == false {
		quote, err := newQuoteBlock(r)
		if err != nil {
			return fmt.Errorf("could not do new quote: %v", err)
		}
		s.quote = quote
		s.quoteKnown = true
		s.currentQuote = quote.quote
	}

	if err := s.quote.paintQuote(r); err != nil {
		return err
	}
	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.clock.destroy()
	s.weather.destroy()
	s.quote.destroy()
}
