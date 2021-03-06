/* Callisto - Yet another Solar System simulator
 *
 * Copyright (c) 2016, Valerian Saliou <valerian@valeriansaliou.name>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *   * Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *   * Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
  "github.com/go-gl/gl/v4.1-core/gl"
)

// ShaderData  Maps shader data
type ShaderData struct {
  VertexAttributes    uint32
  NormalAttributes    uint32
  VertexTextureCoords uint32
}

// InstanceShader  Stores shader data
var InstanceShader ShaderData

// ShaderVertex  Defines code for vertex shader
var ShaderVertex = `
#version 330

uniform mat4 projectionUniform;
uniform mat4 cameraUniform;
uniform mat4 modelUniform;
uniform mat3 normalUniform;

uniform vec3 pointLightingLocationUniform;
uniform vec3 pointLightingColorUniform;

uniform int isLightEmitterUniform;

in vec3 vertexAttributes;
in vec3 vertexNormalAttributes;
in vec2 vertexTextureCoords;

out vec3 N;
out vec2 shaderTextureCoords;
out vec3 vertexLighting;

void main() {
  vec4 modelViewPosition = modelUniform * vec4(vertexAttributes, 1);
  gl_Position = projectionUniform * cameraUniform * modelViewPosition;
  shaderTextureCoords = vertexTextureCoords;

  // Process lighting
  if (isLightEmitterUniform == 1) {
    // Light emitter
    vec3 lightDirection = normalize(pointLightingLocationUniform - modelViewPosition.xyz);

    vec3 transformedNormal = normalUniform * vertexNormalAttributes;
    float directionalLighting = max(dot(transformedNormal, lightDirection), 0.0);
    vertexLighting = pointLightingColorUniform * directionalLighting;
  } else {
    // Light receiver
    vertexLighting = vec3(1.0, 1.0, 1.0);
  }
}
` + "\x00"

// ShaderFragment  Defines code for fragment shader
var ShaderFragment = `
#version 330

uniform sampler2D textureUniform;

uniform int isLightReceiverUniform;

in vec2 shaderTextureCoords;
in vec3 vertexLighting;

out vec4 objectColor;

void main() {
  vec4 objectColorTexture = texture(textureUniform, shaderTextureCoords);

  // Apply lighting to pixel
  if (isLightReceiverUniform == 1) {
    objectColor = vec4(objectColorTexture.rgb * vertexLighting, objectColorTexture.a);
  } else {
    objectColor = objectColorTexture;
  }
}
` + "\x00"

func getShader() (*ShaderData) {
  return &InstanceShader
}

func initializeShaders(program uint32) {
  shader := getShader()

  // Bind buffer to shaders attributes
  shader.VertexAttributes = uint32(gl.GetAttribLocation(program, gl.Str("vertexAttributes\x00")))
  gl.EnableVertexAttribArray(shader.VertexAttributes)

  shader.NormalAttributes = uint32(gl.GetAttribLocation(program, gl.Str("vertexNormalAttributes\x00")))
  gl.EnableVertexAttribArray(shader.NormalAttributes)

  shader.VertexTextureCoords = uint32(gl.GetAttribLocation(program, gl.Str("vertexTextureCoords\x00")))
  gl.EnableVertexAttribArray(shader.VertexTextureCoords)

  // Bind misc. shaders uniforms
  setLightUniforms(program)
  setMatrixUniforms(program)
}
