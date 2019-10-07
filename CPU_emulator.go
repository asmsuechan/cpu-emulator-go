package main

import "fmt"

const (
	MOV = iota
	ADD
	SUB
	AND
	OR
	SL
	SR
	SRA
	LDL
	LDH
	CMP
	JE
	JMP
	LD
	ST
	HLT
)

const (
	REG0 = iota
	REG1
	REG2
	REG3
	REG4
	REG5
	REG6
	REG7
)

var reg [8]uint16
var rom [256]uint16
var ram [256]uint16

func main() {
	var pc, ir, flag_eq uint16

	assembler()

	pc = 0
	flag_eq = 0

	for op_code(ir) != HLT {
		ir = rom[pc]
		fmt.Printf("%5d %5x %5d %5d %5d %5d\n", pc, ir, reg[0], reg[1], reg[2], reg[3])
		pc += 1

		switch op_code(ir) {
		case MOV:
			reg[op_regA(ir)] = reg[op_regB(ir)]
		case ADD:
			reg[op_regA(ir)] = reg[op_regA(ir)] + reg[op_regB(ir)]
		case SUB:
			reg[op_regA(ir)] = reg[op_regA(ir)] - reg[op_regB(ir)]
		case AND:
			reg[op_regA(ir)] = reg[op_regA(ir)] & reg[op_regB(ir)]
		case OR:
			reg[op_regA(ir)] = reg[op_regA(ir)] | reg[op_regB(ir)]
		case SL:
			reg[op_regA(ir)] = reg[op_regA(ir)] << 1
		case SR:
			reg[op_regA(ir)] = reg[op_regA(ir)] >> 1
		case SRA:
			reg[op_regA(ir)] = (reg[op_regA(ir)] & 0x8000) | (reg[op_regA(ir)] >> 1)
		case LDL:
			reg[op_regA(ir)] = (reg[op_regA(ir)] & 0xff00) | (op_data(ir) & 0x00ff)
		case LDH:
			reg[op_regA(ir)] = (op_data(ir) << 8) | (reg[op_regA(ir)] & 0x00ff)
		case CMP:
			if reg[op_regA(ir)] == reg[op_regB(ir)] {
				flag_eq = 1
			} else {
				flag_eq = 0
			}
		case JE:
			if flag_eq == 1 {
				pc = op_addr(ir)
			}
		case JMP:
			pc = op_addr(ir)
		case LD:
			reg[op_regA(ir)] = ram[op_addr(ir)]
		case ST:
			ram[op_addr(ir)] = reg[op_regA(ir)]
		}
	}

	fmt.Println("ram[64] = ", ram[64])
}

// ROMに機械語命令を書き込む関数
func assembler() {
	rom[0] = ldh(REG0, 0)
	rom[1] = ldl(REG0, 0)
	rom[2] = ldh(REG1, 0)
	rom[3] = ldl(REG1, 1)
	rom[4] = ldh(REG2, 0)
	rom[5] = ldl(REG2, 0)
	rom[6] = ldh(REG3, 0)
	rom[7] = ldl(REG3, 10)
	rom[8] = add(REG2, REG1)
	rom[9] = add(REG0, REG2)
	rom[10] = st(REG0, 64)
	rom[11] = cmp(REG2, REG3)
	rom[12] = je(14)
	rom[13] = jmp(8)
	rom[14] = hlt()
}

func mov(ra, rb uint16) uint16 {
	return ((MOV << 11) | (ra << 8) | (rb << 5))
}

func add(ra, rb uint16) uint16 {
	return ((ADD << 11) | (ra << 8) | (rb << 5))
}

func sub(ra, rb uint16) uint16 {
	return ((SUB << 11) | (ra << 8) | (rb << 5))
}

func and(ra, rb uint16) uint16 {
	return ((AND << 11) | (ra << 8) | (rb << 5))
}

func or(ra, rb uint16) uint16 {
	return ((OR << 11) | (ra << 8) | (rb << 5))
}

func sl(ra uint16) uint16 {
	return ((SL << 11) | (ra << 8))
}

func sr(ra uint16) uint16 {
	return ((SR << 11) | (ra << 8))
}

func sra(ra uint16) uint16 {
	return ((SRA << 11) | (ra << 8))
}

func ldh(ra, ival uint16) uint16 {
	return ((LDH << 11) | (ra << 8) | (ival & 0x00ff))
}

func ldl(ra, ival uint16) uint16 {
	return ((LDL << 11) | (ra << 8) | (ival & 0x00ff))
}

func cmp(ra, rb uint16) uint16 {
	return ((CMP << 11) | (ra << 8) | (rb << 5))
}

func je(addr uint16) uint16 {
	return ((JE << 11) | (addr & 0x00ff))
}

func jmp(addr uint16) uint16 {
	return ((JMP << 11) | (addr & 0x00ff))
}

func ld(ra, addr uint16) uint16 {
	return ((LD << 11) | (ra << 8) | (addr & 0x00ff))
}

func st(ra, addr uint16) uint16 {
	return ((ST << 11) | (ra << 8) | (addr & 0x00ff))
}

func hlt() uint16 {
	return (HLT << 11)
}

func op_code(ir uint16) uint16 {
	return (ir >> 11)
}

func op_regA(ir uint16) uint16 {
	return ((ir >> 8) & 0x0007)
}

func op_regB(ir uint16) uint16 {
	return ((ir >> 5) & 0x0007)
}

func op_data(ir uint16) uint16 {
	return (ir & 0x00ff)
}

func op_addr(ir uint16) uint16 {

	return (ir & 0x00ff)
}
