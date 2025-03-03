import re

def partial_eq(params, target, with_concat_operator):
  if len(params) == 1:
    if params[0] == target:
      return str(target)
    else:
      return None

  if len(params) >= 2:
    plus_attempt = partial_eq([params[0] + params[1]] + params[2:], target, with_concat_operator)
    if plus_attempt:
      return f"{params[0]} + {plus_attempt}"
  
    times_attempt = partial_eq([params[0] * params[1]] + params[2:], target, with_concat_operator)
    if times_attempt:
      return f"{params[0]} * {times_attempt}"
    
    if with_concat_operator:
      concat = int(str(params[0]) + str(params[1]))
      concat_attempt = partial_eq([concat] + params[2:], target, with_concat_operator)
      if concat_attempt:
        return f"{params[0]} | {times_attempt}"
  
  return None

def check_equations(fname, with_concat_operator):
  checksum = 0

  with open(fname) as f:
    lines = f.readlines()
    for line in lines:
      result, remainder = line.split(":")
      result = int(result.strip())
      params = list(map(int, remainder.split()))

      eq = partial_eq(params, result, with_concat_operator)
      if eq:
        checksum += result
        #print(eq)

  return checksum
  

if __name__ == '__main__':
  print(check_equations("day7/day7-input-easy.txt", False))
  print(check_equations("day7/day7-input.txt", False))

  print(check_equations("day7/day7-input-easy.txt", True))
  print(check_equations("day7/day7-input.txt", True))