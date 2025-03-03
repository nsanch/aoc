import heapq
import re

numeric_graph = {
  "7": {"8": ">", "4": "v"},
  "8": {"7": "<", "9": ">", "5": "v"},
  "9": {"8": "<", "6": "v"},
  "4": {"7": "^", "5": ">", "1": "v"},
  "5": {"8": "^", "4": "<", "6": ">", "2": "v"},
  "6": {"9": "^", "5": "<", "3": "v"},
  "1": {"4": "^", "2": ">"},
  "2": {"1": "<", "5": "^", "3": ">", "0": "v"},
  "3": {"2": "<", "6": "^", "A": "v"},
  "0": {"2": "^", "A": ">"},
  "A": {"3": "^",  "0": "<"}
}

direction_graph = {
  "^": {"A": ">", "v": "v"},
  "<": {"v": ">"},
  "v": {">": ">", "<": "<", "^": "^"},
  ">": {"A": "^", "v": "<"},
  "A": {"^": "<", ">": "v"}
}
def modified_dijkstra(graph, start, ends):
  if start in ends:
    return 0, ["A"]

  distances = {node: float('inf') for node in graph}
  distances[start] = 0
  best_paths = {}
  for node in graph:
    best_paths[node] = []
  best_paths[start].append([])
  priority_queue = [(0, start)]  # (distance, node)

  while priority_queue:
    current_distance, current_node = heapq.heappop(priority_queue)

    if current_distance > distances[current_node]:
      continue

    #if current_node == end:
    #  break

    for neighbor, direction in graph[current_node].items():
      weight = len(direction)
      distance = current_distance + weight

      if distance < distances[neighbor]:
        distances[neighbor] = distance
        heapq.heappush(priority_queue, (distance, neighbor))

        for path in best_paths[current_node]:
          best_paths[neighbor].append(path + [direction])
  
      elif distance == distances[neighbor]:
        for path in best_paths[current_node]:
          best_paths[neighbor].append(path + [direction])

  best = None
  for e in ends:
    if best is None or distances[e] < best[0]:
      best = distances[e], [(''.join(p) + 'A') for p in best_paths[e]]
  return best

def all_keypad_paths(graph):
  all_paths = {}
  for start in graph:
    all_paths[(start, start)] = ["A"]
    for end in graph:
      if start != end:
        distance, paths = modified_dijkstra(graph, start, end)
        all_paths[(start, end)] = paths
  return all_paths

numeric_full_paths = all_keypad_paths(numeric_graph)
directional_full_paths = all_keypad_paths(direction_graph)

def num_repetitions(s):
  last = None
  repeats = 0
  for c in s:
    if c == last:
      repeats += 1
    last = c
  return repeats

class Keypad(object):
  def __init__(self):
    self.super_cache = {}

  def all_paths_one_digit(self, last, digit):
    pass

  def all_paths_to_full_code(self, digits):
    if digits in self.super_cache:
      return self.super_cache[digits]

    last = "A"
    best = None
    for d in digits:
      paths = self.all_paths_one_digit(last, d)
      if best is None:
        best = max(paths, key=lambda x: num_repetitions(x))
      else:
        #print(f"Combining paths. Orig: {len(ret)}, New: {len(paths)}")
        best = max([best + p for p in paths], key=lambda x: num_repetitions(x))
      last = d
      #assert min(map(len, ret)) == max(map(len, ret))

    self.super_cache[digits] = [best]
    return [best]
    #return [max(ret, key=lambda x: num_repetitions(x))]
    #print(f"There are {len(ret)} paths to full code {digits}")

    #return ret

class DirectKeypad(Keypad):
  def __init__(self, all_paths):
    super().__init__()
    self.all_paths = all_paths
  
  def all_paths_one_digit(self, last, digit):
    return self.all_paths[(last, digit)]


class IndirectKeypad(Keypad):
  def __init__(self, underlying, interface, depth):
    super().__init__()
    self.underlying = underlying
    self.interface = interface
    self.depth = depth
    self.cache = {}

  def all_paths_one_digit(self, last, digit):
    paths_in_underlying = self.underlying.all_paths_one_digit(last, digit)
    if (last, digit) in self.cache:
      return self.cache[(last, digit)]

    best = None
    best_distance = None
    for path in paths_in_underlying:
      indirect_paths = self.interface.all_paths_to_full_code(path)
      if best_distance is None:
        best_distance = len(indirect_paths[0])
        best = indirect_paths[0]
      if len(indirect_paths[0]) == best_distance:
        best = max([best] + indirect_paths, key=lambda x: num_repetitions(x))
      elif len(indirect_paths[0]) < best_distance:
        #print(f"One example that we're discarding. {best_paths[0]}: {len(ret[0])}, {ret[0]}")
        #print(f"New best distance                  {path}: {len(indirect_paths[0])}, {indirect_paths[0]}")
        best = max(indirect_paths, key=lambda x: num_repetitions(x))

    #print(f"Depth {self.depth}, last: {last}, digit: {digit}, best distance: {best_distance}, num_paths: {len(ret)}")

    self.cache[(last, digit)] = [best]
    return [best]
 
num_keypad = DirectKeypad(numeric_full_paths)
dir_keypad = DirectKeypad(directional_full_paths)
up_one = IndirectKeypad(num_keypad, dir_keypad, 1)
up_two = IndirectKeypad(up_one, dir_keypad, 2)

print(num_keypad.all_paths_to_full_code("029A")[0])
print(up_one.all_paths_to_full_code("029A")[0])
print(up_two.all_paths_to_full_code("029A")[0])

def part1(fname):
  result = 0
  with open(fname, "r") as file:
    codes = [line.strip() for line in file.readlines()]
    for code in codes:
      directions = up_two.all_paths_to_full_code(code)[0]
      cost = len(directions) * int(re.match(r"(\d+)", code).group(1))
      result += cost
      print(f"Code: {code}, Len: {len(directions)}, Cost: {cost}")
  return result

def part2(fname, depth):
  with open(fname, "r") as file:
    codes = [line.strip() for line in file.readlines()]
  codes = ["029A"]
  num_keypad = DirectKeypad(numeric_full_paths)
  dir_keypad = DirectKeypad(directional_full_paths)
  last_keypad = num_keypad
  result = 0
  for i in range(depth):
    last_keypad = IndirectKeypad(last_keypad, dir_keypad, i+1)
  
    for code in codes:
      all_paths = last_keypad.all_paths_to_full_code(code)
      directions = all_paths[0]
      cost = len(directions) * int(re.match(r"(\d+)", code).group(1))
      result += cost
      print(f"Code: {code}, Len: {len(directions)}, Cost: {cost}")
    print(result)
  
  return result

print(part1("day21/day21-input-easy.txt"))
print(part1("day21/day21-input.txt"))

print(part2("day21/day21-input.txt", depth=25))


