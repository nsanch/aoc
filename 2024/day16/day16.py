import heapq

def dijkstra(graph, start, ends):
  distances = {node: float('inf') for node in graph}
  distances[start] = 0
  best_paths = {}
  for node in graph:
    best_paths[node] = []
  best_paths[start].append([start])
  priority_queue = [(0, start)]  # (distance, node)

  while priority_queue:
    current_distance, current_node = heapq.heappop(priority_queue)

    if current_distance > distances[current_node]:
      continue

    #if current_node == end:
    #  break

    for neighbor, weight in graph[current_node].items():
      distance = current_distance + weight

      if distance < distances[neighbor]:
        distances[neighbor] = distance
        heapq.heappush(priority_queue, (distance, neighbor))

        for path in best_paths[current_node]:
          best_paths[neighbor].append(path + [neighbor])
  
      elif distance == distances[neighbor]:
        for path in best_paths[current_node]:
          best_paths[neighbor].append(path + [neighbor])

  best = None
  for e in ends:
    if best is None or distances[e] < best[0]:
      best = distances[e], best_paths[e]
  return best


def parse_input(fname):
  with open(fname) as f:
    grid = [l.strip() for l in f.readlines()]
    graph = {}

    start_pos = None
    end_pos = None
    orientations = [(0, 1), (0, -1), (1, 0), (-1, 0)]

    for i in range(len(grid)):
      for j in range(len(grid[i])):
        if grid[i][j] == 'S':
          start_node = (i, j, (0, 1))
        
        if grid[i][j] == 'E':
          end_node = (i, j, (0, 1))
      
        if grid[i][j] in ['S', 'E', '.']:
          for o in orientations:
            graph[(i, j, o)] = {}
          
          # it costs 1000 to turn clockwise and counterclockwise
          for o in orientations:
            clockwise = (o[1], -o[0])
            counter_clockwise = (-o[1], o[0])
            graph[(i, j, o)][(i, j, clockwise)] = 1000
            graph[(i, j, o)][(i, j, counter_clockwise)] = 1000

          # it costs 1 to move forward
          for o in orientations:
            next_pos = (i + o[0], j + o[1], o)
            if 0 <= next_pos[0] < len(grid) and 0 <= next_pos[1] < len(grid[i]) and grid[next_pos[0]][next_pos[1]] != '#':
              graph[(i, j, o)][next_pos] = 1


    return graph, start_node, [(end_node[0], end_node[1], o) for o in orientations]

def part1(fname):
  graph, start_node, end_nodes = parse_input(fname)
  orientations = [(0, 1), (0, -1), (1, 0), (-1, 0)]
  shortest_distance, best_paths = dijkstra(graph, start_node, end_nodes)
  
  nodes_along_best_paths = set()
  for path in best_paths:
    for pos in path:
      nodes_along_best_paths.add((pos[0], pos[1]))

  print(f"The shortest distance from {start_node} to {end_nodes} is: {shortest_distance} with {len(best_paths)} paths comprising {len(nodes_along_best_paths)} nodes")

part1("day16/day16-input-easy.txt")
part1("day16/day16-input.txt")