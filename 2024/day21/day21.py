import re
import unittest

numeric_keypad = {
  "7": { "8": ">", "9": ">>", "4": "v", "5": "v>", "6": "v>>", "1": "vv", "2": "vv>", "3": "vv>>", "0": ">vvv", "A": ">>vvv" },
  "8": { "7": "<", "9": ">", "4": "v<", "5": "v", "6": "v>", "1": "vv<", "2": "vv", "3": "vv>", "0": "vvv", "A": ">vvv" },
  "9": { "7": "<<", "8": "<", "4": "v<<", "5": "v<", "6": "v", "1": "vv<<", "2": "vv<", "3": "vv", "0": "vvv<", "A": "vvv" },
  "4": { "7": "^", "8": "^>", "9": "^>>", "5": ">", "6": ">>", "1": "v", "2": "v>", "3": "v>>", "0": ">vv", "A": ">>vv" },
  "5": { "7": "^<", "8": "^", "9": "^>", "4": "<", "6": ">", "1": "v<", "2": "v", "3": "v>", "0": "vv", "A": "vv>" },
  "6": { "7": "^<<", "8": "^<", "9": "^", "4": "<<", "5": "<", "1": "v<<", "2": "v<", "3": "v", "0": "vv<", "A": "vv" },
  "1": { "7": "^^", "8": "^^>", "9": "^^>>", "4": "^", "5": "^>", "6": "^>>", "2": ">", "3": ">>", "0": ">v", "A": ">>v" },
  "2": { "7": "^^<", "8": "^^", "9": ">^^", "4": "^<", "5": "^", "6": "^>", "1": "<", "3": ">", "0": "v", "A": "v>" },
  "3": { "7": "^^<<", "8": "^^<", "9": "^^", "4": "^<<", "5": "^<", "6": "^", "1": "<<", "2": "<", "0": "v<", "A": "v" },
  "0": { "7": "^^^<", "8": "^^^", "9": "^^^>", "4": "^^<", "5": "^^", "6": "^^>", "1": "^<", "2": "^", "3": "^>", "A": ">" },
  "A": { "7": "^^^<<", "8": "^^^<", "9": "^^^", "4": "^^<<", "5": "^^<", "6": "^^", "1": "^<<", "2": "^<", "3": "^", "0": "<" }
}

directional_keypad = {
  "^": { "A": ">", "<": "<v", "v": "v", ">": "v>" },
  "<": { "A": ">>^", "^": ">^", "v": ">", ">": ">>" },
  "v": { "A": ">^", "^": "^", "<": "<", ">": ">" },
  ">": { "A": "^", "^": "<^", "<": "<<", "v": "<" },
  "A": { "^": "<", "<": "v<<", "v": "<v", ">": "v" }
}

class Keypad(object):
  def __init__(self, keypad):
    self.keypad = keypad
    self.current = "A"

  def path_to(self, digit):
    if digit != self.current:
      presses = self.keypad[self.current][digit]
      return presses
    return ""

  def move_to(self, digit):
    path = self.path_to(digit)
    self.current = digit
    return path

  def reset(self):
    self.current = "A"

def directions_for_keypad(code, keypad):
  ret = []
  for digit in code:
    if digit != keypad.current:
      ret += keypad.move_to(digit)
    ret.append("A")
  return ret

keypad_1 = Keypad(numeric_keypad)
keypad_2 = Keypad(directional_keypad)
keypad_3 = Keypad(directional_keypad)

def full_directions(code):
  up_one = directions_for_keypad(list(code), keypad_1)
  up_two = directions_for_keypad(up_one, keypad_2)
  up_three = directions_for_keypad(up_two, keypad_3)
  return up_three

perms = {}
def permutations(s):
  if s in perms:
    return perms[s]
  if len(s) == 1:
    return [s]
  ret = set()
  for i in range(len(s)):
    sub_perms = permutations(s[:i] + s[i+1:])
    for perm in sub_perms:
      ret.add(s[i] + perm)
  ret = list(ret)
  perms[s] = ret
  return ret

def read_fully_for_one_digit(digit):
  paths_up_one = permutations(keypad_1.path_to(digit))
  print(f"paths to {digit} up one: {paths_up_one}")

read_fully_for_one_digit("0")
read_fully_for_one_digit("1")  

def part1(fname):
  keypad_1.reset()
  keypad_2.reset()
  keypad_3.reset()
  result = 0
  with open(fname, "r") as file:
    codes = [line.strip() for line in file.readlines()]
    for code in codes:
      directions = full_directions(code)
      cost = len(directions) * int(re.match(r"(\d+)", code).group(1))
      result += cost
      print(f"Code: {code}, Len: {len(directions)}, Cost: {cost}")
  return result

class TestDay21(unittest.TestCase):
  def test_029A(self):
    up_one = directions_for_keypad(list("029A"), keypad_1)
    up_two = directions_for_keypad(up_one, keypad_2)
    up_three = directions_for_keypad(up_two, keypad_3)
    self.assertEqual(up_one, list("<A^A>^^AvvvA"))
    self.assertEqual(up_two, list("v<<A>>^A<A>AvA<^AA>A<vAAA>^A"))
    self.assertEqual(up_three, list("<vA<AA>>^AvAA<^A>A<v<A>>^AvA^A<vA>^A<v<A>^A>AAvA^A<v<A>A>^AAAvA<^A>A"))

  def test_379A(self):
    up_one = directions_for_keypad(list("379A"), keypad_1)
    up_two = directions_for_keypad(up_one, keypad_2)
    up_three = directions_for_keypad(up_two, keypad_3)
    self.assertEqual(''.join(up_one), "^A^^<<A>>AvvvA")
    self.assertEqual(''.join(up_two), "<A>A<A>A")
    self.assertEqual(up_three, list(""))

#unittest.main()
print(part1("day21/day21-input-easy.txt"))
print(part1("day21/day21-input.txt"))