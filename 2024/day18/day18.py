import heapq

def dijkstra(graph, start, end):
  distances = {node: float('inf') for node in graph}
  distances[start] = 0
  #best_paths = {}
  #for node in graph:
  #  best_paths[node] = []
  #best_paths[start].append([start])
  priority_queue = [(0, start)]  # (distance, node)

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

        #for path in best_paths[current_node]:
        #  best_paths[neighbor].append(path + [neighbor])
  
      #elif distance == distances[neighbor]:
      #  for path in best_paths[current_node]:
      #    best_paths[neighbor].append(path + [neighbor])

  return distances[end]


# the grid is N x N
# each point at each point in time is a node in the graph.
# when a byte falls, it means that the point disappears from that time onward.
# the cost of moving from one point to an adjacent one at the next tick is 1.
# read the grid. start at 0,0. end at N,N.
def part1(fname, n, num_bytes_to_corrupt):
  with open(fname) as f:
    falling_byte_positions = []
    for line in f.readlines():
      x,y = line.strip().split(',')
      falling_byte_positions.append((int(x), int(y)))

    fallen_bytes = set(falling_byte_positions[:num_bytes_to_corrupt])

    graph = {}
    for y in range(n+1):
      for x in range(n+1):
        graph[(x, y)] = {}
    for x in range(n+1):
      for y in range(n+1):
        if (x, y) not in graph:
          continue
        # we have to skip the bytes that have been corrupted.
        edges = graph[(x, y)]
        if (x > 0 and (x - 1, y) not in fallen_bytes):
          edges[(x - 1, y)] = 1
        if (x < n and (x + 1, y) not in fallen_bytes):
          edges[(x + 1, y)] = 1
        if (y > 0 and (x, y - 1) not in fallen_bytes):
          edges[(x, y - 1)] = 1
        if (y < n and (x, y + 1) not in fallen_bytes):
          edges[(x, y + 1)] = 1
    
    print("grid complete")
  
    result = dijkstra(graph, (0,0), (n,n))

    print(result)
    if result is None:
      print(f"no path found")
      return None
    distance = result #    distance, paths = result
    #if paths is None or len(paths) == 0:
    #  print(f"no path found")
    #  return None
    #else:
    print("Found a path from (0,0) to (n,n)")
    print(distance)
    #grid = [['.' for _ in range(n+1)] for _ in range(n+1)]
    #for x,y in paths[0]:
    #  grid[y][x] = 'O'
    #for b in fallen_bytes:
    #  x,y = b
    #  grid[y][x] = '#'
    #for row in grid:
    #  print(''.join(row))
    #print(paths[0])
    return distance
  
def do_with_given_fallen_bytes(fallen_bytes, n):
  graph = {}
  for y in range(n+1):
    for x in range(n+1):
      graph[(x, y)] = {}
  for x in range(n+1):
    for y in range(n+1):
      if (x, y) not in graph:
        continue
      # we have to skip the bytes that have been corrupted.
      edges = graph[(x, y)]
      if (x > 0 and (x - 1, y) not in fallen_bytes):
        edges[(x - 1, y)] = 1
      if (x < n and (x + 1, y) not in fallen_bytes):
        edges[(x + 1, y)] = 1
      if (y > 0 and (x, y - 1) not in fallen_bytes):
        edges[(x, y - 1)] = 1
      if (y < n and (x, y + 1) not in fallen_bytes):
        edges[(x, y + 1)] = 1
    
  distance = dijkstra(graph, (0,0), (n,n))
  return distance
      
def part2(fname, n):
  with open(fname) as f:
    falling_byte_positions = []
    for line in f.readlines():
      x,y = line.strip().split(',')
      falling_byte_positions.append((int(x), int(y)))

    for num_bytes_to_corrupt in range(len(falling_byte_positions)):
      fallen_bytes = set(falling_byte_positions[:num_bytes_to_corrupt + 1])
      distance = do_with_given_fallen_bytes(fallen_bytes, n)
      if distance != float('inf'):
        print(f"good as of time {num_bytes_to_corrupt}")
      else:
        print(f"bad as of time {num_bytes_to_corrupt}")
        print(falling_byte_positions[num_bytes_to_corrupt])
        break

part1("day18/day18-input-easy.txt", 6, 12)
part1("day18/day18-input.txt", 70, 1024)

part2("day18/day18-input-easy.txt", 6)
part2("day18/day18-input.txt", 70)