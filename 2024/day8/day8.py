class Map(object):
  def is_valid(self, pos, max_0, max_1):
    return (pos[0] >= 0 and pos[0] < max_0 and pos[1] >= 0 and pos[1] < max_1)

  def __init__(self, lines):
    self.nodes = {}
    self.antinodes = set()

    for y in range(len(lines)):
      for x in range(len(lines[y])):
        if lines[y][x] != ".":
          antenna_id = lines[y][x]
          if antenna_id not in self.nodes:
            self.nodes[antenna_id] = []
          
          self.nodes[antenna_id].append((y, x))

    for key in self.nodes.keys():
      positions = self.nodes[key]
      for i in range(len(positions)):
        for j in range(i+1, len(positions)):
          pos1 = positions[i]
          pos2 = positions[j]
          
          delta = (pos2[0] - pos1[0], pos2[1] - pos1[1])
          antinode1 = (pos1[0] - delta[0], pos1[1] - delta[1])

          antinode2 = (pos1[0] + delta[0], pos2[1] + delta[1])
          if self.is_valid(antinode1, len(lines), len(lines[0])):
            self.antinodes.add(antinode1)
          if self.is_valid(antinode2, len(lines), len(lines[0])):
            self.antinodes.add(antinode2)
  
  def num_antinodes(self):
    return len(self.antinodes)

def parse_map(fname):
  with open(fname) as f:
    lines = f.readlines()
    the_map = Map([l.strip() for l in lines])
    #print(the_map.antinodes)
    return the_map.num_antinodes()
  
# Part 2
class Map2(object):
  def is_valid(self, pos, max_0, max_1):
    return (pos[0] >= 0 and pos[0] < max_0 and pos[1] >= 0 and pos[1] < max_1)

  def __init__(self, lines):
    self.nodes = {}
    self.antinodes = set()

    for y in range(len(lines)):
      for x in range(len(lines[y])):
        if lines[y][x] != ".":
          antenna_id = lines[y][x]
          if antenna_id not in self.nodes:
            self.nodes[antenna_id] = []
          
          self.nodes[antenna_id].append((y, x))

    for key in self.nodes.keys():
      positions = self.nodes[key]
      for i in range(len(positions)):
        for j in range(i+1, len(positions)):
          pos1 = positions[i]
          pos2 = positions[j]
          
          delta = (pos2[0] - pos1[0], pos2[1] - pos1[1])
          while self.is_valid(pos1, len(lines), len(lines[0])):
            self.antinodes.add(pos1)
            pos1 = (pos1[0] - delta[0], pos1[1] - delta[1])

          while self.is_valid(pos2, len(lines), len(lines[0])):
            self.antinodes.add(pos2)
            pos2 = (pos2[0] + delta[0], pos2[1] + delta[1])
  
  def num_antinodes(self):
    return len(self.antinodes)

def parse_map2(fname):
  with open(fname) as f:
    lines = f.readlines()
    the_map = Map2([l.strip() for l in lines])
    #print(the_map.antinodes)
    return the_map.num_antinodes()
  
if __name__ == "__main__":
  print(parse_map("day8/day8-input-easy.txt"))
  print(parse_map("day8/day8-input.txt"))

  print(parse_map2("day8/day8-input-easy.txt"))
  print(parse_map2("day8/day8-input.txt"))
