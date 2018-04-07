package main

import (
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

	log.Printf("Starting hello cube!")

	window, err := initGlfw()
	if err != nil {
		log.Fatalf("Error initializing glfw: %v", err)
	}

	program, err := initOpenGL()
	if err != nil {
		log.Fatalf("Error initializing openGL: %v", err)
	}

	// model, view, projection := getTransformations(program)
	// map to screen coordinates
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(winWidth)/winHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// map to camera coordinates
	view := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	// transform from world coordinates
	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// initialize our vertex attribute object and bind it
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// load our data into the buffer that the vao will look at
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(cubeVertices), gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	// set the attribute in teh vao for where the points are located
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(vertAttrib)

	// enable depth of field and general constants
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	// draw a wireframe instead of filling
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time
		// draw(vao, modelUniform, model, window, program, elapsed)

		// claculate new angle
		angle += elapsed
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// render
		gl.UseProgram(program)
		// This sends the updated model to the shaders so we get rotation
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(vao) // note this line is not needed now but will probably be needed when we have multiple vaos
		gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

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

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	gl.UseProgram(program)
	return program, nil
}
