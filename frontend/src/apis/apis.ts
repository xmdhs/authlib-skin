import type { tokenData, ApiUser, YggProfile, ApiConfig, List, UserInfo, EditUser } from '@/apis/model'
import { apiGet } from '@/apis/utils'
import root from '@/utils/root'

export async function login(email: string, password: string, captchaToken: string) {
    const v = await fetch(root() + "/api/v1/user/login", {
        method: "POST",
        body: JSON.stringify({
            "email": email,
            "password": password,
            "CaptchaToken": captchaToken
        })
    })
    return await apiGet<tokenData>(v)
}

export async function register(email: string, username: string, password: string, captchaToken: string, code: string) {
    const v = await fetch(root() + "/api/v1/user/reg", {
        method: "POST",
        body: JSON.stringify({
            "Email": email,
            "Password": password,
            "Name": username,
            "CaptchaToken": captchaToken,
            "EmailJwt": code,
        })
    })
    return await apiGet<tokenData>(v)
}

export async function userInfo(token: string) {
    if (token == "") return
    const v = await fetch(root() + "/api/v1/user", {
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    return await apiGet<ApiUser>(v)
}

export async function yggProfile(uuid: string) {
    if (uuid == "") return
    const v = await fetch(root() + "/api/yggdrasil/sessionserver/session/minecraft/profile/" + uuid)
    const data = await v.json()
    if (!v.ok) {
        throw new Error(data?.errorMessage)
    }
    return data as YggProfile
}

export async function upTextures(token: string, textureType: 'skin' | 'cape', model: 'slim' | '', file: File) {
    const f = new FormData()
    f.set("file", file)
    f.set("model", model)

    const r = await fetch(root() + "/api/v1/user/skin/" + textureType, {
        method: "PUT",
        body: f,
        headers: {
            "Authorization": "Bearer " + token
        }
    })
    return await apiGet<unknown>(r)
}

export async function changePasswd(old: string, newpa: string, token: string) {
    const r = await fetch(root() + "/api/v1/user/password", {
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
    const r = await fetch(root() + "/api/v1/config")
    return await apiGet<ApiConfig>(r)
}

export async function changeName(name: string, token: string) {
    const r = await fetch(root() + "/api/v1/user/name", {
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
    const u = new URL(root() + "/api/v1/admin/users")
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

export async function editUser(u: EditUser, token: string, uid: string) {
    const r = await fetch(root() + "/api/v1/admin/user/" + uid, {
        method: "PATCH",
        headers: {
            "Authorization": "Bearer " + token
        },
        body: JSON.stringify(u)
    })
    return await apiGet<unknown>(r)
}

export async function sendRegEmail(email: string, captchaToken: string) {
    const r = await fetch(root() + "/api/v1/user/reg_email", {
        method: "POST",
        body: JSON.stringify({
            "email": email,
            "captchaToken": captchaToken
        })
    })
    return await apiGet<unknown>(r)
}