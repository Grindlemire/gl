package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// VertexBufferObject wraps the openGL VBO
type VertexBufferObject struct {
	addr uint32
}

// NewVBO creates a vertex buffer object and copies the vertices into it
func NewVBO(vertices []float32) (vbo VertexBufferObject) {
	gl.GenBuffers(1, &vbo.addr)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.addr)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	return vbo
}

// VertexArrayObject wraps the openGL VAO
type VertexArrayObject struct {
	addr uint32
}

// NewVAO creates a vertex array object
func NewVAO() (vao VertexArrayObject) {
	gl.GenVertexArrays(1, &vao.addr)
	gl.BindVertexArray(vao.addr)
	return vao
}

// MapAttribute maps data to a specific attribute from the VAO
// Take the data in the VAO (it points to the data loaded into the VBO) and map it to some
// input passed to the shaders. This takes a pointer to the program, the name of the input in GLSL,
// the offset into the data you set, the number of elements in the data you set, and the stride (how
// many floats between instances of this data)
func (vao VertexArrayObject) MapAttribute(program uint32, name string, offset int, size, stride int32) {
	attributeAddress := uint32(gl.GetAttribLocation(program, gl.Str(fmt.Sprintf("%s\x00", name))))
	gl.VertexAttribPointer(attributeAddress, size, gl.FLOAT, false, stride*4, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(attributeAddress)
}
