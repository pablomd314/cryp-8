package cpu

import (
  "math/rand"
  "time"
  "fmt"
  "errors"
)

type CPU struct {
  memory  [4096]uint8
  i     uint16
  pc    uint16
  v   [16]uint8
  stack [16]uint16
  sp    uint8
  dtimer  uint8
  stimer  uint8
  key   [16]bool
  display [64*32]bool
  RefreshScreen bool
}

func max(a, b uint8) uint8 {
    if a > b {
        return a
    }
    return b
}

func (cpu *CPU) SetKey(k uint8) {
  cpu.key[k] = true
  fmt.Printf("key set %x\n",k)
}

func (cpu *CPU) Display() []bool {
  return cpu.display[:]
}

func getAddress(instruction uint16) uint16 {
  return instruction & 0x0FFF 
}

func get8BitConstant(instruction uint16) uint8 {
  return uint8(instruction & 0x00FF)
}

func get4BitConstant(instruction uint16) uint8 {
  return uint8(instruction & 0x000F)
}

func getX(instruction uint16) uint8 {
  return uint8((instruction & 0x0F00) >> 8)
}

func getY(instruction uint16) uint8 {
  return uint8((instruction & 0x00F0) >> 4)
}

func fontAddress(font uint8) uint16 {
  if font > 0xF {
    panic(errors.New("Don't have fonts beyond 0xF"))
  }
  return uint16(5*font)
}

func NewCPU() CPU {
  rand.Seed(time.Now().UTC().UnixNano())
  var cpu CPU
  cpu.pc = 0x200
  copy(cpu.memory[:], fonts[:])
  cpu.RefreshScreen = false
  return cpu
}

func (cpu *CPU) LoadRom(buff []uint8) {
  for i, val := range buff {
    cpu.memory[0x200 + i] = val
  } 
}

func (cpu *CPU) RunCycle() {
  instruction := uint16(cpu.memory[cpu.pc]) << 8 | uint16(cpu.memory[cpu.pc + 1]);
  cpu.executeInstruction(instruction)
  
  cpu.dtimer = max(0, cpu.dtimer - 1)
  if cpu.stimer == 1 {
    // fmt.Println("BOOP")
  }
  
  cpu.stimer = max(0, cpu.stimer - 1)
}

func (cpu *CPU) clearKeys() {
  for i := range cpu.key {
    cpu.key[i] = false
  }
}

func (cpu *CPU) executeInstruction(instruction uint16) {
  // fmt.Printf("instruction %x\n", instruction)
  switch 0xF000 & instruction {
    case 0x0000:
      switch 0x00FF & instruction {
        case 0x00E0:
          for i, _ := range cpu.display {
            cpu.display[i] = false
          }
          cpu.RefreshScreen = true
          cpu.pc += 2
        case 0x00EE:
          cpu.sp--
          cpu.pc = cpu.stack[cpu.sp] + 2
      }
    case 0x1000:
      cpu.pc = getAddress(instruction)
    case 0x2000:
      cpu.stack[cpu.sp] = cpu.pc
      cpu.sp++
      cpu.pc = getAddress(instruction)
    case 0x3000:
      if cpu.getRegister(getX(instruction)) == get8BitConstant(instruction) {
        cpu.pc  += 2
      }
      cpu.pc    += 2
    case 0x4000:
      if cpu.getRegister(getX(instruction)) != get8BitConstant(instruction) {
        cpu.pc  += 2
      }
      cpu.pc    += 2
    case 0x5000:
      if cpu.getRegister(getX(instruction)) == cpu.getRegister(getY(instruction)) {
        cpu.pc  += 2
      }
      cpu.pc    += 2
    case 0x6000:
      cpu.setRegister(getX(instruction), get8BitConstant(instruction))
      cpu.pc    += 2
    case 0x7000:
      cpu.setRegister(getX(instruction), cpu.getRegister(getX(instruction)) + get8BitConstant(instruction))
      cpu.pc    += 2
    case 0x8000:
      vx := cpu.getRegister(getX(instruction))
      vy := cpu.getRegister(getY(instruction))
      switch 0x000F & instruction {
        case 0x0000:
          cpu.setRegister(getX(instruction), vy)
        case 0x0001:
          cpu.setRegister(getX(instruction), vx | vy)
        case 0x0002:
          cpu.setRegister(getX(instruction), vx & vy)
        case 0x0003:
          cpu.setRegister(getX(instruction), vx ^ vy)
        case 0x0004:
          cpu.setRegister(getX(instruction), vx + vy)
          cpu.setRegister(0xF, 0)
          if vx > 0xFF - vy {
            cpu.setRegister(0xF, 1)
          }
        case 0x0005:
          cpu.setRegister(getX(instruction), vx - vy)
          cpu.setRegister(0xF, 0)
          if vy > vx {
            cpu.setRegister(0xF, 1)
          }
        case 0x0006:
          cpu.setRegister(getX(instruction), vy >> 1)
          cpu.setRegister(0xF, vy & 0x1)
        case 0x0007:
          cpu.setRegister(getX(instruction), vy - vx)
          cpu.setRegister(0xF, 0)
          if vx > vy {
            cpu.setRegister(0xF, 1)
          }
        case 0x000E:
          cpu.setRegister(getX(instruction), vy << 1)
          cpu.setRegister(getY(instruction), vy << 1)
          cpu.setRegister(0xF, vy & 0x80)
      }
      cpu.pc += 2
    case 0x9000:
      if cpu.getRegister(getX(instruction)) != cpu.getRegister(getY(instruction)) {
        cpu.pc  += 2
      }
      cpu.pc    += 2
    case 0xA000:
      cpu.i     = instruction & getAddress(instruction)
      cpu.pc    += 2
    case 0xB000:
      cpu.pc    = uint16(cpu.getRegister(0)) + getAddress(instruction)
    case 0xC000:
      cpu.setRegister(getX(instruction), uint8(rand.Uint32()) & get8BitConstant(instruction))
      cpu.pc    += 2
    case 0xD000:
      vx := uint16(cpu.getRegister(getX(instruction)))
      vy := uint16(cpu.getRegister(getY(instruction)))
      n  := uint16(get4BitConstant(instruction))
      var pixel uint8
      cpu.setRegister(0xF, 0)
      for j := uint16(0); j < n; j++ {
        pixel = cpu.memory[cpu.i + uint16(j)]
        for k := uint16(0); k < 8; k++ {
          if (pixel & (0x80 >> k)) == (0x80 >> k) { //pixel is set
            // fmt.Printf("drawing pixel %v, row  %v\n", k, j)
            if cpu.display[vx + k + (vy + j)*64] {
              cpu.setRegister(0xF, 1)
            }
            // fmt.Printf("display index %v\n", vx + k + (vy + j)*64)
            cpu.display[vx + k + (vy + j)*64] = !cpu.display[vx + k + (vy + j)*64]
          } 
        }
      }
      cpu.RefreshScreen = true
      cpu.pc += 2
    case 0xE000:
      fmt.Printf("got key %x\n", cpu.getRegister(getX(instruction)))     
      switch 0x00FF & instruction {
        case 0x009E:
          if cpu.keyPressed(cpu.getRegister(getX(instruction))) {
            cpu.pc += 2
          }
          cpu.key[cpu.getRegister(getX(instruction))] = false
          cpu.pc += 2
        case 0x00A1:
          if !cpu.keyPressed(cpu.getRegister(getX(instruction))) {
            cpu.pc += 2
          }
          cpu.key[cpu.getRegister(getX(instruction))] = false
          cpu.pc += 2
      }
    case 0xF000:
      switch 0x00FF & instruction {
        case 0x0007:
          cpu.setRegister(getX(instruction), cpu.dtimer)
          cpu.pc += 2
        case 0x000A:
          if k := cpu.getKey(); k != 0xFF {
            fmt.Printf("got key %x\n", k)
            cpu.clearKeys()
            cpu.setRegister(getX(instruction), k)
            cpu.pc  += 2
          }
          //no key pressed 
        case 0x0015:
          cpu.dtimer = cpu.getRegister(getX(instruction))
          cpu.pc    += 2
        case 0x0018:
          cpu.stimer = cpu.getRegister(getX(instruction))
          cpu.pc    += 2
        case 0x001E:
          cpu.i     += uint16(cpu.getRegister(getX(instruction)))
          cpu.pc    += 2
          // check notes on wiki, VF might be set
        case 0x0029:
          cpu.i      = fontAddress(cpu.getRegister(getX(instruction)))
          cpu.pc    += 2
        case 0x0033:
          vx := cpu.getRegister(getX(instruction))
          cpu.memory[cpu.i]   = vx / 100 
          cpu.memory[cpu.i+1] = (vx / 10) % 10
          cpu.memory[cpu.i+2] = (vx % 100) % 10
          cpu.pc += 2
        case 0x0055:
          x := getX(instruction)
          var j uint8
          for j = 0; j <= x; j++ {
            cpu.memory[cpu.i+uint16(j)] = cpu.getRegister(j)
          }
          cpu.pc    += 2
        case 0x0065:
          x := getX(instruction)
          var j uint8
          for j = 0; j <= x; j++ {
            cpu.setRegister(j, cpu.memory[cpu.i+uint16(j)])
          }
          cpu.pc    += 2
      }
  }
}

func (cpu *CPU) getRegister(register uint8) uint8 {
  if register > 15 {
    panic(errors.New("Only have 16 registers, asking for register index ?"))
  }
  return cpu.v[register]
}

func (cpu *CPU) getKey() uint8 {
  for i, x := range cpu.key {
    if x {
      return uint8(i)
    }
  }
  return 0xFF
}

func (cpu *CPU) keyPressed(key uint8) bool {
  if key > 15 {
    return false
  }
  return cpu.key[key]
}

func (cpu *CPU) setRegister(register uint8, value uint8) {
  if register > 15 {
    panic(errors.New("Only have 16 registers, asking for register index ?"))
  }
  cpu.v[register] = value
}
