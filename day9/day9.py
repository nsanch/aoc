class FakeFile(object):
  def __init__(self, id, size, kind):
    self.id = id
    self.size = size
    self.kind = kind
  
  def grow(self):
    self.size += 1
  
  def shrink(self, inc=1):
    self.size -= inc

  def __str__(self):
    return f"FakeFile({self.id}, {self.size}, {self.kind})"

class Filesystem(object):
  def __init__(self):
    self.stuff = [] 

  def add(self, thing):
    self.stuff.append(thing)

  def __str__(self):
    ret = ""
    for x in self.stuff:
      ret = ret + str(x.size)
      if x.size <= 0:
        ret = ret + f"{(x.id, x.size)}"
    return ret
  
  def longstr(self):
    ret = ""
    for x in self.stuff:
      ret = ret + (str(x.id) * x.size)
      #if x.size <= 0:
      #  ret = ret + f"{(x.id, x.size)}"
    return ret
  
  def checksum(self):
    ret = 0
    pos = 0
    for x in self.stuff:
      if x.kind == "empty":
        pos += x.size
      else:
        for i in range(x.size):
          ret = ret + pos * x.id
          pos += 1
    return ret
    
  def compact_by_one(self):
    source = None
    source_idx = None
    last_empty = None
    if self.stuff[-1].kind == "empty":
      last_empty = self.stuff[-1]

    for i in range(len(self.stuff) - 1, 0, -1):
      if self.stuff[i].kind == "file":
        source = self.stuff[i]
        source_idx = i
        break
    if not source:
      return False

    destination_block = None
    destination_idx = None
    for idx in range(min(len(self.stuff), source_idx)):
      curr = self.stuff[idx]
      if curr.kind == "empty":
        destination_block = curr
        destination_idx = idx
        break

    if not destination_block:
      return False

    # add an empty block at end if needed
    # grow empty block at end
    # shrink the destination
    # delete the destination if empty
    # add the new block if necessary
    # grow the new block
    # shrink the source
    # delete the source if empty
    if not last_empty:
      last_empty = FakeFile(".", 0, "empty")
      self.stuff.append(last_empty)

    # consider this like moving the empty spot from the middle to the end
    last_empty.grow()
    destination_block.shrink()

    if destination_block.size == 0:
      source_idx -= 1
      del self.stuff[destination_idx]
  
    new_block = None
    if (destination_idx >= 1 and self.stuff[destination_idx - 1].id == source.id):
      new_block = self.stuff[destination_idx - 1]
      new_block.grow()
    else:
      new_block = FakeFile(source.id, 1, "file")
      #print(f"adding at {destination_idx}")
      self.stuff.insert(destination_idx, new_block)
      source_idx += 1
      #print(self.longstr())

    source.shrink()
    if source.size == 0:
      del self.stuff[source_idx]

    return True
  
  def compact_whole_files(self):
    all_files_in_reverse_order = [x for x in self.stuff if x.kind == "file"]
    all_files_in_reverse_order.reverse()

    for source in all_files_in_reverse_order:
      source_idx = self.stuff.index(source)

      # find a `destination` empty spot closest to beginning that can fit `source`
      # shrink the empty spot by the new `source` file size
      # move `source` before `destination`
      # if destination is empty, remove it
      destination = None
      destination_idx = None
      for idx in range(source_idx):
        if self.stuff[idx].kind == "empty" and self.stuff[idx].size >= source.size:
          destination = self.stuff[idx]
          destination_idx = idx
          break

      if destination:
        del self.stuff[source_idx]
        self.stuff.insert(source_idx, FakeFile(".", source.size, "empty"))
        self.stuff.insert(destination_idx, source)
        destination.shrink(source.size)
        if destination.size == 0:
          del self.stuff[destination_idx + 1]

def read_compact_checksum(line):
  fs = Filesystem()
  id = 0
  idx = 0
  while idx < len(line):
    filesize = int(line[idx])
    if filesize > 0:
      fs.add(FakeFile(id, filesize, "file"))
    idx += 1
    id += 1
    if idx < len(line):
      emptysize = int(line[idx])
      if emptysize > 0:
        fs.add(FakeFile(".", emptysize, "empty"))
      idx += 1
  
  return fs

def part1(fname):
  with open(fname) as f:
    l = f.readlines()[0].strip()
    fs = read_compact_checksum(l)
    print(fs)
    print(fs.longstr())
    while fs.compact_by_one():
      #print(fs.longstr()) 
      pass
    return fs.checksum()
  

def part2(fname):
  with open(fname) as f:
    l = f.readlines()[0].strip()
    fs = read_compact_checksum(l)
    print(fs)
    print(fs.longstr())
    fs.compact_whole_files()
    print(fs.longstr())
    return fs.checksum()

if __name__ == "__main__":
  print(part1("day9/day9-input-easy.txt"))
  print(part1("day9/day9-input.txt"))

  print(part2("day9/day9-input-easy.txt"))
  print(part2("day9/day9-input.txt"))