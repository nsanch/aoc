import re

class Wire(object):
  def __init__(self, name, initial_value):
    self.name = name
    self.value = initial_value
  
  def has_value(self):
    return self.value is not None
  
  def get_value(self):
    return self.value
  
  def set_value(self, value):
    #print(f"Setting wire {self.name} to {value}")
    self.value = value

  def __str__(self):
    return f"Wire({self.name}, {self.value})"

class Gate(object):
  def __init__(self, input1, input2, output):
    self.input1 = input1
    self.input2 = input2
    self.output = output
  
  def is_done(self):
    return self.output.has_value()
  
  def is_ready(self):
    return self.input1.has_value() and self.input2.has_value() and not self.output.has_value()
  
  def reset(self):
    self.output.set_value(None)

  #def swap_outputs(self, other):

  def execute(self):
    pass

  def __hash__(self):
    return hash((type(self).__name__, self.input1.name, self.input2.name, self.output.name))

class AndGate(Gate):
  def __init__(self, input1, input2, output):
    super().__init__(input1, input2, output)
  
  def execute(self):
    self.output.set_value(self.input1.get_value() & self.input2.get_value())

class OrGate(Gate):
  def __init__(self, input1, input2, output):
    super().__init__(input1, input2, output)

  def execute(self):
    self.output.set_value(self.input1.get_value() | self.input2.get_value())

class XorGate(Gate):
  def __init__(self, input1, input2, output):
    super().__init__(input1, input2, output)
  
  def execute(self):
    self.output.set_value(self.input1.get_value() ^ self.input2.get_value())

def parse_file(fname) -> tuple[dict[str, Wire], list[Gate]]:
  wires = {}
  with open(fname) as f:
    while True:
      line = f.readline().strip()
      if line == '':
        break
      name, value = line.split(": ")
      wires[name] = Wire(name, int(value))
    
    gates = []
    gate_re = re.compile(r"(.*) (AND|XOR|OR) (.*) -> (.*)")
    while True:
      line = f.readline().strip()
      m = gate_re.match(line)
      if m is None:
        break
      if m.group(1) not in wires:
        wires[m.group(1)] = Wire(m.group(1), None)
      if m.group(3) not in wires:
        wires[m.group(3)] = Wire(m.group(3), None)
      if m.group(4) not in wires:
        wires[m.group(4)] = Wire(m.group(4), None)

      input1 = wires[m.group(1)]
      input2 = wires[m.group(3)]
      output = wires[m.group(4)]
      operation = m.group(2)
      if operation == "AND":
        gates.append(AndGate(input1, input2, output))
      elif operation == "XOR":
        gates.append(XorGate(input1, input2, output))
      elif operation == "OR":
        gates.append(OrGate(input1, input2, output))

  return wires, gates

def make_model_computer(bits):
  wires = {}
  gates = []
  for i in range(bits):
    wires[f"x{i:02}"] = Wire(f"x{i:02}", None)
    wires[f"y{i:02}"] = Wire(f"y{i:02}", None)
    wires[f"z{i:02}"] = Wire(f"z{i:02}", None)

    # https://en.wikipedia.org/wiki/Adder_(electronics)

    # direct sum of x + y. for bit z this goes straight to z00 because there's no carry bit.
    wires[f"s{i:02}"] = Wire(f"s{i:02}", None)
    if i == 0:
      gates.append(XorGate(wires[f"x{i:02}"], wires[f"y{i:02}"], wires[f"z{i:02}"]))
    else:
      # sum = x ^ y ^ c
      gates.append(XorGate(wires[f"x{i:02}"], wires[f"y{i:02}"], wires[f"s{i:02}"]))
      gates.append(XorGate(wires[f"s{i:02}"], wires[f"c{i:02}"], wires[f"z{i:02}"]))

    #c_1 is the carry bit of x_0 + y_0
    wires[f"c{(i+1):02}"] = Wire(f"c{(i+1):02}", None)
    gates.append(AndGate(wires[f"x{i:02}"], wires[f"y{i:02}"], wires[f"c{(i+1):02}"]))

  wires[f"z{bits:02}"] = Wire(f"z{bits:02}", None)

  return wires, gates

def wires_to_number(all_wires, prefix):
  the_wires = sorted(list(filter(lambda w: w.name.startswith(prefix), all_wires)), key=lambda w: w.name, reverse=True)
  #print(','.join(map(str, the_wires)))
  v = ''.join([str(w.get_value()) for w in the_wires])
  v_int = int(v, base=2)
  return v_int

def set_wires_from_number(wires, prefix, new_number):
  the_wires = sorted(list(filter(lambda w: w.name.startswith(prefix), wires)), key=lambda w: w.name)
  for bit in range(len(the_wires)):
    the_wires[bit].set_value((new_number >> bit) & 1)
  assert new_number == wires_to_number(wires, prefix)
  #print(f"Setting {the_wires[bit].name} to {(new_number >> bit) & 1}")

def execute_computer(wires, gates):
  waiting_gates = list(filter(lambda g: not g.is_done(), gates))
  while len(waiting_gates) > 0:
    for gate in waiting_gates:
      if gate.is_ready():
        gate.execute()
        waiting_gates.remove(gate)

  output_int = wires_to_number(wires.values(), "z")
  return output_int

def part1(fname):
  wires, gates = parse_file(fname)
  return execute_computer(wires, gates)

print(part1("day24/day24-input-easy.txt"))
print(part1("day24/day24-input-easy2.txt"))
print(part1("day24/day24-input.txt"))

def find_gates_upstream_of_wire(wire, gates):
  ret = set()
  for g in gates:
    if g.input1 == wire or g.input2 == wire:
      ret.add(g)
      ret.update(find_gates_upstream_of_wire(g.output, gates))
  return ret

def find_gates_downstream_of_wire(wire, gates):
  ret = set()
  for g in gates:
    if g.output == wire:
      ret.add(g)
      ret.update(find_gates_upstream_of_wire(g.input1, gates))
      ret.update(find_gates_upstream_of_wire(g.input2, gates))
  return ret

def try_computer(wires, gates, input1_int, input2_int):
  set_wires_from_number(wires.values(), "x", input1_int)
  set_wires_from_number(wires.values(), "y", input2_int)
  for g in gates:
    g.reset()
  bad_bits = []
  output_int = execute_computer(wires, gates)
  if output_int != (input1_int + input2_int):
    expected = input1_int + input2_int
    for i in range(64):
      if (output_int >> i) & 1 != (expected >> i) & 1:
        bad_bits.append(i)
    #print(f"Got {bin(input1_int)} + {bin(input2_int)} = {bin(output_int)} instead of {bin(input1_int + input2_int)}")
    return False, output_int, bad_bits

  return True, output_int, []

def is_one_bit_working(wires, gates, bit):
  known_bad_bits = []
  result, out, bad_bits_for_zero_add = try_computer(wires, gates, 0, 0)
  if not result:
    print(f"0b0 + 0b0 = {out}. bad bits: {bad_bits_for_zero_add}")

  result, out, bad_bits_for_one_zero = try_computer(wires, gates, 1 << bit, 0)
  if not result:
    print(f"{1 << bit} + 0b0 = {out}. bad bits: {bad_bits_for_one_zero}")
  
  result, out, bad_bits_for_one_zero = try_computer(wires, gates, 1 << bit, 0)
  if not result:
    print(f"{1 << bit} + 0b0 = {out}. bad bits: {bad_bits_for_one_zero}")

  result, out, bad_bits_for_zero_one = try_computer(wires, gates, 0, 1 << bit)
  if not result:
    print(f"0b0 + {1 << bit} = {out}. bad bits: {bad_bits_for_zero_one}")

  result, out, bad_bits_for_one_one = try_computer(wires, gates, 1 << bit, 1 << bit)
  if not result:
    print(f"{1 << bit} + {1 << bit} = {out}. bad bits: {bad_bits_for_one_one}")

  known_bad_bits += bad_bits_for_zero_add
  known_bad_bits += bad_bits_for_one_zero
  known_bad_bits += bad_bits_for_zero_one
  known_bad_bits += bad_bits_for_one_one

  return list(set(known_bad_bits))

def swap_outputs(g1, g2):
  temp = g1.output
  g1.output = g2.output
  g2.output = temp

def try_fixing_computer(wires, gates, bit, upstream_of_bad, downstream_of_bad):
  print(f"Trying to fix bit {bit} by swapping from {len(upstream_of_bad)} and {len(downstream_of_bad)} gates") 

  for g1 in upstream_of_bad:
    for g2 in downstream_of_bad:
      if g1 == g2:
        continue
      swap_outputs(g1, g2)
      if is_one_bit_working(wires, gates, bit):
        print(f"Swapping {g1.output} and {g2.output} fixed the computer")
        return True      
      swap_outputs(g1, g2)
  return False

def is_computer_working(wires, gates):
  num_input_bits = len(list(filter(lambda w: w.name.startswith("x"), wires.values())))
  #print([w.name for w in wires.values()])
  known_bad_bits = []
  computer_working = True
  for bit in range(num_input_bits):
    bad_bits = is_one_bit_working(wires, gates, bit)
    known_bad_bits += bad_bits
    if len(bad_bits) > 0:
      computer_working = False

      # Find all gates upstream of x_bit and y_bit and downstream of z_bit for all bad bits.
      # try swapping all outputs one at a time.
      bad_upstream = []
      bad_downstream = []
      for bb in bad_bits:
        bad_upstream += list(find_gates_upstream_of_wire(wires[f"x{bb:02}"], gates).union(find_gates_upstream_of_wire(wires[f"y{bb:02}"], gates)))
        bad_downstream += list(find_gates_downstream_of_wire(wires[f"z{bb:02}"], gates))
      try_fixing_computer(wires, gates, bit, list(set(bad_upstream)), list(set(bad_downstream)))

  known_bad_bits = list(set(known_bad_bits))
  print(f"Known bad bits: {known_bad_bits}" )
  return len(known_bad_bits) == 0

def all_possible_pairs_of_gates(gates):
  ret = []
  for i in range(len(gates)):
    for j in range(i+1, len(gates)):
      yield gates[i], gates[j]

def pairs_of_gates(gates, num):
  if num == 1:
    return all_possible_pairs_of_gates(gates)
  
  ret = []
  pairs = all_possible_pairs_of_gates(gates)
  sub_pairs = pairs_of_gates(gates, num-1)
  for sp in sub_pairs:
    for p in pairs:
      yield sp + [p]

def part2(fname):
  wires, gates = parse_file(fname)
  input1_int = wires_to_number(wires.values(), "x")
  input2_int = wires_to_number(wires.values(), "y")
  output_int = execute_computer(wires, gates)
  print(f"{input1_int} + {input2_int} = {input1_int + input2_int}. We get {output_int}.")
  print(f"{bin(input1_int)} + {bin(input2_int)} = {bin(input1_int + input2_int)}. We get {bin(output_int)}.")

  if is_computer_working(wires, gates):
    print("Computer is working")
  else:
    print("Computer is not working")
  #for set_of_pairs in pairs_of_gates(gates, 2):
  #  for g1, g2 in set_of_pairs:
      


#part2("day24/day24-input-easy3.txt")
#part2("day24/day24-input-easy.txt")
part2("day24/day24-input.txt")