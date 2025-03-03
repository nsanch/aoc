class Map(object):
  def __init__(self, lines):
    self.lines = lines
    for i in range(len(self.lines)):
      for j in range(len(self.lines[i])):
        if self.lines[i][j] in ['^', 'v', '<', '>']:
          self.location = (i, j)
          dir = self.lines[i][j]
          if dir == '^':
            self.direction = (-1, 0)
          elif dir == 'v':
            self.direction = (1, 0)
          elif dir == '<':
            self.direction = (0, -1)    
          elif dir == '>': 
            self.direction = (0, 1)
  
    self.trail = set([self.location])
    self.trail_for_cycles = set([(self.location, self.direction)])

  def has_obstruction(self, i, j):
    return self.lines[i][j] == '#'

  def copy_wih_new_obstruction(self, i, j):
    new_lines = [l.copy() for l in self.lines]
    if new_lines[i][j] == '.':
      new_lines[i][j] = '#'
    return Map(new_lines)

  def __str__(self):
    return "\n".join(["".join(l) for l in self.lines])

  def turn_right(self):
    if self.direction == (-1, 0):
      self.direction = (0, 1)
    elif self.direction == (1, 0):
      self.direction = (0, -1)
    elif self.direction == (0, -1):
      self.direction = (-1, 0)
    elif self.direction == (0, 1):
      self.direction = (1, 0)

  def step(self):
    new_location = (self.location[0] + self.direction[0], self.location[1] + self.direction[1])
    if (new_location[0] >= len(self.lines) or new_location[1] >= len(self.lines[0]) or
      new_location[0] < 0 or new_location[1] < 0):
      return None
    if self.lines[new_location[0]][new_location[1]] == '#':
      self.turn_right()
      return self.step()
    else:
      self.trail_for_cycles.add((new_location, self.direction))
      self.trail.add(new_location)
      self.lines[new_location[0]][new_location[1]] = 'X'
      self.location = new_location
      return new_location
    
  def has_cycle(self):
    last_trail_size = len(self.trail)
    while self.step():
      if len(self.trail) == last_trail_size:
        return True
      last_trail_size = len(self.trail)
    return False

def part1(fname):
  with open(fname) as f:
    lines = list(map(lambda s: list(s.strip()), f.readlines()))
    the_map = Map(lines)
    while the_map.step():
      pass
    print(the_map)
    return len(the_map.trail)
  
def part2(fname):
  num_cycles = 0

  with open(fname) as f:
    lines = list(map(lambda s: list(s.strip()), f.readlines()))
    the_map = Map(lines)
    for i in range(len(lines)):
      for j in range(len(lines[0])):
        if the_map.has_obstruction(i, j):
          continue
        new_map = the_map.copy_wih_new_obstruction(i, j)
        if new_map.has_cycle():
          print(f"Cycle detected {num_cycles}")
          num_cycles += 1

    return num_cycles

if __name__ == "__main__":
  print(part1("day6/day6-input-easy.txt"))
  print(part1("day6/day6-input.txt"))

  #print(part2("day6/day6-input-easy.txt"))
  #print(part2("day6/day6-input.txt"))