def is_match(grid, x1, y1, x2, y2, x3, y3, x4, y4):
  if grid[y1][x1] == 'X' and grid[y2][x2] == 'M' and grid[y3][x3] == 'A' and grid[y4][x4] == 'S':
    #print(x1, y1)

    return True
  return False

def search(grid):
  num_matches = 0
  for y in range(len(grid)):
    for x in range(len(grid[y])):
      if grid[y][x] != 'X':
        continue
      if y >= 3:
        # straight up
        if is_match(grid, x, y, x, y-1, x, y-2, x, y-3):
          num_matches += 1
        # up and to the left
        if x >= 3:
          if is_match(grid, x, y, x-1, y-1, x-2, y-2, x-3, y-3):
            num_matches += 1
        # up and to the right
        if x < len(grid[y]) - 3:
          if is_match(grid, x, y, x+1, y-1, x+2, y-2, x+3, y-3):
            num_matches += 1
      # to the right
      if x < len(grid[y]) - 3:
        if is_match(grid, x, y, x+1, y, x+2, y, x+3, y):
          num_matches += 1
      # to the left
      if x >= 3:
        if is_match(grid, x, y, x-1, y, x-2, y, x-3, y):
          num_matches += 1
      # down and to the left
      if y < len(grid) - 3 and x >= 3:
        if is_match(grid, x, y, x-1, y+1, x-2, y+2, x-3, y+3):
          num_matches += 1
      # down and to the right
      if y < len(grid) - 3 and x < len(grid[y]) - 3:
        if is_match(grid, x, y, x+1, y+1, x+2, y+2, x+3, y+3):
          num_matches += 1
      # straight down
      if y < len(grid) - 3:
        if is_match(grid, x, y, x, y+1, x, y+2, x, y+3):
          num_matches += 1

  return num_matches

def count_matches(filename):
  with open(filename, "r") as file:
    grid = [list(line.strip()) for line in file]
    return search(grid)

# part 2

def search_pt2(grid):
  num_matches = 0
  for y in range(1, len(grid) - 1):
    for x in range(1, len(grid[y]) - 1):
      if grid[y][x] != 'A':
        continue

      # normal MAS in the shape of an X
      check = ''.join([grid[y-1][x-1],
        grid[y][x],
        grid[y+1][x+1],
        grid[y-1][x+1],
        grid[y][x],
        grid[y+1][x-1]])
      #print(check)  
      options = set(['MASMAS', 'SAMSAM', 'MASSAM', 'SAMMAS'])
      if check in options:
        num_matches += 1
        
  return num_matches


def count_matches_pt2(filename):
  with open(filename, "r") as file:
    grid = [list(line.strip()) for line in file]
    return search_pt2(grid)


if __name__ == "__main__":
  print(count_matches("day4-input-easy2.txt"))
  print(count_matches("day4-input-easy.txt"))
  print(count_matches("day4-input.txt"))

  print(count_matches_pt2("day4-input-easy3.txt"))
  print(count_matches_pt2("day4-input.txt"))