package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type quoteBlock struct {
	color sdl.Color
	x, y  int32
	h, w  int32
	font  *ttf.Font
	quote string
	auth  string
}

type quoteStruct []struct {
	ID         int    `json:"ID"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Link       string `json:"link"`
	CustomMeta struct {
		Source string `json:"Source"`
	} `json:"custom_meta"`
}

func newQuote() (string, error) {
	resp, err := http.Get("http://quotesondesign.com/wp-json/posts?filter[orderby]=rand&filter[posts_per_page]=1")
	if err != nil {
		return "", fmt.Errorf("Error getting quote: %v", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	return bodyString, err
}

var quoteable quoteStruct

func createQuoteStruct() error {
	quotejson, err := newQuote()
	if err != nil {
		return fmt.Errorf("error with quote json: %v", err)
	}
	error := json.Unmarshal([]byte(quotejson), &quoteable)
	return error
}

func todaysQuote() (string, string, error) {
	if err := createQuoteStruct(); err != nil {
		return "", "", err
	}
	quote := *&quoteable[0].Content
	auth := *&quoteable[0].Title
	return quote, auth, nil
}

func newQuoteBlock(r *sdl.Renderer) (*quoteBlock, error) {
	quotecolor := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	quotefont, err := ttf.OpenFont("fonts/LemonMilk.ttf", 250)
	if err != nil {
		return nil, fmt.Errorf("could not open quote font %v", err)
	}
	todaysQuote, auth, _ := todaysQuote()
	return &quoteBlock{quote: todaysQuote,
		font:  quotefont,
		color: quotecolor,
		auth:  auth,
		x:     100, y: 900, h: 100, w: 1000}, nil
}

func (quoteBlock *quoteBlock) paintQuote(r *sdl.Renderer) error {
	rect := &sdl.Rect{X: quoteBlock.x, Y: quoteBlock.y, W: quoteBlock.w, H: quoteBlock.h}
	s, err := quoteBlock.font.RenderUTF8Solid(quoteBlock.quote, quoteBlock.color)
	if err != nil {
		return fmt.Errorf("could not copy quote texture: %v", err)
	}
	texture, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create quote texture: %v", err)
	}
	if err := r.Copy(texture, nil, rect); err != nil {
		return fmt.Errorf("could not copy quote texture: %v", err)
	}
	r.Present()
	return nil
}

func (quoteBlock *quoteBlock) destroy() {
	quoteBlock.font.Close()
}
