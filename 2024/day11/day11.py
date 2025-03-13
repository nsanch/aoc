import math

def iterate_item(item):
  next_version = []
  if item == 0:
    next_version.append(1)
  else:
    num_digits = math.floor(math.log(item, 10)) + 1
    if num_digits % 2 == 0:
      splitter = 10 ** (num_digits / 2)
      left_half = int(item / splitter)
      right_half = int(item % splitter)
      next_version.append(left_half)
      next_version.append(right_half)
    else:
      next_version.append(item * 2024)
  return next_version

def iterate(lst):
  next_version = []
  for item in lst:
    next_version += iterate_item(item)
  return next_version

three_deep = {
  '0': [4048, 1, 4048, 8096],
}

def iterate5(lst):
  next_version = []
  for item in lst:
    if item in three_deep:
      next_version += three_deep[item]
    else:
      computed = iterate(iterate(iterate(iterate(iterate_item(item)))))
      #print(f"computed {item}")
      three_deep[item] = computed
      next_version += computed
  return next_version


def part1(fname, depth):
  with open(fname) as f:
    lst = [int(x) for x in f.readlines()[0].strip().split()]
    print(lst)
    for i in range(depth):
      lst = iterate(lst)
      #print(lst)
      print(f"{i}: {len(lst)}")

precomputed = {}

def iterate_fully(item, depth):
  if item < 1000:
    if item not in precomputed:
      precomputed[item] = {}
    if depth in precomputed[item]:
      return precomputed[item][depth]

  total = 0
  lst = [([item], depth)]
  while len(lst) > 0:
    (iteration, level) = lst.pop()
    if level == 0:
      total += 1
      if total % 1000000 == 0:
        num_at_orig_depth = len([x for x in lst if x[1] == depth])
        print(f"passing {total} items. {len(lst)} in list. {num_at_orig_depth} still left alone.")
    elif level > 5:
      next_iter = iterate5(iteration)
      for item in next_iter:
        lst.append(([item], level - 5))
    else:
      next_iter = iterate(iteration)
      for item in next_iter:
        lst.append(([item], level - 1))
  
  #print(f"computed {depth} for {item}: {total}")
  if item in precomputed:
    precomputed[item][depth] = total
  return total

def iterate_fully_recursive(item, depth):
  if depth == 0:
    return 1
  
  if item < 1000:
    if item not in precomputed:
      precomputed[item] = {}
    if depth in precomputed[item]:
      return precomputed[item][depth]
  
  total = 0
  next_iteration = iterate_item(item)
  for next_item in next_iteration:
    total += iterate_fully_recursive(next_item, depth - 1)
  
  if item in precomputed:
    precomputed[item][depth] = total

  return total

def part2(fname, depth):
  total = 0
  with open(fname) as f:
    lst = []
    for x in f.readlines()[0].strip().split():
      total += iterate_fully_recursive(int(x), depth)
  print(f"{depth}: {total}")
  return total

print(part2("day11/day11-input.txt", 75))