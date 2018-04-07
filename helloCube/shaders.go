package main

var vertexShaderSrc = `
	#version 410

	uniform mat4 projection;
	uniform mat4 view;
	uniform mat4 model;

	in vec3 vert;

	void main() {
		gl_Position = projection * view * model * vec4(vert, 1.0);
	}
` + "\x00"

var fragShaderSrc = `
	#version 410
	out vec4 outputColor;
	void main(){
		outputColor = vec4(0.0, 1.0, 0.0, 1.0);
	}
` + "\x00"
