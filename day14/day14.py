import re
import sys

class Robot(object):

  def __init__(self, position, direction, bounds):
    self.position = position
    self.direction = direction
    self.bounds = bounds
    self.t = 0

  def step(self):
    self.position[0] = (self.position[0] + self.direction[0]) % self.bounds[0]
    self.position[1] = (self.position[1] + self.direction[1]) % self.bounds[1]
    self.t += 1

  def step_n(self, n):  
    self.position[0] = (self.position[0] + (self.direction[0]*n)) % self.bounds[0]
    self.position[1] = (self.position[1] + (self.direction[1]*n)) % self.bounds[1]
    self.t += n

  def __str__(self):
    return f"Robot: position={self.position}, direction={self.direction}, bounds={self.bounds}, t={self.t}"
  
def look_in_string_for_sequence(s):
  longest_run = 0
  i = 0 
  in_run = False
  curr_run = 0
  while i < len(s):
    if s[i] != ".":
      if in_run:
        curr_run += 1
      else:
        in_run = True
        curr_run = 1
    else:
      if in_run:
        if curr_run > longest_run:
          longest_run = curr_run
        in_run = False
    i += 1

  return longest_run

def render_robots(robots, bounds, print_them):
  grid = [['.' for x in range(bounds[0])] for y in range(bounds[1])]
  for r in robots:
    curr = grid[r.position[1]][r.position[0]]
    if curr == ".":
      grid[r.position[1]][r.position[0]] = str(1)
    else:
      grid[r.position[1]][r.position[0]] = str(int(curr) + 1)


  sum_of_runs = 0
  for row in grid:
    sum_of_runs += look_in_string_for_sequence("".join(row))
    if print_them:
      print("".join(row))

  return sum_of_runs

def consume_robots(fname, bounds):
  robot_re = re.compile(r"p=(\d+),(\d+) v=(\-?\d+),(\-?\d+)")

  with open(fname) as f:
    lines = f.readlines()

    robots = []
    for line in lines:
      m = robot_re.match(line)
      if m:
        position = [int(m.group(1)), int(m.group(2))]
        direction = [int(m.group(3)), int(m.group(4))]
        robot = Robot(position, direction, bounds)
        #print(str(robot))
        robots.append(robot)

    return robots

def part1(fname, bounds):
  robots = consume_robots(fname, bounds)

  quad1 = 0
  quad2 = 0
  quad3 = 0
  quad4 = 0
  for r in robots:
    r.step_n(100)
    #print(r.position)
    if r.position[0] < (bounds[0]-1)/2 and r.position[1] < (bounds[1]-1)/2:
      quad1 += 1
    elif r.position[0] > (bounds[0]-1)/2 and r.position[1] < (bounds[1]-1)/2:
      quad2 += 1
    elif r.position[0] < (bounds[0]-1)/2 and r.position[1] > (bounds[1]-1)/2:
      quad3 += 1
    elif r.position[0] > (bounds[0]-1)/2 and r.position[1] > (bounds[1]-1)/2:
      quad4 += 1

  render_robots(robots, bounds, False)
  print("final quad counts: ", quad1, quad2, quad3, quad4)
  return quad1 * quad2 * quad3 * quad4
    
def part2(fname, bounds):
  robots = consume_robots(fname, bounds)

  t = 0
  prev_max = 0

  for r in robots: r.step_n(6875)

  while True:
    t += 1
    for r in robots:
      r.step()

    #print("===================")
    #print("t=", t)
    #print("===================")
    sum_of_runs = render_robots(robots, bounds, False)
    print(f"t={t}, sum of runs: {sum_of_runs}") 
    if sum_of_runs > prev_max:
      print(f"broke our record! {prev_max} -> {sum_of_runs}")
      prev_max = sum_of_runs
      print("Continue? (y/n)")
      if sys.stdin.readline().strip() == "n":
        render_robots(robots, bounds, True)
        break

print(part1("day14/day14-input-easy.txt", [11, 7]))
print(part1("day14/day14-input.txt", [101, 103]))

part2("day14/day14-input.txt", [101, 103])