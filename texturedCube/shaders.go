package main

var vertexShaderSrc = `
	#version 410

	uniform mat4 projection;
	uniform mat4 view;
	uniform mat4 model;

	in vec2 vertTexCoord;
	in vec3 vert;

	out vec2 fragTexCoord;

	void main() {
		fragTexCoord = vertTexCoord;
		gl_Position = projection * view * model * vec4(vert, 1.0);
	}
` + "\x00"

var fragShaderSrc = `
	#version 410
	uniform sampler2D texSampler;

	in vec2 fragTexCoord;

	out vec4 outputColor;

	void main(){
		outputColor = texture(texSampler, fragTexCoord);

	}
` + "\x00"
