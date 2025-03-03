class Map(object):
  def __init__(self, lines):
    self.grid = []
    for line in lines:
      self.grid.append([int(c) for c in line.strip()])
  
  def elevation_at(self, y, x):
    if y >= 0 and y < len(self.grid):
      if x >= 0 and x < len(self.grid[y]):
        return self.grid[y][x]
    return None
  
  def find_all_trails(self, y, x, previous_trail):
    # if curr_elevation = 9, return [current]
    # else return [current, next..] for every next.. that the recursion gives
    # returns an array of trails, where each trail is an array of points

    curr_elevation = self.elevation_at(y, x)
    if curr_elevation == 9:
      return [[(y, x)]]

    ret = []
    for possible_next_step in [(y-1, x), (y + 1, x), (y, x-1), (y, x+1)]:
      if possible_next_step in previous_trail:
        next
        
      next_el = self.elevation_at(possible_next_step[0], possible_next_step[1])
      if next_el and (next_el == curr_elevation + 1):
        remaining_trails = self.find_all_trails(possible_next_step[0], possible_next_step[1], previous_trail + [(y, x)])
        for rt in remaining_trails:
          rt.insert(0, (y, x))
          ret.append(rt)
    return ret
  
  def find_trailheads(self):
    ret = []
    for y in range(len(self.grid)):
      for x in range(len(self.grid[y])):
        if self.elevation_at(y, x) == 0:
          ret.append((y, x))
    #print(f"trailheads are {ret}")
    return ret
  
  def count_all_full_trailheads(self):
    scores = 0
    ratings = 0
    for th in self.find_trailheads():
      ts = self.find_all_trails(th[0], th[1], [])
      final_destinations = set([t[-1] for t in ts])
      #print(f"trailhead {th} has {len(ts)} leading to {len(final_destinations)} unique destinations")
      scores += len(final_destinations)
      ratings += len(ts)
    return (scores, ratings)
  
def part1(fname):
  with open(fname) as f:
    lines = f.readlines()
    the_map = Map(lines)
    return the_map.count_all_full_trailheads()
  
if __name__ == "__main__":
  print(part1("day10/day10-input-easy.txt"))
  print(part1("day10/day10-input.txt"))