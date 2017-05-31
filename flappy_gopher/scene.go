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
	pipe *pipe
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

	p, err := newPipe(r)
	if err != nil {
		return nil, err
	}
	return &Scene{Bg: bg, bird: b, pipe: p}, nil
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
				s.Update()
				if s.bird.isDead() {
					drawTitle(r, "Game Over")
					time.Sleep(time.Second)
					s.Restart()
				}
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

//Update updates the bird's y axis
func (s *Scene) Update() {
	s.bird.update()
	s.pipe.update()
	s.bird.touch(s.pipe)
}

//Restart the scene
func (s *Scene) Restart() {
	s.bird.restart()
	s.pipe.restart()
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
	if err := s.pipe.paint(r); err != nil {
		return err
	}
	r.Present()
	return nil
}

//Destroy of Scene s
func (s *Scene) Destroy() {
	s.Bg.Destroy()
	s.bird.destroy()
	s.pipe.destroy()
}
