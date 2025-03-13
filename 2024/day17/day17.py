import re
import sys

class  Computer(object):
  def __init__(self, a, b, c):
    self.orig_reg_a = a
    self.orig_reg_b = b
    self.orig_reg_c = c
    self.reg_a = a
    self.reg_b = b
    self.reg_c = c
    self.reg_a_sym = "A"
    self.reg_b_sym = "0"
    self.reg_c_sym = "0"
    self.instruction_pointer = 0
    self.outputs = []
    self.output_syms = []
  
  def reset_with_a(self, a):
    self.reg_a = a
    self.reg_b = self.orig_reg_b
    self.reg_c = self.orig_reg_c
    self.outputs = []
    self.instruction_pointer = 0

  def operand_value(self, operand):
    if operand >= 0 and operand <= 3:
      return operand, str(operand)
    elif operand == 4:
      return self.reg_a, self.reg_a_sym
    elif operand == 5:
      return self.reg_b, self.reg_b_sym
    elif operand == 6:
      return self.reg_c, self.reg_c_sym
    else:
      print(f"Invalid operand: {operand}")
      return  None

  def adv(self, operand):
    numerator = self.reg_a
    denominator, denom_sym = self.operand_value(operand)
    self.reg_a = numerator // 2**denominator
    self.reg_a_sym = f"(({self.reg_a_sym}) >> {denom_sym})"
    self.instruction_pointer += 2
  
  def bdv(self, operand):
    numerator = self.reg_a
    denominator, denom_sym = self.operand_value(operand)
    self.reg_b = numerator // 2**denominator
    self.reg_b_sym = f"(({self.reg_a_sym}) >> {denom_sym})"
    self.instruction_pointer += 2
    
  def cdv(self, operand):
    numerator = self.reg_a
    denominator, denom_sym = self.operand_value(operand)
    self.reg_c = numerator // 2**denominator
    self.reg_c_sym = f"(({self.reg_a_sym}) >> {denom_sym})"
    self.instruction_pointer += 2
  
  # The bxl instruction (opcode 1) calculates the bitwise XOR of register B and the
  # instruction's literal operand, then stores the result in register B.
  def bxl(self, operand):
    self.reg_b = self.reg_b ^ operand
    self.reg_b_sym = f"(({self.reg_b_sym}) ^ {operand})"
    self.instruction_pointer += 2

  # The bst instruction (opcode 2) calculates the value of its combo operand modulo 8
  # (thereby keeping only its lowest 3 bits), then writes that value to the B register.
  def bst(self, operand):
    op_value, op_value_sym = self.operand_value(operand)
    self.reg_b = op_value % 8
    self.reg_b_sym = f"(({op_value_sym}) % 8)"
    self.instruction_pointer += 2
  
  # The jnz instruction (opcode 3) does nothing if the A register is 0. However, if the
  # A register is not zero, it jumps by setting the instruction pointer to the value of
  # its literal operand; if this instruction jumps, the instruction pointer is not
  # increased by 2 after this instruction.
  def jnz(self, operand):
    if self.reg_a == 0:
      self.instruction_pointer += 2
    else:
      self.instruction_pointer = operand

  # The bxc instruction (opcode 4) calculates the bitwise XOR of register B and register C,
  # then stores the result in register B. (For legacy reasons, this instruction reads an
  # operand but ignores it.)
  def bxc(self, operand):
    self.reg_b = self.reg_b ^ self.reg_c
    self.reg_b_sym = f"(({self.reg_b_sym}) ^ ({self.reg_c_sym}))"
    self.instruction_pointer += 2

  # The out instruction (opcode 5) calculates the value of its combo operand modulo 8, then
  # ouputs that value. (If a program outputs multiple values, they are separated by commas.)
  def out(self, operand):
    op_value, op_value_sym = self.operand_value(operand)
    self.outputs.append(op_value % 8)
    self.output_syms.append(f"{op_value_sym} % 8")
    self.instruction_pointer += 2

  def print_registers(self):
    print("A: " + str(self.reg_a) + ", B: " + str(self.reg_b) + ", C: " + str(self.reg_c))

  def run_instruction(self, opcode, operand):
    #self.print_registers()
    #print("Instruction: " + str(opcode) + " " + str(operand))
    if opcode == 0:
      self.adv(operand)
    elif opcode == 1:
      self.bxl(operand)
    elif opcode == 2:
      self.bst(operand)
    elif opcode == 3:
      self.jnz(operand)
    elif opcode == 4:
      self.bxc(operand)
    elif opcode == 5:
      self.out(operand)
    elif opcode == 6:
      self.bdv(operand)
    elif opcode == 7:
      self.cdv(operand)

    #self.print_registers()

  def run_program(self, program):
    while self.instruction_pointer < len(program):
      opcode = program[self.instruction_pointer]
      operand = program[self.instruction_pointer + 1]
      self.run_instruction(opcode, operand)
    return self.outputs, self.output_syms
  
def parse(fname):
  reg_a_re = re.compile(r"Register A: (\d+)")
  reg_b_re = re.compile(r"Register B: (\d+)")
  reg_c_re = re.compile(r"Register C: (\d+)")
  program_re = re.compile(r"Program: (.+)")
  with open(fname, "r") as file:
    reg_a = int(reg_a_re.match(file.readline().strip()).group(1))
    reg_b = int(reg_b_re.match(file.readline().strip()).group(1))
    reg_c = int(reg_c_re.match(file.readline().strip()).group(1))
    file.readline()
    program = [int(x) for x in program_re.match(file.readline().strip()).group(1).split(",")]
    return reg_a, reg_b, reg_c, program
  
  
def day17_part1(fname):
  reg_a, reg_b, reg_c, program = parse(fname)
  comp = Computer(reg_a, reg_b, reg_c)
  print(','.join([str(o) for o in comp.run_program(program)[0]]))
  out, out_syms = comp.run_program(program)
  print(','.join([str(o) for o in out]))
  #for i in range(len(out)):
    #print(f"Output {i}: {bin(out[i])} === {out_syms[i]}")

def day17_part2(fname):
  reg_a, reg_b, reg_c, program = parse(fname)

  program_to_bits = 0
  for i in program:
    program_to_bits = (program_to_bits << 3) | i

  reg_a = 2**(len(bin(program_to_bits)) - 3)
  print(f"Starting at {reg_a}")


  #print(f"{bin(program_to_bits)} {len(bin(program_to_bits))}")
  
  while True:
    if reg_a % 100_000 == 0:
      print("Trying reg_a: " + str(reg_a))
    comp = Computer(reg_a, reg_b, reg_c)
    new_program, _ = comp.run_program(program)


    #new_program_to_bits = 0
    #for i in new_program:
    #  new_program_to_bits = (new_program_to_bits << 3) | i
    
    #print(f"{bin(reg_a)} {len(bin(reg_a))} --> {bin(new_program_to_bits)} {len(bin(new_program_to_bits))}")

    if new_program == program:
      print("Found program that repeats itself: " + str(reg_a))
      break
    reg_a += 1


#day17_part1("day17/day17-input-easy.txt")
day17_part1("day17/day17-input.txt")

day17_part2("day17/day17-input-easy2.txt")
#day17_part2("day17/day17-input.txt")



# 2,4 # B = A & 0b111. B is bottom three bits of A
# 1,3 # B = B ^ 3 = (A&0b111) ^ 0b11
# 7,5 # C = A >> B = A >> (A&0b111 ^ 0b11)
# 1,5 # B = B ^ 5 = (A&0b111) ^ 0b11 ^ 0b101
# 0,3 # A = A >> 3
# 4,1 # B = B ^ C = (A&0b111) ^ 0b11 ^ 0b101 ^ (A >> (A&0b111 ^ 0b11))
# 5,5 # output B & 0b111
# 3,0 # if A != 0, jump to 0. we know this jump will happen exactly 15 times

def one_pass(a_1):
  b_1 = a_1 & 0b111
  b_2 = b_1 ^ 0b11
  c_1 = a_1 >> b_2
  b_3 = b_2 ^ 0b101
  #a_2 = a_1 >> 3
  b_4 = b_3 ^ c_1
  b_5 = b_4 & 0b111
  output = b_5
  return output

best_i = -1

def search_all_passes(program, goal, i, input, reference_computer):
  global best_i

  if i >= len(goal):
    print(f"Found a winner! {input}")

    reference_computer.reset_with_a(input)
    output, _ = reference_computer.run_program(program)
    print(f"Output from reference computer (0): {output}")

    sys.exit(1)

    return [input]

  wins = []
  for bottom_3_bits_of_a in range(8):
    output = one_pass((input << 3) | bottom_3_bits_of_a)
    if output == goal[i]: 
      #print(f"success: {input} << 3 | {bottom_3_bits_of_a} = {output}")
      win = (input << 3) | bottom_3_bits_of_a
      reference_computer.reset_with_a(win)
      output, _ = reference_computer.run_program(program)
      #print(f"Output from reference computer (0) for first {i}: {win} --> {output[:i+1]}")
      wins.append(win)
  if len(wins) == 0:
    #print(f"nothing worked to get ouput={goal[:i+1]}")
    return []
  if i > best_i:
    #print(f"these ones work to get output={goal[:i+1]}. inputs={wins}")
    best_i = i
  ret = []
  for winner in wins:
    #print(f"trying {winner}")
    next_wins = search_all_passes(program, goal, i+1, winner, reference_computer)
    for nw in next_wins:
      ret.append(next_wins)
  return ret

def part2(fname):
  reg_a, reg_b, reg_c, program = parse(fname)
  reference_computer = Computer(reg_a, reg_b, reg_c)
  goal = program.copy()
  goal.reverse()
  starting_a = 0
#  while True:
  wins = search_all_passes(program, goal, 0, starting_a, reference_computer)
#    if len(wins) > 0:
#      break
#    starting_a += 1
  print(f"Found {len(wins)} winning inputs: {wins}")

part2("day17/day17-input.txt")