export const createRandomUniqueKey = (label: string, id: string | number): string => {
  return `my-unique-key_${label}_${id}_${new Date().toISOString().replace(/[:.]/g, '-')}`
}