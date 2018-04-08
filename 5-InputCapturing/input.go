package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Key manages a key and whether it is currently being pressed or not
type Key struct {
	glfw.Key
	Pressed bool
}

// Keys is the global state of the captured keys we are listening for
var Keys = map[string]*Key{
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
}

// Mouse positions
var (
	XPos     float64
	DeltaX   float64
	YPos     float64
	DeltaY   float64
	firstPos = true

	Yaw   float64
	Pitch float64
)

// HandleKeyPress is the callback that is called anytime the program detects a key action
func HandleKeyPress(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
	camera.ProcessKeyPress(key, action, mods)
}

// HandleCursorMove handles the global state of the cursor
func HandleCursorMove(w *glfw.Window, xpos float64, ypos float64) {
	camera.ProcessMouseMove(xpos, ypos)
}
