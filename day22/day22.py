class FakeRNG(object):
  def __init__(self, seed, id):
    self.value = seed
    self.id = id
  
  def next(self):
    orig = self.value
    nextval = ((orig * 64) ^ orig) % 16777216
    nextval = ((nextval // 32) ^ nextval) % 16777216 
    nextval = ((nextval * 2048) ^ nextval) % 16777216
    self.value = nextval
    return self.value
  
  def curr(self):
    return self.value

def part1(fname):
  rngs = []
  id = 0
  with open(fname) as f:
    for line in f.readlines():
      seed = int(line.strip())
      rngs.append(FakeRNG(seed, id))
      id += 1
  for r in rngs:
    for i in range(2000):
      r.next()
    #print(r.curr())
  return sum([r.curr() for r in rngs])

print(f"Result: {part1('day22/day22-input-easy.txt')}")
print(f"Result: {part1('day22/day22-input.txt')}")

sequence_values = {}

def sequence(a,b,c,d):
  return (a+10) << 24 | (b+10) << 16 | (c+10) << 8 | (d+10)

def unpack_sequence(seq):
  return ((seq >> 24) - 10,
    ((seq >> 16) & 0xff) - 10,
    ((seq >> 8) & 0xff) - 10,
    (seq & 0xff) - 10
  )

def part2(fname):
  sequence_values = {}
  
  rngs = []
  id = 0
  with open(fname) as f:
    for line in f.readlines():
      seed = int(line.strip())
      rngs.append(FakeRNG(seed, id))
      id += 1

  for r in rngs:
    prices = []
    for i in range(2000):
      prices.append(r.next() % 10)

    for i in range(1, len(prices)-3):
      seq = sequence(prices[i] - prices[i-1],
                     prices[i+1] - prices[i],
                     prices[i+2] - prices[i+1],
                     prices[i+3] - prices[i+2])
      if seq not in sequence_values:
        sequence_values[seq] = {}
      if r.id not in sequence_values[seq]:
        sequence_values[seq][r.id] = prices[i+3]
    
  best = 0
  for sequence_id, id_to_price_map in sequence_values.items():
    total_price = sum(id_to_price_map.values())
    if total_price > best:
      best = total_price
      print(f"New best: {best}. Sequence: {sequence_id} {unpack_sequence(sequence_id)}")
      print(id_to_price_map)
  return best

#print(part2('day22/day22-input-easy2.txt'))
print(part2('day22/day22-input.txt'))