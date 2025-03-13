import math

def compute_distances_part1():
  col1 = []
  col2 = []
  with open("day1-input.txt") as infile:
    data = infile.read().splitlines()
    for line in data:
      (x, y) = line.split()
      col1.append(int(x))
      col2.append(int(y))

  col1.sort()
  col2.sort()
  #print(col1)
  #print(col2)

  distances = [abs(t1 - t2) for t1, t2 in zip(col1, col2)]
  #print(distances)
  cumsum = sum(distances)
  return cumsum

def compute_simscore_part2(fname):
  col1 = []
  col2 = []
  with open(fname) as infile:
    data = infile.read().splitlines()
    for line in data:
      (x, y) = line.split()
      col1.append(int(x))
      col2.append(int(y))

  frequencies = {}
  for y in col2:
    if y in frequencies:
      frequencies[y] += 1
    else:
      frequencies[y] = 1
  
  simscore = 0
  for x in col1:
    if x in frequencies:
      simscore += x * frequencies[x]
  return simscore


if __name__ == "__main__":
  print(compute_distances_part1())
  print(compute_simscore_part2("day1-input.txt"))
  