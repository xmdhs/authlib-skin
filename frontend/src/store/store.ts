import { atomWithStorage } from 'jotai/utils'

export const token = atomWithStorage("token", "")
export const user = atomWithStorage("username", {
    name: "",
    uuid: ""
})