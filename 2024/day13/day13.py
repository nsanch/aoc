import re
import math
import sys
import numpy as np

class Button(object):
  def __init__(self, x, y):
    self.x = x
    self.y = y

class PrizePosition(object):
  def __init__(self, x, y):
    self.x = x
    self.y = y

class Machine(object):
  global_counter = 1
  def __init__(self, button_a, button_b, prize_position):
    self.counter = Machine.global_counter
    Machine.global_counter += 1
    self.button_a = button_a
    self.button_b = button_b
    self.prize_position = prize_position

  def __str__(self):
    return f"Machine {self.counter} Button A: {self.button_a.x},{self.button_a.y} Button B: {self.button_b.x},{self.button_b.y} Prize: {self.prize_position.x},{self.prize_position.y}"

button_re = re.compile(r'Button .: X\+(\d+), Y\+(\d+)')
prize_re = re.compile(r'Prize: X=(\d+), Y=(\d+)')

def consume_machine(f, increment=0):
  button_a = f.readline().strip()
  if button_a == '':
    return None
  a_match = button_re.match(button_a)
  button_a = Button(int(a_match.group(1)), int(a_match.group(2)))
  button_b = f.readline().strip()
  b_match = button_re.match(button_b)
  button_b = Button(int(b_match.group(1)), int(b_match.group(2)))
  position = f.readline().strip()
  p_match = prize_re.match(position)
  prize_position = PrizePosition(int(p_match.group(1)) + increment, int(p_match.group(2)) + increment)
  empty = f.readline().strip()
  return Machine(button_a, button_b, prize_position)

# i'm a dum-dum, the problem is just solving a system of two equations with two unknowns.
# numpy does this with linear algebra built-in
def starting_over(machine: Machine) -> tuple[int, int]:
  # Define the coefficient matrix
  A = np.array([[machine.button_a.x, machine.button_b.x], [machine.button_a.y, machine.button_b.y]]) 

  # Define the constant vector
  B = np.array([machine.prize_position.x, machine.prize_position.y])

  # Solve the system of equations
  X = np.linalg.solve(A, B)

  # Print the solution
  #print(X) # Output: [x, y]
  
  a_presses = np.round(X[0])
  b_presses = np.round(X[1])

  x_pos = a_presses*machine.button_a.x + b_presses*machine.button_b.x
  y_pos = a_presses*machine.button_a.y + b_presses*machine.button_b.y

  if x_pos == machine.prize_position.x and y_pos == machine.prize_position.y:
    return (a_presses, b_presses, 3*a_presses + b_presses)
  else:
    #print (f"Ignoring: {machine} {X} {a_presses} {b_presses}  n {x_pos} {y_pos}")
    return None

def compute_result(machine: Machine) -> int:
  best = None

  max_presses_x = machine.prize_position.x / min(machine.button_a.x, machine.button_b.x)
  max_presses_y = machine.prize_position.y / min(machine.button_a.y, machine.button_b.y)
  max_presses = max(max_presses_x, max_presses_y)

  for a_presses in range(int(max_presses) + 2):
    for b_presses in range(int(max_presses) + 2):
      pos_x = a_presses * machine.button_a.x + b_presses * machine.button_b.x
      pos_y = a_presses * machine.button_a.y + b_presses * machine.button_b.y
      if pos_x == machine.prize_position.x and pos_y == machine.prize_position.y:
        cost = 3*a_presses + b_presses
        if best is None or cost < best[2]:
          best = (a_presses, b_presses, 3*a_presses + b_presses)
  return best


def part1(fname):
  cost = 0
  with open(fname) as f:
    while machine := consume_machine(f):
      ret = compute_result(machine)
      #print(ret)
      if ret:
        cost += ret[2]
  print(cost)

def part2(fname, increment=0):
  cost = 0
  with open(fname) as f:
    while machine := consume_machine(f, increment):
      ret = starting_over(machine)
      if ret:
        #print(ret)
        cost += ret[2]
  print(cost)

part1('day13/day13-input-easy.txt')
part2('day13/day13-input-easy.txt')
part2('day13/day13-input-easy.txt',increment=10_000_000_000_000)

#part1('day13/day13-input.txt')
part2('day13/day13-input.txt')
part2('day13/day13-input.txt', increment=10_000_000_000_000)