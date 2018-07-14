package cpu

import (
  "testing"
  "fmt"
)

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
  fmt.Println("add v0 v1")

}
