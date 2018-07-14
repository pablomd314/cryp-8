package cpu

import (
  "testing"
  // "fmt"
)

func checkMem(cpu *CPU, addr uint16, val uint8, t *testing.T) {
  if cpu.memory[addr] != val {
    t.Errorf("Incorrect memory[%v]. Got %v, wanted %v", addr, cpu.memory[addr], val)
  }
}

func checkStimer(cpu *CPU, val uint8, t *testing.T) {
  if cpu.stimer != val {
    t.Errorf("Incorrect stimer. Got %v, wanted %v", cpu.stimer, val)

  }
}

func checkDtimer(cpu *CPU, val uint8, t *testing.T) {
  if cpu.dtimer != val {
    t.Errorf("Incorrect dtimer. Got %v, wanted %v", cpu.dtimer, val)

  }
}

func checkI(cpu *CPU, val uint16, t *testing.T) {
  if cpu.i != val {
    t.Errorf("Incorrect i. Got %v, wanted %v", cpu.i, val)

  }
}

func checkReg(cpu *CPU, reg int, val uint8, t *testing.T) {
  if cpu.v[reg] != val {
    t.Errorf("Incorrect v[%v]. Got %v, wanted %v", reg, cpu.v[reg], val)

  }
}

func checkStack(cpu *CPU, stack_idx int, pc uint16, t *testing.T) {
  if cpu.stack[stack_idx] != pc {
    t.Errorf("Incorrect stack[%v]. Got %v, wanted %v", stack_idx, cpu.stack[stack_idx], pc)
  }
}

func checkSP(cpu *CPU, sp uint8, t *testing.T) {
  if cpu.sp != sp {
    t.Errorf("Incorrect stack pointer. Got %v, wanted %v", cpu.sp, sp)
  }
}

func checkPC(cpu *CPU, pc uint16, t *testing.T) {
  if cpu.pc != pc {
    t.Errorf("Incorrect PC. Got %v, wanted %v", cpu.pc, pc)
  }
}

func TestFlowSubroutine(t *testing.T) {
  cpu := NewCPU()
  ogpc := cpu.pc

  cpu.executeInstruction(0x2abc)
  checkPC(&cpu, 0xabc, t)
  checkSP(&cpu, 1, t)
  checkStack(&cpu, 0, ogpc, t)

  cpu.executeInstruction(0x00ee)
  checkPC(&cpu, ogpc+2, t)
  checkSP(&cpu, 0, t)
}

func TestFlowSkip(t *testing.T) {
  cpu := NewCPU()
  ogpc := cpu.pc
  cpu.setRegister(0, 0xab)

  cpu.executeInstruction(0x30ab) /* SEQ V0 0xab */
  checkPC(&cpu, ogpc+4, t)

  ogpc = cpu.pc
  cpu.setRegister(0, 0xab)

  cpu.executeInstruction(0x30ff) /* SNE V0 0xff */
  checkPC(&cpu, ogpc+2, t) 

  ogpc = cpu.pc
  cpu.setRegister(0, 0xab)

  cpu.executeInstruction(0x40ab) /* SNE V0 0xab */
  checkPC(&cpu, ogpc+2, t)

  ogpc = cpu.pc
  cpu.setRegister(0, 0xab)

  cpu.executeInstruction(0x40ff) /* SNE V0 0xff */
  checkPC(&cpu, ogpc+4, t)

  ogpc = cpu.pc
  cpu.setRegister(0, 0xab)
  cpu.setRegister(1, 0xab)

  cpu.executeInstruction(0x5010) /* SEQ V0 V1 */
  checkPC(&cpu, ogpc+4, t)

  ogpc = cpu.pc
  cpu.setRegister(0, 0xab)
  cpu.setRegister(1, 0xbc)

  cpu.executeInstruction(0x5010) /* SEQ V0 V1 */
  checkPC(&cpu, ogpc+2, t)

  ogpc = cpu.pc
  cpu.setRegister(0, 0xab)
  cpu.setRegister(1, 0xab)

  cpu.executeInstruction(0x9010) /* SNE V0 V1 */
  checkPC(&cpu, ogpc+2, t)

  ogpc = cpu.pc
  cpu.setRegister(0, 0xab)
  cpu.setRegister(1, 0xbc)
  
  cpu.executeInstruction(0x9010) /* SNE V0 V1 */
  checkPC(&cpu, ogpc+4, t)
}

func TestFlowSkipKeys(t *testing.T) {
  cpu := NewCPU()
  var ogpc uint16
  var key uint8
  for key=0; key < 16; key++ {
    ogpc = cpu.pc
    cpu.setRegister(0, key) 
    cpu.key[key] = false
    cpu.executeInstruction(0xe09e)
    checkPC(&cpu, ogpc+2, t)

    ogpc = cpu.pc
    cpu.setRegister(0, key) 
    cpu.key[key] = true
    cpu.executeInstruction(0xe09e)
    checkPC(&cpu, ogpc+4, t)
  }
}

func TestFlowJump(t *testing.T) {
  cpu := NewCPU()
  cpu.executeInstruction(0x1abc)
  checkPC(&cpu, 0xabc, t)
  cpu.setRegister(0, 0xab)
  cpu.executeInstruction(0xbcde)
  checkPC(&cpu, 0xab + 0xcde, t)
}

func TestMathAdd(t *testing.T) {
  cpu := NewCPU()
  cpu.setRegister(0, 0x12)

  cpu.executeInstruction(0x7034) /* ADD V0 0x34 */
  checkReg(&cpu, 0, 0x12 + 0x34, t)
  checkReg(&cpu, 0xf, 0, t)

  cpu.setRegister(0, 0x12)

  cpu.executeInstruction(0x70ff) /* ADD V0 0xff */
  checkReg(&cpu, 0, 17, t)
  checkReg(&cpu, 0xf, 0, t)

  cpu.setRegister(0, 0x12)
  cpu.setRegister(1, 0x34)

  cpu.executeInstruction(0x8014) /* ADD V0 V1 */
  checkReg(&cpu, 0, 0x12 + 0x34, t)
  checkReg(&cpu, 0xf, 0, t)

  cpu.setRegister(0, 0x12)
  cpu.setRegister(1, 0xff)

  cpu.executeInstruction(0x8014) /* ADD V0 V1 */
  checkReg(&cpu, 0, 17, t)
  checkReg(&cpu, 0xf, 1, t)

  cpu.i = 0x123
  cpu.setRegister(0, 0x45)

  cpu.executeInstruction(0xf01e)
  checkI(&cpu, 0x123 + 0x45, t)
}

func TestMathSub(t *testing.T) {
  cpu := NewCPU()
  cpu.setRegister(0, 0x43)
  cpu.setRegister(1, 0x21)

  cpu.executeInstruction(0x8015) // sub v0 v1
  checkReg(&cpu, 0, 0x43-0x21, t)
  checkReg(&cpu, 0xf, 0, t)

  cpu.setRegister(0, 0x12)
  cpu.setRegister(1, 0x34)

  cpu.executeInstruction(0x8015) // sub v0 v1
  checkReg(&cpu, 0, 222, t)
  checkReg(&cpu, 0xf, 1, t)

  cpu.setRegister(0, 0x34)
  cpu.setRegister(1, 0x12)

  cpu.executeInstruction(0x8017) // sub v0 v1
  checkReg(&cpu, 0, 222, t)
  checkReg(&cpu, 0xf, 1, t)

  cpu.setRegister(0, 0x12)
  cpu.setRegister(1, 0x34)

  cpu.executeInstruction(0x8017) // sub v0 v1
  checkReg(&cpu, 0, 0x34-0x12, t)
  checkReg(&cpu, 0xf, 0, t)
}

func TestMathBitwise(t *testing.T) {
  cpu := NewCPU()
  cpu.setRegister(0, 0x12)
  cpu.setRegister(1, 0x34)
  
  cpu.executeInstruction(0x8011)
  checkReg(&cpu, 0, 0x34|0x12, t)

  cpu.setRegister(0, 0x12)
  cpu.setRegister(1, 0x34)
  
  cpu.executeInstruction(0x8012)
  checkReg(&cpu, 0, 0x34&0x12, t)

  cpu.setRegister(0, 0x12)
  cpu.setRegister(1, 0x34)
  
  cpu.executeInstruction(0x8013)
  checkReg(&cpu, 0, 0x34^0x12, t)
}

func TestLoadReg(t *testing.T) {
  cpu := NewCPU()

  cpu.executeInstruction(0x60ab)
  checkReg(&cpu, 0, 0xab, t)  

  cpu.setRegister(0, 0)
  cpu.setRegister(1, 0xff)

  cpu.executeInstruction(0x8010)
  checkReg(&cpu, 0, 0xff, t)

  cpu.i = 0x200

  cpu.executeInstruction(0xafff)
  checkI(&cpu, 0xfff, t)
}

func TestLoadTimers(t *testing.T) {
  cpu := NewCPU()
  cpu.setRegister(0, 0xab)

  cpu.executeInstruction(0xf015)
  checkDtimer(&cpu, 0xab, t)

  cpu.executeInstruction(0xf018)
  checkStimer(&cpu, 0xab, t)

  cpu.setRegister(0,0)
  
  cpu.executeInstruction(0xf007)
  checkReg(&cpu, 0, 0xab, t)
}

func TestLoadBcd(t *testing.T) {
  cpu := NewCPU()
  cpu.setRegister(0, 0xab)

  cpu.executeInstruction(0xf033)
  checkMem(&cpu, cpu.i+0, 0xab/100, t)
  checkMem(&cpu, cpu.i+1, (0xab % 100)/10, t)
  checkMem(&cpu, cpu.i+2, (0xab % 10), t) 
}

func TestLoadMem(t *testing.T) {
  cpu := NewCPU()
  cpu.setRegister(0, 0xde)
  cpu.setRegister(1, 0xad)
  cpu.setRegister(2, 0xbe)
  cpu.setRegister(3, 0xef)
  cpu.i = 0x123

  cpu.executeInstruction(0xf355)
  checkMem(&cpu, cpu.i+0, 0xde, t)
  checkMem(&cpu, cpu.i+1, 0xad, t)
  checkMem(&cpu, cpu.i+2, 0xbe, t)
  checkMem(&cpu, cpu.i+3, 0xef, t)

  cpu.setRegister(0, 0)
  cpu.setRegister(1, 0)
  cpu.setRegister(2, 0)
  cpu.setRegister(3, 0)

  cpu.executeInstruction(0xf365)
  checkReg(&cpu, 0, 0xde, t)
  checkReg(&cpu, 1, 0xad, t) 
  checkReg(&cpu, 2, 0xbe, t) 
  checkReg(&cpu, 3, 0xef, t) 
}

