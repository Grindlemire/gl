package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
)

// width and height of the window we are creating
const (
	winWidth  = 500
	winHeight = 500
)

// runs the program
func main() {
	runtime.LockOSThread()

	log.Printf("Starting textured cube!")

	window, err := initGlfw()
	if err != nil {
		log.Fatalf("Error initializing glfw: %v", err)
	}

	program, err := initOpenGL()
	if err != nil {
		log.Fatalf("Error initializing openGL: %v", err)
	}

	// create our transformations
	model := NewModel(program, "model")
	_ = NewView(program, "view", mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	_ = NewProjection(program, "projection")

	// load our data into our buffers
	vao := NewVAO()
	_ = NewVBO(cubeVertices) // we don't need the vbo after initializing it

	// load our texture
	texture, err := NewTexture(program, "texSampler", "wall.jpg")
	if err != nil {
		log.Fatalf("Error generating texture: %v\n", err)
	}

	// map our data into the shader
	vao.MapAttribute(program, "vert", 0, 3, 5)
	vao.MapAttribute(program, "vertTexCoord", 3, 2, 5)

	// enable depth of field and general constants
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	// draw a wireframe instead of filling
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture.textureID)

		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		// claculate new angle
		angle += elapsed
		model.UpdateMatrix(mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0}))

		// render
		gl.UseProgram(program)
		// This sends the updated model transformation to the shaders so we get rotation
		model.UpdateUniform()

		gl.BindVertexArray(vao.addr) // note this line is not needed now but will probably be needed when we have multiple vaos
		gl.DrawArrays(gl.TRIANGLES, 0, 6*6)

		window.SwapBuffers()
		glfw.PollEvents()
	}

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
func initOpenGL() (program uint32, err error) {
	err = gl.Init()
	if err != nil {
		return 0, errors.Wrap(err, "unable to initialize openGL")
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("OpenGL version: %s\n", version)

	vertexShader, err := compileShader(vertexShaderSrc, gl.VERTEX_SHADER)
	if err != nil {
		return 0, errors.Wrap(err, "unable to compile vertex shader")
	}

	fragmentShader, err := compileShader(fragShaderSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, errors.Wrap(err, "unable to compile fragment shader source")
	}

	program = gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		l := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(l))

		return 0, fmt.Errorf("failed to link program: %v", l)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	gl.UseProgram(program)
	return program, nil
}
