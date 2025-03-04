def make_graph(lines: list[str]) -> dict[str, list[str]]:
  graph = {}
  for line in lines:
    start, end = line.split('-')
    if start not in graph:
      graph[start] = []
    if end not in graph:
      graph[end] = []
    graph[start].append(end)
    graph[end].append(start)
  return graph

def look_for_groups_of_three(graph: dict[str, list[str]], limit_to_t) -> list[str]:
  ret = set()
  for node in graph:
    if limit_to_t and not node.startswith('t'):
      continue

    neighbors = graph[node]
    for i in range(len(neighbors)):
      for j in range(i + 1, len(neighbors)):
        if neighbors[i] in graph[neighbors[j]]:
          ret.add(','.join(sorted([node, neighbors[i], neighbors[j]])))
  return list(ret)

def embiggen_groups(graph: dict[str, list[str]], groups: set[str]) -> set[str]:
  one_bigger: set[str] = set()
  for group_str in groups:
    group: list[str] = group_str.split(',')
    new_members: set[str] = set(graph[group[0]])
    for i in range(1, len(group)):
      new_members.intersection_update(graph[group[i]])
    for nm in new_members:
      one_bigger.add(','.join(sorted(group + [nm])))
    #one_bigger.update([(group_str + ',' + nm) for nm in new_members])

  if len(one_bigger) == 0:
    return groups
  else:
    one_bigger_list = list(one_bigger)
    print(f"Recursing with {len(one_bigger)} groups of size {one_bigger_list[0].count(',') + 1}")
    return embiggen_groups(graph, one_bigger)

def part1(fname: str) -> str:
  with open(fname) as f:
    lines = f.read().splitlines()

  graph = make_graph(lines)
  groups = sorted(look_for_groups_of_three(graph, True))
  print(f"{fname}: {len(groups)} groups: {groups[:50]}")
  return groups

def part2(fname: str) -> str:
  with open(fname) as f:
    lines = f.read().splitlines()

  graph = make_graph(lines)
  groups = sorted(look_for_groups_of_three(graph, False))

  groups = embiggen_groups(graph, groups)
  
  print(f"{fname}: {len(groups)} groups: {list(groups)[:50]}")
  print(groups.pop())
#  return groups


part1("day23/day23-input-easy.txt")
part1("day23/day23-input.txt")

part2("day23/day23-input-easy.txt")
part2("day23/day23-input.txt")
