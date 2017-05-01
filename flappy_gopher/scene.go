package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

//Scene is a struct to represent a single frame
type Scene struct {
	Bg *sdl.Texture
}

//NewScene returns a scene structure by loading a background png and assigning the texture to the scene
func NewScene(r *sdl.Renderer) (*Scene, error) {
	//t = texture
	t, err := img.LoadTexture(r, "res/imgs/background0.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	return &Scene{Bg: t}, nil
}

func (s *Scene) paint(r *sdl.Renderer) error {
	r.Clear()
	if err := r.Copy(s.Bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	r.Present()
	return nil
}

func (s *Scene) destroy() {
	s.Bg.Destroy()
}
