import re
import sys

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
    self.fired = False
  
  def is_done(self):
    return self.output.has_value()
  
  def is_ready(self):
    return self.input1.has_value() and self.input2.has_value() and not self.output.has_value()
  
  def reset(self):
    self.fired = False
    self.output.set_value(None)

  def __str__(self):
    return f"{type(self).__name__}({self.input1.name}, {self.input2.name}, {self.output.name})"

  def execute(self):
    pass

  def __hash__(self):
    return hash((type(self).__name__, self.input1.name, self.input2.name, self.output.name))

class AndGate(Gate):
  def __init__(self, input1, input2, output):
    super().__init__(input1, input2, output)
  
  def execute(self):
    self.fired = True
    self.output.set_value(self.input1.get_value() & self.input2.get_value())

class OrGate(Gate):
  def __init__(self, input1, input2, output):
    super().__init__(input1, input2, output)

  def execute(self):
    self.fired = True
    self.output.set_value(self.input1.get_value() | self.input2.get_value())

class XorGate(Gate):
  def __init__(self, input1, input2, output):
    super().__init__(input1, input2, output)
  
  def execute(self):
    self.fired = True
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

def wires_to_number(all_wires, prefix):
  the_wires = sorted(list(filter(lambda w: w.name.startswith(prefix), all_wires)), key=lambda w: w.name, reverse=True)
  v = ''.join([str(w.get_value()) for w in the_wires])
  v_int = int(v, base=2)
  return v_int

def set_wires_from_number(wires, prefix, new_number):
  the_wires = sorted(list(filter(lambda w: w.name.startswith(prefix), wires)), key=lambda w: w.name)
  for bit in range(len(the_wires)):
    the_wires[bit].set_value((new_number >> bit) & 1)
  assert new_number == wires_to_number(wires, prefix)

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
      ret.update(find_gates_downstream_of_wire(g.input1, gates))
      ret.update(find_gates_downstream_of_wire(g.input2, gates))
  return ret

def try_computer(wires, gates, input1_int, input2_int, expectedFunction):
  set_wires_from_number(wires.values(), "x", input1_int)
  set_wires_from_number(wires.values(), "y", input2_int)
  for g in gates:
    g.reset()
  bad_output_bits = []
  output_int = execute_computer(wires, gates)
  expected = expectedFunction(input1_int, input2_int)
  gates_fired = [g for g in gates if g.fired]
  if output_int != expected:
    for i in range(64):
      if (output_int >> i) & 1 != (expected >> i) & 1:
        bad_output_bits.append(i)
    #print(f"Got {bin(input1_int)} + {bin(input2_int)} = {bin(output_int)} instead of {bin(expected)}}")
    return False, output_int, bad_output_bits, gates_fired

  return True, output_int, [], gates_fired

def is_one_bit_working(wires, gates, bit, expectedFunction):
  known_bad_output_bits = []
  all_gates_used = []
  result, out, bad_bits_for_zero_add, gates_zero_add = try_computer(wires, gates, 0, 0, expectedFunction)
  #if not result:
  #  print(f"0b0 + 0b0 = {out}. bad bits: {bad_bits_for_zero_add}")

  result, out, bad_bits_for_one_zero, gates_one_zero = try_computer(wires, gates, 1 << bit, 0, expectedFunction)
  #if not result:
  #  print(f"{1 << bit} + 0b0 = {out}. bad bits: {bad_bits_for_one_zero}")

  result, out, bad_bits_for_zero_one, gates_zero_one = try_computer(wires, gates, 0, 1 << bit, expectedFunction)
  #if not result:
  #  print(f"0b0 + {1 << bit} = {out}. bad bits: {bad_bits_for_zero_one}")

  result, out, bad_bits_for_one_one, gates_one_one = try_computer(wires, gates, 1 << bit, 1 << bit, expectedFunction)
  #if not result:
  #  print(f"{1 << bit} + {1 << bit} = {out}. bad bits: {bad_bits_for_one_one}")

  known_bad_output_bits += bad_bits_for_zero_add
  known_bad_output_bits += bad_bits_for_one_zero
  known_bad_output_bits += bad_bits_for_zero_one
  known_bad_output_bits += bad_bits_for_one_one

  return list(set(known_bad_output_bits))

def swap_outputs(g1, g2):
  temp = g1.output
  g1.output = g2.output
  g2.output = temp
  g1.reset()
  g2.reset()

def has_cycle(gates):
  wires_to_gates = {}
  for g in gates:
    if g.input1 not in wires_to_gates:
      wires_to_gates[g.input1] = set()
    if g.input2 not in wires_to_gates:
      wires_to_gates[g.input2] = set()
    wires_to_gates[g.input1].add(g)
    wires_to_gates[g.input2].add(g)
  for starting_point in gates:
    to_visit = set()
    to_visit.add(starting_point)
    first = True
    visited = set()
    while len(to_visit) > 0:
      current = to_visit.pop()
      if not first:
        if current == starting_point:
          return True
      if current in visited:
        continue
      first = False
      gates_with_current_as_input = wires_to_gates.get(current.output, set())
      to_visit.update(gates_with_current_as_input)
      visited.add(current)
  return False

class SwapSequence(object):
  def __init__(self):
    self.swaps = []

  def copy_and_add_swap(self, g1, g2):
    new_seq = SwapSequence()
    new_seq.swaps = self.swaps.copy()
    new_seq.swaps.append((g1, g2))
    return new_seq
  
  def is_valid(self, g1, g2):
    for s in self.swaps:
      if g1 in s or g2 in s:
        return False
    return True
  
  def apply(self):
    for p in self.swaps:
      swap_outputs(p[0], p[1])
  
  def __str__(self):
    return ','.join([f"{str(p[0])}->{str(p[1])}" for p in self.swaps])

def try_fixing_computer(wires, gates, bit, upstream_of_bad, downstream_of_bad, expectedFunction, swap_seqs):
  print(f"Trying to fix bit {bit} by swapping from {len(upstream_of_bad)} and {len(downstream_of_bad)} gates") 

  assert not has_cycle(gates)

  new_swap_seqs = []

  counter = 0
  for g1 in upstream_of_bad:
    for g2 in downstream_of_bad:
      counter += 1
      if counter % 1000 == 0:
        print(f"on iteration {counter}")
      if g1 == g2:
        continue
      #print(f"Swapping {g1} and {g2}")

      swap_outputs(g1, g2)
      if has_cycle(gates):
        #print("has cycle, skipping")
        swap_outputs(g1, g2)
        continue
      fixed = len(is_one_bit_working(wires, gates, bit, expectedFunction)) == 0
      is_fully_working_now = False
      if len(swap_seqs) == 0:
        is_fully_working_now, _ = is_computer_working(wires, gates, expectedFunction, bit)
        if is_fully_working_now:
          swap_seq = SwapSequence()
          swap_seq = swap_seq.copy_and_add_swap(g1, g2)
          new_swap_seqs.append(swap_seq)
          print(f"Swapping {g1} and {g2} fixed the computer up to bit {bit}")
      else:
        for seq in swap_seqs:
          if seq.is_valid(g1, g2):
            seq.apply()
            fully_better, _ = is_computer_working(wires, gates, expectedFunction, bit)
            seq.apply()

            if fully_better:
              new_swap_seq = seq.copy_and_add_swap(g1, g2)
              new_swap_seqs.append(new_swap_seq)
              print(f"Swapping {g1} and {g2} fixed the computer up to bit {bit} as part of SwapSeq {str(new_swap_seq)}")
              is_fully_working_now = True
              break
            else:
              #print(f"Swapping {g1} and {g2} did not fix the computer up to bit {bit} when we also swapped {str(seq)}")
              pass
  
      swap_outputs(g1, g2)
  
  return len(new_swap_seqs) > 0, new_swap_seqs

def is_computer_working(wires, gates, expectedFunction, up_to_bit=None):
  num_input_bits = len(list(filter(lambda w: w.name.startswith("x"), wires.values())))
  known_bad_input_bits = []

  for bit in range(num_input_bits):
    if up_to_bit is not None and bit > up_to_bit:
      continue
    bad_output_bits = is_one_bit_working(wires, gates, bit, expectedFunction)
    if len(bad_output_bits) > 0:
      known_bad_input_bits.append(bit)

  return len(known_bad_input_bits) == 0, known_bad_input_bits

def fix_computer(wires, gates, expectedFunction, known_bad_input_bits):
  working_swap_seqs = []
  computer_working = True
  still_bad_input_bits = []
  for bit in known_bad_input_bits:
    swap_seqs_that_already_fixed_this_bit = []
    for ss in working_swap_seqs:
      ss.apply()
      bad_output_bits = is_one_bit_working(wires, gates, bit, expectedFunction)
      ss.apply()
      if len(bad_output_bits) == 0:
        print(f"Bit {bit} is already fixed by {ss}")
        swap_seqs_that_already_fixed_this_bit.append(ss)
    if len(swap_seqs_that_already_fixed_this_bit) > 0:
      working_swap_seqs = swap_seqs_that_already_fixed_this_bit
      continue

    # need to recompute here because earlier swaps maybe/likely fixed this gate.
    bad_output_bits = is_one_bit_working(wires, gates, bit, expectedFunction)
    if len(bad_output_bits) > 0:

      # Find all gates upstream of x_bit and y_bit and downstream of z_bit for all bad bits.
      # try swapping all outputs one at a time.
      to_consider = set()
      #print(f"Considering paths from input bit {bit} to output bits {bad_output_bits}")

      upstreams = set()
      upstreams.update(find_gates_upstream_of_wire(wires[f"x{bit:02}"], gates))
      upstreams.update(find_gates_upstream_of_wire(wires[f"y{bit:02}"], gates))
      downstreams = set()
      for bb in bad_output_bits:
        #if f"x{bb:02}" in wires:
        #  to_consider.update(find_gates_upstream_of_wire(wires[f"x{bb:02}"], gates).union(find_gates_upstream_of_wire(wires[f"y{bb:02}"], gates)))
        downstreams.update(find_gates_downstream_of_wire(wires[f"z{bb:02}"], gates))
      #consider = upstreams.union(downstreams)
      fixed, seqs = try_fixing_computer(wires, gates, bit, list(upstreams), list(downstreams), expectedFunction, working_swap_seqs)

      if fixed:
        print("Fixed and verified!")
        working_swap_seqs = seqs
      else:
        still_bad_input_bits.append(bit)
        computer_working = False

  still_bad_input_bits = list(set(still_bad_input_bits))
  working_swap_seqs = sorted(list(set(working_swap_seqs)), key=lambda p: str(p))
  print(f"Still known bad bits: {still_bad_input_bits}" )
  print(f"Swapped gates: {','.join([str(p) for p in working_swap_seqs])}")
  return computer_working, working_swap_seqs

def part2(fname, expectedFunction):
  wires, gates = parse_file(fname)
  #classify_gates(wires, gates)
  #return

  def get_original_output(wires, gates):
    in1 = wires_to_number(wires.values(), "x")
    in2 = wires_to_number(wires.values(), "y")
    return in1, in2, execute_computer(wires, gates)
    
  input1_int, input2_int, output_int = get_original_output(wires, gates)
  print(f"{input1_int} + {input2_int} = {expectedFunction(input1_int, input2_int)}. We get {output_int}.")
  print(f"{bin(input1_int)} + {bin(input2_int)} = {bin(expectedFunction(input1_int, input2_int))}. We get {bin(output_int)}.")

  is_working, bad_input_bits = is_computer_working(wires, gates, expectedFunction)
  print(f"Bad input bits: {bad_input_bits}")
  if not is_working:
    pass
    fixed, swap_seqs = fix_computer(wires, gates, expectedFunction, bad_input_bits)
    #fixed = False
    #swapped_gates = []
  else:
    print("It was already working")

  if fixed:
    print(f"Computer is working now that we swapped {[str(p) for p in swap_seqs]}")
    for ss in swap_seqs:
      for g in gates:
        g.reset()
      ss.apply()
      input1_int, input2_int, output_int = get_original_output(wires, gates)
      assert expectedFunction(input1_int, input2_int) == output_int
      ss.apply()
      swapped_outputs = []
      for p in ss.swaps:
        swapped_outputs.append(p[0].output.name)
        swapped_outputs.append(p[1].output.name)
      print(f"Working Sequence: {','.join(sorted(swapped_outputs))}")
  else:
    print(f"Computer is not working even though we swapped {[str(p) for p in swap_seqs]}")

part2("day24/day24-input-easy3.txt", (lambda x, y: x & y))
part2("day24/day24-input.txt", (lambda x, y: x + y))