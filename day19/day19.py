def is_possible(available_towels, desired_towel):
  for avail in available_towels:
    if desired_towel == avail:
      return True
    if desired_towel.startswith(avail):
      if is_possible(available_towels, desired_towel[len(avail):]):
        return True
  return False

ways_to_desired_towel = {}
def all_ways(available_towels, desired_towel):
  if desired_towel in ways_to_desired_towel:
    return ways_to_desired_towel[desired_towel]
  results = 0
  for avail in available_towels:
    if desired_towel == avail:
      results += 1
    elif desired_towel.startswith(avail):
      results += all_ways(available_towels, desired_towel[len(avail):])
  ways_to_desired_towel[desired_towel] = results
  return results

def part1(fname):
  with open(fname) as f:
    available_towels = [t.strip() for t in f.readline().split(',')]
    blank = f.readline()
    desired = [d.strip() for d in f.readlines()]

    #print(desired)
    #print(available_towels)

    count = 0
    for d in desired:
      if is_possible(available_towels, d):
        count += 1
        #print(f"{d} is possible")
      else:
        pass
        #print(f"{d} is not possible")

    return count
  
def part2(fname):
  ways_to_desired_towel.clear()
  with open(fname) as f:
    available_towels = [t.strip() for t in f.readline().split(',')]
    blank = f.readline()
    desired = [d.strip() for d in f.readlines()]

    count = 0
    for d in desired:
      options = all_ways(available_towels, d)
      count += options
      #print(f"{d} has {options} ways")
    return count

print(part1('day19/day19-input-easy.txt'))
print(part1('day19/day19-input.txt'))

print(part2('day19/day19-input-easy.txt'))
print(part2('day19/day19-input.txt'))
