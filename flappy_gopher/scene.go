package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

//Scene is a struct to represent a single frame
type Scene struct {
	Time  int
	Bg    *sdl.Texture
	Birds []*sdl.Texture //a slice of textures
}

//NewScene returns a scene structure by loading a background png and assigning the texture to the scene
func NewScene(r *sdl.Renderer) (*Scene, error) {
	//t = texture
	bg, err := img.LoadTexture(r, "res/imgs/background0.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}
	var birds []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("res/imgs/birdframe%d.png", i)
		bird, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not load bird frame: %v", err)
		}
		birds = append(birds, bird)

	}

	return &Scene{Bg: bg, Birds: birds}, nil
}

//Run implements the process of continuously calling paint()
func (s *Scene) Run(ctx context.Context, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		for range time.Tick(10 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.Paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

//Paint paints the scene into Renderer r
func (s *Scene) Paint(r *sdl.Renderer) error {
	s.Time = (s.Time + 1) //increment time
	r.Clear()

	if err := r.Copy(s.Bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	i := (s.Time / 10) % len(s.Birds)
	//W,H are the dimensions of the bird_frame
	rect := &sdl.Rect{X: 10, Y: 300 - 43/2, W: 50, H: 43}
	if err := r.Copy(s.Birds[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird frame: %v", err)
	}

	r.Present()
	return nil
}

func (s *Scene) Destroy() {
	s.Bg.Destroy()
}
