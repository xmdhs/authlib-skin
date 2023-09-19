import { atomWithStorage } from 'jotai/utils'

export const token = atomWithStorage("token", "")
export const username = atomWithStorage("username", "")