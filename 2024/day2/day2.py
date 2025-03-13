
def is_safe(levels):
  #print(levels)
  last = levels[0]

  asc = False
  if levels[1] > levels[0]:
    asc = True
  
  for next in levels[1:]:
    if asc:
      if next <= last:
        #print("%s is not safe" % levels)
        return False
      if (next - last) > 3:
        #print("%s is not safe" % levels)
        return False
    else: # desc
      if next >= last:
        #print("%s is not safe" % levels)
        return False
      if (last - next) > 3:
        #print("%s is not safe" % levels)
        return False
    last = next

  #print("%s is safe" % levels)
  return True


def part1(fname):
  num_safe = 0
  with open(fname) as infile:
    for report in infile.readlines():
      levels = map(int, report.split())
      if is_safe(levels):
        num_safe += 1
  return num_safe


def part2(fname):
  num_safe = 0
  with open(fname) as infile:
    for report in infile.readlines():
      levels = map(int, report.split())
      if is_safe(levels):
        num_safe += 1
      else:
        dampened = False
        for i in range(len(levels)):
          new_list = levels[:i] + levels[i+1:]
          if is_safe(new_list):
            dampened = True
            break
        if dampened:
          num_safe += 1
  return num_safe

if __name__ == "__main__":
  print(part1("day2-input-easy.txt"))
  print(part1("day2-input.txt"))

  print(part2("day2-input-easy.txt"))
  print(part2("day2-input.txt"))