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
  "log"
  "os"
  "runtime"

  "github.com/go-gl/gl/v4.1-core/gl"
  "github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
  // GLFW event handling must run on the main OS thread
  runtime.LockOSThread()

  dir, err := importPathToDir("github.com/valeriansaliou/callisto")
  if err != nil {
    log.Fatalln("Unable to find Go package in your GOPATH; needed to load assets:", err)
  }

  err = os.Chdir(dir)
  if err != nil {
    log.Panicln("os.Chdir:", err)
  }
}

func main() {
  var (
    err    error
    window *glfw.Window
    vao    uint32
  )

  // Create window
  if err = glfw.Init(); err != nil {
      log.Fatalln("Failed to initialize glfw:", err)
  }
  defer glfw.Terminate()

  glfw.WindowHint(glfw.Resizable, glfw.False)
  glfw.WindowHint(glfw.ContextVersionMajor, 4)
  glfw.WindowHint(glfw.ContextVersionMinor, 1)
  glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
  glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

  // Create window
  monitor := glfw.GetPrimaryMonitor()

  initializeWindow(monitor)

  window, err = glfw.CreateWindow(getWindowData().Width, getWindowData().Height, ConfigWindowTitle, nil, nil)
  if err != nil {
    panic(err)
  }

  // Full-screen window?
  if ConfigWindowFullScreen == true {
    // Adjust window to full-screen mode once we got the screen DPI (Retina screens)
    adjustWindow(window)

    // Re-create window to match full screen size w/ good framebuffer size (keeps high-DPI resolution)
    window.SetFramebufferSizeCallback(handleAdjustWindow)

    window.Destroy()

    window, err = glfw.CreateWindow(getWindowData().Width, getWindowData().Height, ConfigWindowTitle, monitor, nil)
    if err != nil {
      panic(err)
    }
  }

  // Bind window context
  window.MakeContextCurrent()

  // Bind key listeners
  window.SetInputMode(glfw.CursorMode, glfw.CursorHidden);

  window.SetKeyCallback(handleKey)
  window.SetCursorPosCallback(handleMouseCursor)
  window.SetScrollCallback(handleMouseScroll)

  // Initialize OpenGL
  gl.Init()

  // Configure the shaders program
  program, err := newProgram(ShaderVertex, ShaderFragment)
  if err != nil {
    panic(err)
  }

  gl.UseProgram(program)

  // Create environment
  createProjection(program)
  createCamera(program)

  // Create the VAO (Vertex Array Objects)
  // Notice: this stores links between attributes and active vertex data
  gl.GenVertexArrays(1, &vao)

  // Load the map of stellar objects + voidbox (aKa skybox)
  voidbox := loadObjects("voidbox")
  stellar := loadObjects("stellar")

  // Apply orbit traces
  createOrbitTraces(stellar, program, vao)

  // Create all object buffers
  createAllBuffers(voidbox, program, vao)
  createAllBuffers(stellar, program, vao)

  // Initialize shaders
  initializeShaders(program)

  // Configure global settings
  gl.Enable(gl.DEPTH_TEST)
  gl.Enable(gl.TEXTURE_2D)
  gl.Enable(gl.BLEND)
  gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
  gl.DepthFunc(gl.LESS)
  gl.ClearColor(0.0, 0.0, 0.0, 1.0)

  // Render loop
  for !window.ShouldClose() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    // Global routines
    updateElaspedTime(glfw.GetTime())
    gl.UseProgram(program)

    // Initialize stack matrix
    initializeMatrix()

    // Update context
    updateCamera()

    // Bind context
    bindProjection()
    bindCamera()

    // Render skybox
    gl.Disable(gl.DEPTH_TEST)
    renderObjects(voidbox, program)
    gl.Enable(gl.DEPTH_TEST)

    // Render all stellar objects in the map
    renderObjects(stellar, program)

    glfw.PollEvents()
    window.SwapBuffers()

    // Defer next update
    deferSceneUpdate()
  }
}
