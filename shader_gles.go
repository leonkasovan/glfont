//go:build gles2
package glfont

import (
	"fmt"
	"strings"

	gl "github.com/leonkasovan/gl/v3.1/gles2"
)

// newProgram links the frag and vertex shader programs
func (r *FontRenderer_GLES) newProgram(GLSLVersion uint, vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShaderSource = fmt.Sprintf("#version %d es\nprecision mediump float;\n", GLSLVersion) + vertexShaderSource
	fragmentShaderSource = fmt.Sprintf("#version %d es\nprecision mediump float;\n", GLSLVersion) + fragmentShaderSource
	compileShader := func(source string, shaderType uint32) (uint32, error) {
		shader := gl.CreateShader(shaderType)

		csources, free := gl.Strs(source)
		gl.ShaderSource(shader, 1, csources, nil)
		free()
		gl.CompileShader(shader)

		var status int32
		gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

			return 0, fmt.Errorf("%v\nfailed to compile %v: %v", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)), source, log)
		}

		return shader, nil
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("%v\nfailed to link program: %v", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)), log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}
