class Key(object):
  def __init__(self, schematic):
    self.pin_heights = [0] * 5
    for column in range(5):
      min_height = 0
      for row in range(len(schematic) - 1, -1, -1):
        if schematic[row][column] == '#':
          min_height = len(schematic) - row - 1
      self.pin_heights[column] = min_height
    
  def __str__(self):
    return ','.join(str(x) for x in self.pin_heights)

class Lock(object):
  def __init__(self, schematic):
    self.pin_heights = [0] * 5
    for column in range(len(schematic[0])):
      max_height = 0
      for row in range(len(schematic)):
        if schematic[row][column] == '#':
          max_height = row
      self.pin_heights[column] = max_height

  def accepts(self, key):
    for i in range(len(self.pin_heights)):
      if (key.pin_heights[i] + self.pin_heights[i]) > 5:
        return False
    return True
  
  def __str__(self):
    return ','.join(str(x) for x in self.pin_heights)

def part1(fname):
  keys = []
  locks = []
  with open(fname) as f:
    while True:
      schematic = []
      while True:
        line = f.readline().strip()
        if line == "":
          break
        schematic.append(line)
      if len(schematic) == 0:
        break
      if schematic[0][0] == '#':
        locks.append(Lock(schematic))
      else:
        keys.append(Key(schematic))

  ret = 0
  for lock in locks:
    #print(f"Lock: {lock}")
    for key in keys:
      #print(f"Key: {key}")
      if lock.accepts(key):
        print(f"Lock: {lock}, Key: {key} works")
        ret += 1
  return ret

print(part1("day25/day25-input-easy.txt"))
print(part1("day25/day25-input.txt"))
