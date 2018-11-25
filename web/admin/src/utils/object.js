const setPath = ([key, ...next], value, obj) => {
  if (next.length === 0) {
    return { ...obj, [key]: value }
  }
  return { ...obj, [key]: setPath(next, value, obj[key]) }
}

export const set = (path, value, obj) => setPath(path.split('.'), value, obj)

const unsetPath = ([key, ...next], obj) => {
  const { [key]: value, ...rest } = obj
  if (next.length === 0) {
    return rest
  }
  return { ...rest, [key]: unsetPath(next, value) }
}

export const unset = (path, obj) => unsetPath(path.split('.'), obj)
