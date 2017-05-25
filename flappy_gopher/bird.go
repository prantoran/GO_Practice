package main

import (
	"fmt"
	"sync"

	img "github.com/veandco/go-sdl2/sdl_image"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	gravity   = 0.2
	jumpSpeed = -5
)

type bird struct {
	mu sync.RWMutex

	time     int
	textures []*sdl.Texture

	y, speed float64
	dead     bool
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

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.time++
	b.y += b.speed
	if b.y > 600 {
		b.dead = true
		b.speed = -0.8 * b.speed
		b.y = 600
	}
	b.speed += gravity
}

func (b *bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dead
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	i := (b.time / 10) % len(b.textures)
	//W,H are the dimensions of the bird_frame
	rect := &sdl.Rect{X: 10, Y: int32(b.y) - 43/2, W: 50, H: 43}
	if err := r.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird frame: %v", err)
	}
	return nil
}

func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, t := range b.textures {
		t.Destroy()
	}
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.speed = jumpSpeed
}
