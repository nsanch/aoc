from enum import Enum

class WallTile(object):
  def __init__(self, position):
    self.position = position

  def move(self, the_map, direction):
    return False
  
  def can_move(self, the_map, direction):
    return False
  
  def __str__(self):
    return "#"

class EmptyTile(object):
  def __init__(self, position):
    self.position = position

  def move(self, the_map, direction):
    return True
  
  def can_move(self, the_map, direction):
    return True

  def __str__(self):
    return "."

class MovableTile(object):
  def __init__(self, positions):
    self.positions = positions # assume this is left then right

  def can_move(self, the_map, direction):
    next_positions = set([(p[0] + direction[0], p[1] + direction[1]) for p in self.positions]) - set(self.positions)
    for np in next_positions:
      next_tile = the_map.grid[np[0]][np[1]]
      if not next_tile.can_move(the_map, direction):
        return False
    return True

  def move(self, the_map, direction):
    if not self.can_move(the_map, direction):
      return False

    next_positions = set([(p[0] + direction[0], p[1] + direction[1]) for p in self.positions]) - set(self.positions)
    for np in next_positions:
      next_tile = the_map.grid[np[0]][np[1]]
      next_tile.move(the_map, direction)

    rng = range(len(self.positions))
    # if we're moving to the right we have to move the right-most position first.
    if direction == (0, 1):
      rng = range(len(self.positions) - 1, -1, -1)     
    for i in rng:
      p = self.positions[i]
      np = (p[0] + direction[0], p[1] + direction[1])
      the_map.update_tile_pos(self, p, np)
      self.positions[i] = np

    return True
  
class BoxTile(MovableTile):
  def __init__(self, positions):
    super().__init__(positions)
  
  def __str__(self):
    return "O"

class RobotTile(MovableTile):
  def __init__(self, position):
    super().__init__([position])
  
  def __str__(self):
    return "@"

class Map(object):
  def __init__(self, grid, double_wide):
    built_grid = []
    for i in range(len(grid)):
      row = []
      built_grid.append(row)
      for j in range(len(grid[i])):
        pos = (i, j)
        if double_wide:
          pos = (i, j*2)
        if grid[i][j] == '#':
          row.append(WallTile(pos))
          if double_wide:
            row.append(WallTile((pos[0], pos[1]+1)))
        elif grid[i][j] == '.':
          row.append(EmptyTile(pos))
          if double_wide:
            row.append(EmptyTile((pos[0], pos[1]+1)))
        elif grid[i][j] == 'O':
          if not double_wide:
            row.append(BoxTile([pos]))
          else:
            row.append(BoxTile([pos, (pos[0], pos[1]+1)]))
            row.append(row[-1])
        elif grid[i][j] == "@":
          row.append(RobotTile(pos))
          self.robot_location = row[-1]
          if double_wide:
            row.append(EmptyTile((pos[0], pos[1]+1)))
        else:
          print("Invalid character in grid: {grid[i][j]}")
  

    self.grid = built_grid
    self.double_wide = double_wide
    self.print_map()
  
  def update_tile_pos(self, tile, old_pos, new_pos):
    self.grid[old_pos[0]][old_pos[1]] = EmptyTile(old_pos)
    self.grid[new_pos[0]][new_pos[1]] = tile

  def accept_move(self, direction, print_it):
    self.robot_location.move(self, direction)
    if print_it:
      print(f"Move {direction}")
      self.print_map()
      print()

  def print_map(self):
    for row in self.grid:
      print(''.join([str(tile) for tile in row]))
  
  def box_positions(self):
    ret = []
    for row in self.grid:
      i = 0
      while i < len(row):
        tile = row[i]
        if isinstance(tile, BoxTile):
          ret.append(tile.positions[0])
          i += len(tile.positions)
        else:
          i += 1
    return ret

def direction_from_string(s):
  if s == '<':
    return (0, -1)
  elif s == '>':
    return (0, 1)
  elif s == 'v':
    return (1, 0)
  elif s == '^':
    return (-1, 0)
  else:
    print("Invalid direction string: {s}")

def part1(fname, print_moves):
  with open(fname, 'r') as f:
    grid = []
    while True:
      line = f.readline()
      if line.strip() == "":
        break
      grid.append(list(line.strip()))
    m = Map(grid, False)
    moves = ''.join([l.strip() for l in f.readlines()])
    #print(moves)

    for move in moves:
      m.accept_move(direction_from_string(move), print_moves)

    box_positions = m.box_positions()
    the_sum = sum([(100*pos[0]) + pos[1] for pos in box_positions])
    print(the_sum)
    return the_sum
  
def day15(fname, print_moves, double_wide):
  with open(fname, 'r') as f:
    grid = []
    while True:
      line = f.readline()
      if line.strip() == "":
        break
      grid.append(list(line.strip()))
    m = Map(grid, double_wide)
    moves = ''.join([l.strip() for l in f.readlines()])
    #print(moves)

    for move in moves:
      m.accept_move(direction_from_string(move), print_moves)

    box_positions = m.box_positions()
    the_sum = sum([(100*pos[0]) + pos[1] for pos in box_positions])
    print(the_sum)
    return the_sum

if __name__ == "__main__":
  day15("day15/day15-input-easy.txt", False, False)
  day15("day15/day15-input-easy2.txt", False, False)
  day15("day15/day15-input.txt", False, False)

  day15("day15/day15-input-easy2.txt", False, True)
  day15("day15/day15-input-easy3.txt", False, True)
  day15("day15/day15-input.txt", False, True)
