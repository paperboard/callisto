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

type LightData struct {
  PointLightingLocationUniform int32
  PointLightingColorUniform    int32

  IsLightEmitterUniform        int32
  IsLightReceiverUniform       int32
}

var LIGHT LightData

func setLightUniforms(program uint32) {
  LIGHT.PointLightingLocationUniform = gl.GetUniformLocation(program, gl.Str("pointLightingLocationUniform\x00"))
  LIGHT.PointLightingColorUniform = gl.GetUniformLocation(program, gl.Str("pointLightingColorUniform\x00"))

  LIGHT.IsLightEmitterUniform = gl.GetUniformLocation(program, gl.Str("isLightEmitterUniform\x00"))
  LIGHT.IsLightReceiverUniform = gl.GetUniformLocation(program, gl.Str("isLightReceiverUniform\x00"))
}