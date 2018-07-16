package main

import (

  _ "image/png"
  "cryp-8/cpu"
  "fmt"
  "time"
  "log"
  "strings"

  "math/rand"
  "os"
  "runtime"
  "github.com/go-gl/gl/v2.1/gl"
  "github.com/go-gl/glfw/v3.2/glfw"
)

type app struct {
  cpu *cpu.CPU
}

const (
  width  = 512
  height = 256

  vertexShaderSource = `
    #version 410
    in vec3 vp;
    void main() {
      gl_Position = vec4(vp, 1.0);
    }
  ` + "\x00"

  fragmentShaderSource = `
    #version 410
    out vec4 frag_colour;
    void main() {
      frag_colour = vec4(1, 1, 1, 1.0);
    }
  ` + "\x00"

  rows    = 64
  columns = 32
  threshold = 0.15
  fps = 60
)

var (
  square = []float32{
    -0.5, 0.5, 0,
    -0.5, -0.5, 0,
    0.5, -0.5, 0,

    -0.5, 0.5, 0,
    0.5, 0.5, 0,
    0.5, -0.5, 0,
  }
)

type cell struct {
    drawable uint32


    alive     bool
    aliveNext bool

    x int
    y int
}

// checkState determines the state of the cell for the next tick of the game.
func (c *cell) checkState(cells [][]*cell) {
    c.alive = c.aliveNext
    c.aliveNext = c.alive
    
    liveCount := c.liveNeighbors(cells)
    if c.alive {
        // 1. Any live cell with fewer than two live neighbours dies, as if caused by underpopulation.
        if liveCount < 2 {
            c.aliveNext = false
        }
        
        // 2. Any live cell with two or three live neighbours lives on to the next generation.
        if liveCount == 2 || liveCount == 3 {
            c.aliveNext = true
        }
        
        // 3. Any live cell with more than three live neighbours dies, as if by overpopulation.
        if liveCount > 3 {
            c.aliveNext = false
        }
    } else {
        // 4. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
        if liveCount == 3 {
            c.aliveNext = true
        }
    }
}

func (h *app) onKey(w *glfw.Window, key glfw.Key, scancode int,
  action glfw.Action, mods glfw.ModifierKey) {
  if key == glfw.KeyEscape && action == glfw.Press {
    w.SetShouldClose(true)
  }
  switch key {
    case glfw.Key0:
      h.cpu.SetKey(0)
    case glfw.Key1:
      h.cpu.SetKey(1)
    case glfw.Key2:
      h.cpu.SetKey(2)
    case glfw.Key3:
      h.cpu.SetKey(3)
    case glfw.Key4:
      h.cpu.SetKey(4)
    case glfw.Key5:
      h.cpu.SetKey(5)
    case glfw.Key6:
      h.cpu.SetKey(6)
    case glfw.Key7:
      h.cpu.SetKey(7)
    case glfw.Key8:
      h.cpu.SetKey(8)
    case glfw.Key9:
      h.cpu.SetKey(9)
    case glfw.KeyA:
      h.cpu.SetKey(0xa)
    case glfw.KeyB:
      h.cpu.SetKey(0xb)
    case glfw.KeyC:
      h.cpu.SetKey(0xc)
    case glfw.KeyD:
      h.cpu.SetKey(0xd)
    case glfw.KeyE:
      h.cpu.SetKey(0xe)
    case glfw.KeyF:
      h.cpu.SetKey(0xf)
  }
}


// liveNeighbors returns the number of live neighbors for a cell.
func (c *cell) liveNeighbors(cells [][]*cell) int {
    var liveCount int
    add := func(x, y int) {
        // If we're at an edge, check the other side of the board.
        if x == len(cells) {
            x = 0
        } else if x == -1 {
            x = len(cells) - 1
        }
        if y == len(cells[x]) {
            y = 0
        } else if y == -1 {
            y = len(cells[x]) - 1
        }
        
        if cells[x][y].alive {
            liveCount++
        }
    }
    
    add(c.x-1, c.y)   // To the left
    add(c.x+1, c.y)   // To the right
    add(c.x, c.y+1)   // up
    add(c.x, c.y-1)   // down
    add(c.x-1, c.y+1) // top-left
    add(c.x+1, c.y+1) // top-right
    add(c.x-1, c.y-1) // bottom-left
    add(c.x+1, c.y-1) // bottom-right
    
    return liveCount
}

func main() {

  runtime.LockOSThread()

  window := initGlfw()
  defer glfw.Terminate()
  program := initOpenGL()

  fmt.Println("Welcome to cryp-8, the only chip-8 emulator in existence.")
  cpu := cpu.NewCPU()
  h := app{&cpu}
  window.SetKeyCallback(h.onKey)

  f, err := os.Open("./15.ch8")

  b1 := make([]byte, 4096-0x200)
  _, err = f.Read(b1)
  if err != nil {
    panic(err)
  }
  cpu.LoadRom(b1)
  var iteration_times [100]float64
    cells := makeCells()

  for i := 0; !window.ShouldClose(); i %= 100 {
    t := time.Now() 
    cpu.RunCycle()
    glfw.PollEvents()

    if cpu.RefreshScreen {        
      for x := range cells {
        for y, c := range cells[x] {
          c.alive = cpu.Display()[64*(31 - y) + (x)]
        }
      }
      draw(cells, window, program)
      draw_(&cpu)
      cpu.RefreshScreen = false
    }

    if i == 99 {
      fmt.Printf("cycles per second: %v\n", calcCPS(iteration_times[:]))
    }
    time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
    iteration_times[i] = time.Since(t).Seconds()
    i++
  }
}


func draw_(cpu *cpu.CPU) {
  f, err := os.Create("./out.img")
  defer f.Close()
  if err != nil {
    panic(err)
  }
  out := []byte{}
  // fmt.Println("len", len(cpu.Display()))
  for _, val := range cpu.Display() {
    if val {
      out = append(out, byte(255))
    } else {
      out = append(out, byte(0))
    }
  }
  f.Write(out)
}

func calcCPS(x []float64) float64 {
  var total float64 = 0
  for _, value:= range x {
    total += value
  }
  return float64(len(x))/total
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
  if err := glfw.Init(); err != nil {
    panic(err)
  }
  glfw.WindowHint(glfw.Resizable, glfw.False)
  glfw.WindowHint(glfw.ContextVersionMajor, 4)
  glfw.WindowHint(glfw.ContextVersionMinor, 1)
  glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
  glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

  window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
  if err != nil {
    panic(err)
  }
  window.MakeContextCurrent()

  return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
  if err := gl.Init(); err != nil {
    panic(err)
  }
  version := gl.GoStr(gl.GetString(gl.VERSION))
  log.Println("OpenGL version", version)

  vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
  if err != nil {
    panic(err)
  }

  fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
  if err != nil {
    panic(err)
  }

  prog := gl.CreateProgram()
  gl.AttachShader(prog, vertexShader)
  gl.AttachShader(prog, fragmentShader)
  gl.LinkProgram(prog)
  return prog
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
  var vbo uint32
  gl.GenBuffers(1, &vbo)
  gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
  gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

  var vao uint32
  gl.GenVertexArrays(1, &vao)
  gl.BindVertexArray(vao)
  gl.EnableVertexAttribArray(0)
  gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
  gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

  return vao
}

func compileShader(source string, shaderType uint32) (uint32, error) {
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

    return 0, fmt.Errorf("failed to compile %v: %v", source, log)
  }

  return shader, nil
}

func draw(cells [][]*cell, window *glfw.Window, program uint32) {
  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
  gl.UseProgram(program)

  for x := range cells {
    for _, c := range cells[x] {
      c.draw()
    }
  }

  glfw.PollEvents()
  window.SwapBuffers()
}


func makeCells() [][]*cell {
    rand.Seed(time.Now().UnixNano())

    cells := make([][]*cell, rows)
    for x := 0; x < rows; x++ {
        for y := 0; y < columns; y++ {
            c := newCell(x, y)
            
            c.alive = true
            c.aliveNext = c.alive
            
            cells[x] = append(cells[x], c)
        }
    }

  return cells
}

func newCell(x, y int) *cell {
  // fmt.Println(x,y)
  points := make([]float32, len(square), len(square))
  copy(points, square)

  for i := 0; i < len(points); i++ {
    var position float32
    var size float32
    switch i % 3 {
    case 0:
      size = 1.0 / float32(rows)
      position = float32(x) * size
      // fmt.Println("x=",position)
    case 1:
      size = 1.0 / float32(columns)
      position = float32(y) * size
      // fmt.Println("y=",position)
    default:
      continue
    }

    if points[i] < 0 {
      points[i] = (position * 2) - 1
    } else {
      points[i] = ((position + size) * 2) - 1
    }
  }

  return &cell{
    drawable: makeVao(points),

    x: x,
    y: y,
  }
}
func (c *cell) draw() {
    if !c.alive {
            return
    }

    gl.BindVertexArray(c.drawable)
    gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}