package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

//Scene is a struct to represent a single frame
type Scene struct {
	Bg   *sdl.Texture
	bird *bird
}

//NewScene returns a scene structure by loading a background png and assigning the texture to the scene
func NewScene(r *sdl.Renderer) (*Scene, error) {
	//t = texture
	bg, err := img.LoadTexture(r, "res/imgs/background0.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}
	b, err := newBird(r)
	if err != nil {
		return nil, err
	}

	return &Scene{Bg: bg, bird: b}, nil
}

//Run implements the process of continuously calling paint()
func (s *Scene) Run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				if err := s.Paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

func (s *Scene) handleEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		s.bird.jump()
	case *sdl.MouseMotionEvent, *sdl.WindowEvent:
	default:
		log.Printf("unknown event %T\n", e)
	}
	return false
}

//Paint paints the scene into Renderer r
func (s *Scene) Paint(r *sdl.Renderer) error {
	r.Clear()

	if err := r.Copy(s.Bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err := s.bird.paint(r); err != nil {
		return err
	}
	r.Present()
	return nil
}

//Destroy of Scene s
func (s *Scene) Destroy() {
	s.Bg.Destroy()
}
