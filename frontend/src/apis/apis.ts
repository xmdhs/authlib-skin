import type { tokenData, ApiUser, ApiServerInfo, YggProfile, ApiConfig, List, UserInfo } from '@/apis/model'
import { apiGet } from '@/apis/utils'

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
        throw new Error(data?.errorMessage)
    }
    return data as tokenData
}

export async function register(email: string, username: string, password: string, captchaToken: string) {
    const v = await fetch(import.meta.env.VITE_APIADDR + "/api/v1/user/reg", {
        method: "PUT",
        body: JSON.stringify({
            "Email": email,
            "Password": password,
            "Name": username,
            "CaptchaToken": captchaToken
        })
    })
    return await apiGet(v)
}

export async function userInfo(token: string) {
    if (token == "") return
    const v = await fetch(import.meta.env.VITE_APIADDR + "/api/v1/user", {
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    return await apiGet<ApiUser>(v)
}


export async function serverInfo() {
    const v = await fetch(import.meta.env.VITE_APIADDR + "/api/yggdrasil")
    return await v.json() as ApiServerInfo
}

export async function yggProfile(uuid: string) {
    if (uuid == "") return
    const v = await fetch(import.meta.env.VITE_APIADDR + "/api/yggdrasil/sessionserver/session/minecraft/profile/" + uuid)
    const data = await v.json()
    if (!v.ok) {
        throw new Error(data?.errorMessage)
    }
    return data as YggProfile
}

export async function upTextures(uuid: string, token: string, textureType: 'skin' | 'cape', model: 'slim' | '', file: File) {
    const f = new FormData()
    f.set("file", file)
    f.set("model", model)

    const r = await fetch(import.meta.env.VITE_APIADDR + "/api/yggdrasil/api/user/profile/" + uuid + "/" + textureType, {
        method: "PUT",
        body: f,
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    if (r.status != 204) {
        throw new Error("上传失败 " + String(r.status))
    }
}

export async function changePasswd(old: string, newpa: string, token: string) {
    const r = await fetch(import.meta.env.VITE_APIADDR + "/api/v1/user/password", {
        method: "POST",
        body: JSON.stringify({
            "old": old,
            "new": newpa
        }),
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    return await apiGet<unknown>(r)
}

export async function getConfig() {
    const r = await fetch(import.meta.env.VITE_APIADDR + "/api/v1/config")
    return await apiGet<ApiConfig>(r)
}

export async function changeName(name: string, token: string) {
    const r = await fetch(import.meta.env.VITE_APIADDR + "/api/v1/user/name", {
        method: "POST",
        body: JSON.stringify({
            "name": name,
        }),
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    return await apiGet<unknown>(r)
}

export async function ListUser(page: number, token: string, email: string, name: string) {
    const u = new URL(import.meta.env.VITE_APIADDR + "/api/v1/admin/users")
    u.searchParams.set("page", String(page))
    u.searchParams.set("email", email)
    u.searchParams.set("name", name)
    const r = await fetch(u.toString(), {
        method: "GET",
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    return await apiGet<List<UserInfo>>(r)
}