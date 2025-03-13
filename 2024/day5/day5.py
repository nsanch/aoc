import re

class Rule(object):
  def __init__(self, a, b):
    self.a = a
    self.b = b

  def check(self, update):
    if self.a in update and self.b in update:
      ret = update.index(self.a) < update.index(self.b)
    else:
      ret = True
    
    return ret
  
  def __str__(self):
    return f"Rule({self.a}, {self.b})"

def check_rules(rules, update):
  for rule in rules:
    if not rule.check(update):
      return False
  return True

def partition_list(list, predicate):
  true_list = []
  false_list = []
  for i in list:
    if predicate(i):
      true_list.append(i)
    else:
      false_list.append(i)
  return true_list, false_list

class RulesGraph(object):
  def __init__(self, rules):
    self.nodes = {}
    self.ants = {}
    for r in rules:
      if r.a not in self.nodes:
        self.nodes[r.a] = []
      self.nodes[r.a].append(r.b)

      if r.b not in self.ants:
        self.ants[r.b] = []
      self.ants[r.b].append(r.a)

  def antecedents(self, n):
    return self.ants.get(n, [])
  
  def followers(self, n):
    return self.nodes.get(n, [])

  def check(self, update):
    if len(update) <= 1:
      return True
    
    for i in range(len(update)):
      # everything in ant has to be before update[i]
      ant = self.antecedents(update[i])
      for j in range(i, len(update)):
        if update[j] in ant:
          return False
      
    return True
  
  def fix(self, update, idx=0):
    if self.check(update[idx:]):
      return update
    
    current = update[idx]
    must_be_before_current = set(self.antecedents(current))
    move_before, leave_alone = partition_list(update[idx+1:], lambda x: x in must_be_before_current)
    move_before = self.fix(move_before, 0)

    if len(move_before) > 0:
      #print(f"Fixing {update[idx:]} by moving {move_before} before {current}")
      update = update[:idx] + move_before + [current] + leave_alone

      return self.fix(update, idx + 1)
    else:
      return self.fix(update, idx + 1)

  def __str__(self):
    return f"RulesGraph({self.nodes})"

def part1(fname):
  rules_re = re.compile(r"(\d+)\|(\d+)")
  rules = []
  checksum = 0

  with open(fname) as f:
    lines = f.readlines()
    for line in lines:
      if line == '\n':
        break

      m = rules_re.match(line)
      if m:
        rules.append(Rule(int(m.group(1)), int(m.group(2))))

    for update in lines[len(rules) + 1:]:
      update = list(map(int, update.strip().split(',')))
      if check_rules(rules, update):
        #print(f"Valid: {update}")
        checksum += update[int(len(update) / 2)]

  return checksum


def part2(fname):
  rules_re = re.compile(r"(\d+)\|(\d+)")
  rules = []
  checksum = 0

  with open(fname) as f:
    lines = f.readlines()
    for line in lines:
      if line == '\n':
        break

      m = rules_re.match(line)
      if m:
        rules.append(Rule(int(m.group(1)), int(m.group(2))))
    
    g = RulesGraph(rules)
    #print(g)

    for update in lines[len(rules) + 1:]:
      update = list(map(int, update.strip().split(',')))
      if g.check(update):
        continue
      else:
        new_update = g.fix(update)
        #print(f"Fixed: {update} -> {new_update}")
        update = new_update
        checksum += update[int(len(update) / 2)]

  return checksum

if __name__ == "__main__":
  #print(part1("day5/day5-input-easy.txt"))
  #print(part1("day5/day5-input.txt"))

  print(part2("day5/day5-input-easy.txt"))
  print(part2("day5/day5-input.txt"))