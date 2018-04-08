package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Transformation is the generic struct defining a transformation matrix used to convert coordinates
// The path for converting coordinates is Model -> World -> Camera -> Screen
// Model handles Model -> World
// View handles World -> Camera
// Projection handles Camera -> Screen
type Transformation struct {
	addr   int32      // the location in memory of the matrix
	matrix mgl32.Mat4 // the actual matrix value
}

// UpdateUniform sends an update to the openGL shader for the transformation matrix
// This is called when the transformation matrix has changed and we want to push that change
// to the shader
func (t *Transformation) UpdateUniform() {
	gl.UniformMatrix4fv(t.addr, 1, false, &t.matrix[0])
}

// UpdateMatrix updates the matrix to the new matrix
func (t *Transformation) UpdateMatrix(newMatrix mgl32.Mat4) {
	t.matrix = newMatrix
}

// GetAddr returns the address of the transformation matrix
func (t *Transformation) GetAddr() int32 {
	return t.addr
}

// GetMatrix returns the value of the transformation matrix
func (t *Transformation) GetMatrix() mgl32.Mat4 {
	return t.matrix
}

// Projection manages the projection matrix (it converts camera coordinates to screen coordinates)
// This maps the world onto a 2-d screen
type Projection struct {
	Transformation
}

// NewProjection creates a projection transformation matrix
// It takes the program pointer and the name of the trasnformation in GLSL
func NewProjection(program uint32, name string) (projection *Projection) {
	// create the transformation matrix
	matrix := mgl32.Perspective(mgl32.DegToRad(45.0), float32(winWidth)/winHeight, 0.1, 100.0)
	// get the location in memory where we need to place it
	addr := gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", name)))
	// load the data into the memory location
	gl.UniformMatrix4fv(addr, 1, false, &matrix[0])

	projection = &Projection{
		Transformation{
			addr:   addr,
			matrix: matrix,
		},
	}

	return projection
}

// View manages the view transformation matrix (it converts world coordinates to camera coordinates)
// This remaps everything in the world with respect to some camera somewhere
type View struct {
	Transformation
}

// NewView creates a view transformation matrix
// It takes the program pointer, name of the transformation in GLSL, and
// 3 3x1 matrices corresponding to where the eye is looking at, located at, and what direction is up
func NewView(program uint32, name string, position, target, up mgl32.Vec3) (view *View) {
	// create the view transformation matrix with
	matrix := mgl32.LookAtV(position, target, up)
	addr := gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", name)))
	gl.UniformMatrix4fv(addr, 1, false, &matrix[0])

	view = &View{
		Transformation{
			addr:   addr,
			matrix: matrix,
		},
	}

	return view
}

// UpdateCameraLocation updates the location and direction of the camera
func (v *View) UpdateCameraLocation(position, target, up mgl32.Vec3) {
	v.matrix = mgl32.LookAtV(position, target, up)
}

// Model handles the model transformation matrix (it converts model coordinates to world coordinates)
// This converts a standard model to placing it somewhere in the world
type Model struct {
	Transformation
}

// NewModel creates a model transformation matrix
// It takes the program pointer and name of the model transformation in GLSL
func NewModel(program uint32, name string) (model *Model) {
	// transform from world coordinates
	matrix := mgl32.Ident4()
	addr := gl.GetUniformLocation(program, gl.Str(fmt.Sprintf("%s\x00", name)))
	gl.UniformMatrix4fv(addr, 1, false, &matrix[0])

	model = &Model{
		Transformation{
			addr:   addr,
			matrix: matrix,
		},
	}

	return model
}
