package main

import (
	"log"
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

// Camera manages the location of the camera
type Camera struct {
	position mgl32.Vec3
	front    mgl32.Vec3
	up       mgl32.Vec3

	yaw   float64
	pitch float64

	xoffset  float64
	yoffset  float64
	xpos     float64
	ypos     float64
	firstPos bool

	view *View

	keys map[string]*Key

	cameraSpeed float32
}

// NewCamera creates a new camera to manage the view matrix
func NewCamera(program uint32, name string, position, front, up mgl32.Vec3) (c *Camera) {
	c = &Camera{
		position: position,
		up:       up,
		front:    front,

		cameraSpeed: 2.5,
		yaw:         -90,
		firstPos:    true,

		view: NewView(program, name, position, position.Add(front), up),

		keys: map[string]*Key{
			"w": &Key{
				Key:     glfw.KeyW,
				Pressed: false,
			},
			"a": &Key{
				Key:     glfw.KeyA,
				Pressed: false,
			},
			"s": &Key{
				Key:     glfw.KeyS,
				Pressed: false,
			},
			"d": &Key{
				Key:     glfw.KeyD,
				Pressed: false,
			},
		},
	}

	return c
}

// ProcessMouseMove updates the camera based on mouse movement
func (c *Camera) ProcessMouseMove(xpos, ypos float64) {
	if c.firstPos {
		c.xpos = xpos
		c.ypos = ypos
		c.firstPos = false
		log.Printf("ENTERED AREA: X: %v | Y: %v", c.xpos, c.ypos)
		return
	}

	c.xoffset = -(c.xpos - xpos) * .2
	c.yoffset = (c.ypos - ypos) * .2
	c.ypos = ypos
	c.xpos = xpos

	c.yaw += c.xoffset
	c.pitch += c.yoffset

	if c.pitch > 89.0 {
		c.pitch = 89.0
	}

	if c.pitch < -89.0 {
		c.pitch = -89.0
	}

	log.Printf("YAW: %v PITCH: %v", c.yaw, c.pitch)

	x := float32(math.Cos(mgl64.DegToRad(c.yaw)) * math.Cos(mgl64.DegToRad(c.pitch)))
	y := float32(math.Sin(mgl64.DegToRad(c.pitch)))
	z := float32(math.Sin(mgl64.DegToRad(c.yaw)) * math.Cos(mgl64.DegToRad(c.pitch)))
	c.front = mgl32.Vec3{x, y, z}.Normalize()

	return
}

// ProcessKeyPress processes a key press for the camera
func (c *Camera) ProcessKeyPress(key glfw.Key, action glfw.Action, mods glfw.ModifierKey) {
	for _, potentialKey := range c.keys {
		if key == potentialKey.Key {
			if action == glfw.Press || action == glfw.Repeat {
				potentialKey.Pressed = true
				return
			}

			if action == glfw.Release {
				potentialKey.Pressed = false
				return
			}
		}
	}
}

// Update updates the camera based on the key press state
func (c *Camera) Update(deltaTime float32) {

	if c.keys["w"].Pressed {
		c.position = c.position.Add(c.front.Mul(c.cameraSpeed * deltaTime))
	}
	if c.keys["a"].Pressed {
		c.position = c.position.Sub(c.front.Cross(c.up).Mul(c.cameraSpeed * deltaTime))
	}
	if c.keys["s"].Pressed {
		c.position = c.position.Sub(c.front.Mul(c.cameraSpeed * deltaTime))
	}
	if c.keys["d"].Pressed {
		c.position = c.position.Add(c.front.Cross(c.up).Mul(c.cameraSpeed * deltaTime))
	}

	c.view.UpdateCameraLocation(c.position, c.position.Add(c.front), c.up)
	c.view.UpdateUniform()
}
