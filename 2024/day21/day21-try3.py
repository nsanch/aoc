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

distance_cache = {}
def distance_at_depth(depth, desired_depth, start, end): 
  if (depth, start, end) in distance_cache:
    return distance_cache[(depth, start, end)]

  paths = []
  if depth == 0:
    paths = numeric_full_paths[(start, end)]
  else:
    paths = directional_full_paths[(start, end)]

  if depth == desired_depth:
    distance_cache[(depth, start, end)] = len(min(paths, key=len))
    return len(paths[0])
  
  best_distance = float('inf')
  for path in paths:
    last = "A"
    path_distance = distance_for_full_code(path, depth + 1, desired_depth)
#    for curr in path:
#      path_distance += distance_at_depth(depth + 1, desired_depth, last, curr)
    if path_distance < best_distance:
      best_distance = path_distance
  
  distance_cache[(depth, start, end)] = best_distance
  return best_distance

def distance_for_full_code(path, depth, desired_depth):
  path_distance = 0
  last = "A"
  for curr in path:
    path_distance += distance_at_depth(depth, desired_depth, last, curr)
    last = curr
  return path_distance

def part1(fname):
  result = 0
  with open(fname, "r") as file:
    codes = [line.strip() for line in file.readlines()]
    for code in codes:
      length = distance_for_full_code(code, 0, 25)
      cost = length * int(re.match(r"(\d+)", code).group(1))
      result += cost
      print(f"Code: {code}, Len: {length}, Cost: {cost}")
  return result

print(part1("day21/day21-input-easy.txt"))
print(part1("day21/day21-input.txt"))