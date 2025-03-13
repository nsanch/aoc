import heapq
import sys

def dijkstra(graph, start, end):
  distances = {node: float('inf') for node in graph}
  distances[start] = 0
  priority_queue = [(0, start)]
  while priority_queue:
    current_distance, current_node = heapq.heappop(priority_queue)

    if current_distance > distances[current_node]:
      continue

    if current_node == end:
      break

    for neighbor, weight in graph[current_node].items():
      distance = current_distance + weight

      if distance < distances[neighbor]:
        distances[neighbor] = distance
        heapq.heappush(priority_queue, (distance, neighbor))

  return distances[end]

def make_graph(grid):
  start_node = None
  end_node = None
  graph = {}
  for y in range(len(grid)):
    for x in range(len(grid[y])):
      graph[(y, x)] = {}
      edges = graph[(y, x)]
      if grid[y][x] == "#":
        continue

      if grid[y][x] == "S":
        start_node = (y, x)
      if grid[y][x] == "E":
        end_node = (y, x)

      if y > 0 and grid[y-1][x] != "#":
        edges[(y-1, x)] = 1
      if y < len(grid) - 1 and grid[y+1][x] != "#":
        edges[(y+1, x)] = 1
      if x > 0 and grid[y][x-1] != "#":
        edges[(y, x-1)] = 1
      if x < len(grid[y]) - 1 and grid[y][x+1] != "#":
        edges[(y, x+1)] = 1
  return graph, start_node, end_node

def cheat(orig_grid, removed_wall):
  grid = [row.copy() for row in orig_grid]
  if grid[removed_wall[0]][removed_wall[1]] == "#":
    grid[removed_wall[0]][removed_wall[1]] = "X"
  graph, start_node, end_node = make_graph(grid)
  return dijkstra(graph, start_node, end_node), grid

def part1(fname, savings_threshold):
  grid = []
  with open(fname) as f:
    grid = [list(line.strip()) for line in f.readlines()]

  graph = {}
  start_node = None
  end_node = None
  graph, start_node, end_node = make_graph(grid)
  
  original_distance = dijkstra(graph, start_node, end_node)

  print(f"Original distance: {original_distance}")
  #for row in grid:
  #  print("".join(row))

  savings = {}
  for y in range(len(grid)):
    for x in range(len(grid[y])):
      if grid[y][x] == "#":
        new_distance, new_grid = cheat(grid, (y, x))
        win = original_distance - new_distance
        if win > savings_threshold:
          print(f"Removed wall at {y}, {x}. New distance: {new_distance}. Win: {win}")
          #for row in new_grid:
          #  print("".join(row))
          #print()
          #print()
          savings[win] = savings.get(win, 0) + 1
  
  total_big_wins = 0
  for win in sorted(savings.keys()):
    count = savings[win]
    #print(f"There are {count} ways to save {win} picoseconds.")
    total_big_wins += count
  return total_big_wins

print(part1("day20/day20-input-easy.txt", 0))
#print(part1("day20/day20-input.txt", 100)) 

def part2_try2(fname, savings_threshold):
  grid = []
  with open(fname) as f:
    grid = [list(line.strip()) for line in f.readlines()]

  graph = {}
  graph, start_node, end_node = make_graph(grid)
  original_distance = dijkstra(graph, start_node, end_node)

  print(f"Original distance: {original_distance}")

  distances_to_cheat_start = {}
  for start_y in range(len(grid)):
    for start_x in range(len(grid[start_y])):
      natural_distance = abs(start_y - start_node[0]) + abs(start_x - start_node[1])
      if natural_distance < original_distance:
        distance_from_start = dijkstra(graph, start_node, (start_y, start_x))
        if distance_from_start != float('inf'):
          distances_to_cheat_start[(start_y, start_x)] = distance_from_start

  distances_from_cheat_end = {}
  for end_y in range(len(grid)):
    for end_x in range(len(grid[end_y])):
      natural_distance = abs(end_y - end_node[0]) + abs(end_x - end_node[1])
      if natural_distance < original_distance:
        distance_to_end = dijkstra(graph, (end_y, end_x), end_node)
        if distance_to_end != float('inf'):
          distances_from_cheat_end[(end_y, end_x)] = distance_to_end

  print(f"Found {len(distances_to_cheat_start)} possible cheat starts and {len(distances_from_cheat_end)} possible cheat ends.")
  #print(f"{sorted(distances_to_cheat_start.keys())}")
  #print(f"{sorted(distances_from_cheat_end.keys())}")


  wins = [0 for x in range(original_distance+1)]
  attempts = 0

  for possible_cheat_start, distance_to_start in distances_to_cheat_start.items():
    for possible_cheat_end, distance_from_end in distances_from_cheat_end.items():
      natural_distance = abs(possible_cheat_start[0] - possible_cheat_end[0]) + abs(possible_cheat_start[1] - possible_cheat_end[1])
      if natural_distance <= 20 and (distance_from_end + distance_to_start + natural_distance) <= (original_distance - savings_threshold):
        #print(f"Attempting shortcut from {possible_cheat_start} to {possible_cheat_end}")
        graph[possible_cheat_start][possible_cheat_end] = natural_distance
        if attempts % 100 == 0:
          print(f"Attempted {attempts} shortcuts.")
        attempts += 1

        new_distance = (distance_from_end + distance_to_start + natural_distance)
        win = original_distance - new_distance
        if win >= savings_threshold:
          wins[win] += 1
          if sum(wins) % 100 == 0:
            print(f"{sum(wins)}: Shortcut from {possible_cheat_start} to {possible_cheat_end} saves {win} picoseconds for full path.")

  for i in range(len(wins)):
    if wins[i] > 0:
      print(f"{i}: {wins[i]}")
  return sum(wins)
        
#print(part2("day20/day20-input-easy.txt", 50))
print(part2_try2("day20/day20-input-easy.txt", 50))
print(part2_try2("day20/day20-input.txt", 100))