package main

var vertexShaderSrc = `
	#version 410

	uniform mat4 projection;
	uniform mat4 view;
	uniform mat4 model;

	in vec3 color;
	in vec3 vert;

	out vec3 Color;

	void main() {
		Color = color;
		gl_Position = projection * view * model * vec4(vert, 1.0);
	}
` + "\x00"

var fragShaderSrc = `
	#version 410
	in vec3 Color;

	out vec4 outputColor;

	void main(){
		outputColor = vec4(Color, 1.0);
	}
` + "\x00"
