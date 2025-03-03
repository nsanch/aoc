import math

class Position(object):
  def __init__(self, y, x):
    self.y = y
    self.x = x
  
  def __eq__(self, value):
    return value and (self.y == value.y and self.x == value.x)

  def __str__(self):
    return f"Position ({self.y}, {self.x})"
  
  def __hash__(self):
    return hash((self.y, self.x))
  
  def neighbors(self, max_y, max_x):
    ret = []
    # left
    if self.x > 0:
      ret.append(Position(self.y, self.x - 1))
    # up
    if self.y > 0:
      ret.append(Position(self.y - 1, self.x))
    # right
    if self.x + 1 < max_x:
      ret.append(Position(self.y, self.x + 1))
    # down
    if self.y + 1 < max_y:
      ret.append(Position(self.y + 1, self.x))
    return ret

class Node(object):
  def __init__(self, data):
    self.visited = False
    self.data = data

  def reset_visited(self):
    self.visited = False

  def visit(self):
    self.visited = True

  def __hash__(self):
    return hash(self.data)
  
  def __eq__(self, other):
    return isinstance(other, Node) and self.data == other.data

class Graph(object):
  def __init__(self):
    self.nodes: dict[Position, Node] = {}
    self.edges: dict[Node, list[Node]] = {}

  def get_node(self, data: Position):
    if data in self.nodes:
      return self.nodes[data]

    n = Node(data)
    self.nodes[data] = n
    return n

  def add_node(self, data):
    if data not in self.nodes:
      self.nodes[data] = Node(data)
  
  def add_edge(self, data1, data2):
    n1 = self.get_node(data1)
    n2 = self.get_node(data2)
    if n1 not in self.edges:
      self.edges[n1] = set()
    self.edges[n1].add(n2)
    if n2 not in self.edges:
      self.edges[n2] = set()
    self.edges[n2].add(n1)
  
  def num_edges(self, data: Position):
    if data not in self.nodes:
      return 0
    if self.get_node(data) not in self.edges:
      return 0
    return len(self.edges[self.get_node(data)])
  
  def neighbors(self, data: Position):
    return [e.data for e in self.edges[self.get_node(data)]]
  
  def count_sides(self):
    if len(self.nodes) == 1:
      return 4

    boundary_points = [n for n in self.nodes if self.num_edges(n.data) < 4]

    corners = []
    pokey_ends = []
    for n in self.nodes:
      edges = self.neighbors(n)
      if len(edges) == 2:
        dir1 = (edges[0].y - n.y, edges[0].x - n.x)
        dir2 = (edges[1].y - n.y, edges[1].x - n.x)
        if dir1 != (-1*dir2[0], -1*dir2[1]):
          corners.append(n)
        else:
          # if there's a one-unit side, then this is actually a corner too.
          pass
      if len(edges) == 1:
        pokey_ends.append(n)
    
    print(f"corners: {' '.join(map(str, corners))}")
    print(f"pokey ends: {' '.join(map(str, pokey_ends))}")
    
    return len(corners) + 3*len(pokey_ends)
  
  def get_boundary_points(self):
    return [data for data in self.nodes if self.num_edges(data) < 4]
  
  # take all boundary points. enumerate their neighbors outside the region, including the
  # ones that are off the grid. take that boundary point and see if its in a line with any
  # of the existing perimeters. if so, add it to that line. if not, create a new line.
  # we need to enumerate the boundary points in some order. one way to do that is to put all
  # the boundary points in a graph whose edges are their neighbors on the boundary. 
  def walk_graph_to_get_boundary(self):
    boundary_points = [n for n in self.nodes if self.num_edges(n) < 4]
    curr = boundary_points[0]
    prev_lines = []
    curr_line = Line([curr])
    self.walk_graph_to_get_boundary_helper(boundary_points, None, curr, curr_line, None, prev_lines, set(boundary_points))
    return prev_lines
  
  def walk_graph_to_get_boundary_helper(
      self, boundary_points, prev, curr, curr_line, prev_direction, prev_lines, unvisited_boundary_points):
    print(f"walking at {curr}")

    # examine the next options to find the best next one to visit. there will be up to 3, but we don't want
    # to go backwards -- recursion covers the case where we need to back up.
    # forwards in same direction is preferred.
    # turning in different direction is next best.
    # but we visit anythign we haven't visited before.
    ns = self.neighbors(curr)
    # there can be zero or one forwards options
    forwards = list(filter(lambda n: n != prev and ((n.y - curr.y, n.x - curr.x) == prev_direction), ns))
    # there can be zero, one or two turning options.
    turning = list(filter(lambda n: n != prev and ((n.y - curr.y, n.x - curr.x) != prev_direction), ns))
    preference_order = forwards + turning

    for p in preference_order:
      if p in unvisited_boundary_points:
        curr_line.add(p)
        next_direction = (p.y - curr.y, p.x - curr.x)
        if prev_direction and next_direction != prev_direction:
          # we changed direction. start a new line and add this line to prev_lines.
          prev_lines.append(curr_line)
          next_line = Line([curr])
        else:
          next_line = curr_line
        # recurse
        unvisited_boundary_points.remove(p)
        self.walk_graph_to_get_boundary_helper(boundary_points=boundary_points,
                                                prev = curr,
                                                curr = p,
                                                curr_line = next_line.copy(),
                                                prev_direction = next_direction,
                                                prev_lines = prev_lines,
                                                unvisited_boundary_points = unvisited_boundary_points)
    


class Line(object):
  def __init__(self, points):
    self.points = points
  
  def copy(self):
    return Line(self.points.copy())

  def can_add(self, point):
    my_slope = self.slope()

    if my_slope:
      if Line.compute_slope(point, self.first_point()) == my_slope:
        return 0
      elif Line.compute_slope(self.last_point(), point) == my_slope:
        return len(self.points)
      else:
        return None
    else:
      self_point = self.first_point()
      #print(f"trying to see if {self.first_point()} is next to {point}")
      if Line.compute_slope(point, self_point) in [(0, 1), (1, 0)]:
        #print("insert at beginning")
        return 0
      elif Line.compute_slope(self_point, point) in [(1, 0), (0, 1)]:
        #print("insert at end")
        return 1
      else:
        #print("nope")
        return None

  def add(self, point):
    idx = self.can_add(point)
    if idx is not None:
      self.points.insert(idx, point)
    else:
      print(f"ERROR cant add new point {point} to line {self}")
    
  
  def reverse_copy(self):
    return Line(self.points.copy().reverse())

  @classmethod
  def compute_slope(cls, point_a: Position, point_b: Position):
    return (point_b.y - point_a.y, point_b.x - point_a.x)


  def first_point(self): return self.points[0]
  def last_point(self): return self.points[-1]

  def slope(self):
    if len(self.points) > 1:
      return Line.compute_slope(self.points[0], self.points[1])
    else:
      return None

  def absorb(self, other_line):
    my_slope = self.slope()
    other_slope = other_line.slope()

    if my_slope:
      if other_slope:
        if my_slope != other_slope:
          neg_other_slope = (-1*other_slope[0], -1*other_slope[1])
          if (my_slope == neg_other_slope):
            other_line = other_line.reverse_copy()
            other_slope = neg_other_slope
          else:
            return False

      if Line.compute_slope(other_line.last_point(), self.first_point()) == my_slope:
        self.points = other_line.points + self.points
      elif Line.compute_slope(self.last_point(), other_line.first_point()) == my_slope:
        self.points = self.points + other_line.points
      else:
        return False
      
      return True
    else: # there's only one point on this line
      if other_slope:
        other_acceptance = other_line.can_add(self.first_point())
        if other_acceptance is not None:
          self.points = other_line.points.copy()
          self.points.insert(other_acceptance, self.first_point())
          return True
        else:
          return False
      else: # other line is one point
        pos = self.can_add(other_line.first_point())
        if pos is not None:
          self.points.insert(pos, other_line.first_point())
          return True
        else:
          return False

  def __str__(self):
    return f"Line: {str(self.first_point())} -> {str(self.last_point())}"

class Region(object):
  def __init__(self, plant, area, perimeter, positions):
    self.plant = plant
    self.area = area
    self.perimeter = perimeter
    self.positions = positions

  def price(self):
    return self.area * self.perimeter
  
  def __str__(self):
    return f"Region {self.plant}: area {self.area}, perimeter {self.perimeter}"

class Map(object):
  def at(self, pos):
    return self.grid[pos.y][pos.x]

  def __init__(self, lines):
    self.grid = lines

    partitioned_already = set()
    self.regions: list[Region] = []
    for y in range(len(self.grid)):
      for x in range(len(self.grid[y])):
        pos = Position(y, x)
        if pos in partitioned_already:
          continue

        plant = self.at(pos)

        partition = set()
        self.add_to_full_region(plant, pos, partition)
        partitioned_already = partitioned_already.union(partition)
        self.regions.append(self.partition_to_region(plant, partition))

  def add_to_full_region(self, plant, seed_member, region_acc):
    neighbors = seed_member.neighbors(len(self.grid), len(self.grid[0]))
    region_acc.add(seed_member)
    for n in neighbors:
      if n not in region_acc:
        if self.at(n) == plant:
          self.add_to_full_region(plant, n, region_acc)

  def boundaries_of_point(self, plant: str, pos: Position) -> list[list[Position]]:
    up_sides = []
    down_sides = []
    left_sides = []
    right_sides = []
    left = Position(pos.y, pos.x - 1)
    right = Position(pos.y, pos.x + 1)
    up = Position(pos.y - 1, pos.x)
    down = Position(pos.y + 1, pos.x)
    if left.x < 0 or self.at(left) != plant:
      left_sides.append(Position(pos.y, 2*pos.x - 1))
    if right.x >= len(self.grid[pos.y]) or self.at(right) != plant:
      right_sides.append(Position(pos.y, 2*pos.x + 1))
    if up.y < 0 or self.at(up) != plant:
      up_sides.append(Position(2*pos.y - 1, pos.x))
    if down.y >= len(self.grid) or self.at(down) != plant:
      down_sides.append(Position(2*pos.y + 1, pos.x))

    #print(f"evens: {' '.join(map(str, evens))}")
    #print(f"odds: {' '.join(map(str, odds))}")
    return [left_sides, right_sides, up_sides, down_sides]
#    return horizontals, verticals
  
  def num_boundaries_of_point(self, plant, pos):
    return sum(map(len, self.boundaries_of_point(plant, pos)))
    #return len(xs) + len(ys)

  def partition_to_region(self, plant, partition): 
    region_area = len(partition)
    region_perimeter = 0
    for pos in partition:
      region_perimeter += self.num_boundaries_of_point(plant, pos)
      
    r = Region(plant, region_area, region_perimeter, partition)
    return r

  def region_as_graph(self, region):
    g = Graph()
    for pos in region.positions:
      g.add_node(pos)
      for neighbor in pos.neighbors(len(self.grid), len(self.grid[pos.x])):
        if neighbor in region.positions:
          g.add_edge(pos, neighbor)
    return g
  
def walk_to_corner(grid, starting_point, curr, direction, visited_points):
  visited_points.add((curr, direction))
  if (curr.x == len(grid[curr.y]) - 1):
    return curr
  else:
    next_one = Position(curr.y + direction.x, curr.x + 1) + direction.x
    if grid.at(next_one) != grid.at(starting_point):
      return curr
    else:
      return walk_to_corner(grid, starting_point, next_one)

def walk_from_00(grid, starting_point):
  visited_points = set()
  x_boundary = walk_to_corner(grid, starting_point, starting_point, Position(0, 1), visited_points)
  y_boundary = walk_to_corner(grid, starting_point, starting_point, Position(1, 0), visited_points)
  return

def clump(boundary_points: list[Position], the_map: Map) -> list[Line]:
  partitions_of_side_points: list[set[Position]] = [set(), set(), set(), set()]
  for bp in boundary_points:
    to_add_in = the_map.boundaries_of_point(the_map.at(bp), bp)
    for i in range(0, 4):
      for x in to_add_in[i]:
        partitions_of_side_points[i].add(x)

  clumps = []
  for p in partitions_of_side_points:
    clumps += clump_helper(list(p), the_map)
  return clumps

def clump_helper(boundary_points: list[Position], the_map: Map) -> list[Line]:
  clumps: list[Line] = []
  for p in boundary_points:
    clumps.append(Line([p]))
  should_continue = True
  while should_continue:
    should_continue = False
    #print(f"there are {len(clumps)} clumps")
    i = 0
    while i < len(clumps):
      j = i+1
      while j < len(clumps):
        could_absorb = clumps[i].absorb(clumps[j])
        if could_absorb:
          #print("absorbed")
          del clumps[j]
          should_continue = True
        else:
          #print(f"did not absorb {clumps[i]} and {clumps[j]}")
          j += 1
      i += 1

  #print(f"done clumping. there are {len(clumps)} clumps")
  return clumps



def part1(fname):
  with open(fname) as f:
    lines = f.readlines()
    the_map = Map([l.strip() for l in lines])
    regions = the_map.regions
    #print([(r.plant, r.price()) for r in regions])
    return sum([r.price() for r in regions])
  
def part2(fname):
  with open(fname) as f:
    lines = f.readlines()
    the_map = Map([l.strip() for l in lines])
    regions = the_map.regions
    price = 0
    for r in regions:
      g = the_map.region_as_graph(r)
      bp = g.get_boundary_points()
      clumps = clump(bp, the_map)
      #print(f"Region {r.plant} has {len(horizontals)} horizontals and {len(verticals)} verticals")
      #for c in horizontals:
      #  print(f"Region {r.plant} has horizontal side: {c}")
      #for c in verticals:
      #  print(f"Region {r.plant} has vertical side: {c}")
      price += r.area * len(clumps)#(len(horizontals) + len(verticals))
      #print(f"Region {r.plant} has {r.area} area and {r.perimeter} perimeter with {len(clumps)} sides")
    return price


if __name__ == "__main__":
  print(part1("day12/day12-input-easy.txt"))
  print(part1("day12/day12-input.txt")) 

  print("PART 2 ================") 

  print(part2("day12/day12-input-easy.txt"))
  print(part2("day12/day12-input.txt"))


# a b c
# d X f
# g h i

# X X X
# X X X
# g h i

# ignore the interior points entirely?