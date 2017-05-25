package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

const (
	//TITLEDELAY is multiplier to time.sleep of title
	TITLEDELAY = 2
)

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

func run() error {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("Could not initialize ttf: %v", err)
	}
	defer ttf.Quit()

	//w = window, r = renderer
	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy() //shall destroy window and renderer

	if err = drawTitle(r, "intro"); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(TITLEDELAY * time.Second)

	s, err := NewScene(r)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.Destroy()

	events := make(chan sdl.Event)
	errc := s.Run(events, r)

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err

		}
	}
}

func drawTitle(r *sdl.Renderer, text string) error {
	r.Clear()
	if text == "intro" {
		f, err := ttf.OpenFont("res/fonts/flappy.ttf", 20)
		if err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		defer f.Close()

		c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
		//s = surface
		s, err := f.RenderUTF8_Solid("flappy gopher", c)
		if err != nil {
			return fmt.Errorf("could not render title: %v", err)
		}
		defer s.Free()

		t, err := r.CreateTextureFromSurface(s)
		if err != nil {
			return fmt.Errorf("could not create texture: %v", err)
		}
		defer t.Destroy()

		if err := r.Copy(t, nil, nil); err != nil {
			return fmt.Errorf("could not copy texture: %v", err)
		}
		r.Present()
	}
	if text == "Game Over" {
		f, err := ttf.OpenFont("res/fonts/flappy.ttf", 20)
		if err != nil {
			return fmt.Errorf("could not load font: %v", err)
		}
		defer f.Close()

		c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
		//s = surface
		s, err := f.RenderUTF8_Solid("flappy gopher", c)
		if err != nil {
			return fmt.Errorf("could not render title: %v", err)
		}
		defer s.Free()

		t, err := r.CreateTextureFromSurface(s)
		if err != nil {
			return fmt.Errorf("could not create texture: %v", err)
		}
		defer t.Destroy()

		if err := r.Copy(t, nil, nil); err != nil {
			return fmt.Errorf("could not copy texture: %v", err)
		}
		r.Present()
	}
	return nil
}
