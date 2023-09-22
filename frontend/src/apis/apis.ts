import type { tokenData } from '@/apis/model'

export async function login(username: string, password: string) {
    const v = await fetch(import.meta.env.VITE_APIADDR + "/api/yggdrasil/authserver/authenticate", {
        method: "POST",
        body: JSON.stringify({
            "username": username,
            "password": password,
        })
    })
    const data = await v.json()
    if (!v.ok) {
        throw data?.errorMessage
    }
    return data as tokenData
}

export async function register(email: string, username: string, password: string) {
    const v = await fetch(import.meta.env.VITE_APIADDR + "/api/v1/user/reg", {
        method: "POST",
        body: JSON.stringify({
            "Email": email,
            "Password": password,
            "Name": username
        })
    })
    const data = await v.json()
    if (!v.ok) {
        throw data
    }
    return
}