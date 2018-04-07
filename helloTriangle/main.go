package main

import (
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/pkg/errors"
)

// width and height of the window we are creating
const (
	winWidth  = 500
	winHeight = 500
)

var triangle = []float32{
	0, 0.5, 0, // top
	-0.5, -0.5, 0, // left
	0.5, -0.5, 0, // right
}

// runs the program
func main() {
	runtime.LockOSThread()

	log.Printf("Starting hello triangle!")

	window, err := initGlfw()
	if err != nil {
		log.Fatalf("Error initializing glfw: %v", err)
	}

	program, err := initOpenGL()
	if err != nil {
		log.Fatalf("Error initializing openGL: %v", err)
	}

	vao := makeVAO(triangle)
	for !window.ShouldClose() {
		draw(vao, window, program)
	}

}

// draw draws each frame
func draw(vao uint32, window *glfw.Window, program uint32) {
	// clear previous frame
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	// use our program
	gl.UseProgram(program)

	// bind to the vao so when we call draw it knows which to draw
	gl.BindVertexArray(vao)
	// draw our array. Tell it to draw triangles, start at the first vertex and draw three vertices
	countVertices := int32(len(triangle) / 3)
	gl.DrawArrays(gl.TRIANGLES, 0, countVertices)

	// check for keyboard events
	glfw.PollEvents()
	// swap to our new drawn buffers
	window.SwapBuffers()
}

// makeVAO takes a list of vertices and makes a vertex array object from them
func makeVAO(vertices []float32) (vao uint32) {
	// create vertex buffer object
	var vbo uint32
	// generate the buffer with this memory (we only want 1 of them)
	gl.GenBuffers(1, &vbo)
	// bind the buffer to our variable and say it is an array buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// place the data  in the buffer and tell it how long the array is
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	// generate the vertex array (we only want 1 array)
	gl.GenVertexArrays(1, &vao)
	// bind the vertex array to our variable
	gl.BindVertexArray(vao)
	// bind our vao to the buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// tell the vao about our data. args(index to start at, the number of vertices per point, the type, do we normalize, how much space each point makes, ?)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	// enable the vao
	gl.EnableVertexAttribArray(0)

	return vao
}

// compileShader will take the GLSL raw source and compile it to a shader
func compileShader(source string, shaderType uint32) (shader uint32, err error) {
	// initialize a shader for whatever type we are creating
	shader = gl.CreateShader(shaderType)

	// convert the string source to a c string
	csources, free := gl.Strs(source)

	// point the c lib at the string memeory (we are only using 1 string)
	gl.ShaderSource(shader, 1, csources, nil)
	// free up the c string after the shader has used it
	free()
	// try to compile the GLSL into machine code
	gl.CompileShader(shader)

	// error handling
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var loglength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &loglength)

		// fill a string with a bunch of C nulls so we can null terminate the string
		l := strings.Repeat("\x00", int(loglength+1))
		gl.GetShaderInfoLog(shader, loglength, nil, gl.Str(l))

		return 0, errors.New(l)
	}
	return shader, nil
}

// initGlfw initiallizes our window
func initGlfw() (window *glfw.Window, err error) {
	err = glfw.Init()
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize glfw")
	}

	// lets us resize the window
	glfw.WindowHint(glfw.Resizable, glfw.True)

	// sets the version of openGL we will be using
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	// set the profile for compatibility
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// actually create the window with the title (window and monitor are nil here)
	window, err = glfw.CreateWindow(winWidth, winHeight, "Hello Triangle", nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create the window")
	}
	// bind it to this thread
	window.MakeContextCurrent()

	return window, nil
}

// initOpenGL initializes openGL
func initOpenGL() (prog uint32, err error) {
	err = gl.Init()
	if err != nil {
		return 0, errors.Wrap(err, "unable to initialize openGL")
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("OpenGL version: %s\n", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, errors.Wrap(err, "unable to compile vertex shader")
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, errors.Wrap(err, "unable to compile fragment shader source")
	}

	prog = gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog, nil
}

// vertex shader program that just passes through the posistion. Null terminate it C
var vertexShaderSource = `
	#version 410
	in vec3 vp;
	void main() {
		gl_Position = vec4(vp, 1.0);
	}
` + "\x00"

// fragment shader program that just sets the color of the fragment to white. Null terminate it for C
var fragmentShaderSource = `
	#version 410
	out vec4 frag_color;
	void main(){
		frag_color = vec4(1, 1, 1, 1);
	}
` + "\x00"
