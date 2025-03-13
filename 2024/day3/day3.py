import re

def part1(filename):
  data = ""
  with open(filename) as infile:
    data = infile.read()

  total_sum = 0

  for match in re.finditer(r"mul\((\d+),(\d+)\)", data):
    print(match.group(1), match.group(2))
    total_sum += int(match.group(1)) * int(match.group(2))

  return total_sum


mul_re = re.compile(r"mul\((\d+),(\d+)\)")
do_re = re.compile(r"do\(\)")
dont_re = re.compile(r"don't\(\)")

def consume_next(data, idx, total_sum, enabled):
  mul_m = mul_re.search(data[idx:])
  do_m = do_re.search(data[idx:])
  dont_m = dont_re.search(data[idx:])

  if not mul_m:
    return (len(data), total_sum, enabled)
  
  mul_idx = mul_m.start() if mul_m else len(data)
  do_idx = do_m.start() if do_m else len(data)
  dont_idx = dont_m.start() if dont_m else len(data)

  if (mul_idx < do_idx and mul_idx < dont_idx):
    if enabled:
      total_sum += int(mul_m.group(1)) * int(mul_m.group(2))
    return (idx + mul_m.end(), total_sum, enabled)
  elif (do_idx < mul_idx and do_idx < dont_idx):
    return (idx + do_m.end(), total_sum, True)
  elif (dont_idx < mul_idx and dont_idx < do_idx):
    return (idx + dont_m.end(), total_sum, False)


def part2(filename):
  data = ""
  with open(filename) as infile:
    data = infile.read()

  idx = 0
  total_sum = 0
  enabled = True
  while idx < len(data):
    idx, total_sum, enabled = consume_next(data, idx, total_sum, enabled)

  return total_sum



if __name__ == "__main__":
  print(part1("day3/day3-input-easy.txt"))
  print(part1("day3/day3-input.txt"))

  print(part2("day3/day3-input-easy2.txt"))
  print(part2("day3/day3-input.txt"))