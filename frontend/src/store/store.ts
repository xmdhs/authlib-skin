import { atom } from 'jotai'
import { atomWithStorage } from 'jotai/utils'

export const token = atomWithStorage("token", "")
export const user = atomWithStorage("username", {
    name: "",
    uuid: ""
})

export const LayoutAlertErr = atom("")
