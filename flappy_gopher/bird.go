package main

import (
	"fmt"

	img "github.com/veandco/go-sdl2/sdl_image"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	gravity   = 0.2
	jumpSpeed = -5
)

type bird struct {
	textures []*sdl.Texture
	time     int
	y, speed float64
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("res/imgs/birdframe%d.png", i)
		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not load bird frame: %v", err)
		}
		textures = append(textures, texture)

	}
	return &bird{textures: textures, y: 300, speed: 0}, nil
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.time++
	b.y += b.speed

	if b.y > 600 {
		b.speed = -0.8 * b.speed
		b.y = 600
	}

	b.speed += gravity

	i := (b.time / 10) % len(b.textures)
	//W,H are the dimensions of the bird_frame
	rect := &sdl.Rect{X: 10, Y: int32(b.y) - 43/2, W: 50, H: 43}
	if err := r.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird frame: %v", err)
	}
	return nil
}

func (b *bird) destroy() {
	for _, t := range b.textures {
		t.Destroy()
	}
}

func (b *bird) jump() {
	b.speed = jumpSpeed
}
