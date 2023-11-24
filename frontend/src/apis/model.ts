export interface tokenData {
    token: string
    name: string
    uuid: string
}

export interface Api<T> {
    code: number
    msg: string
    data: T
}

export interface List<T> {
    total: number
    list: T[]
}


interface captcha {
    type: string
    siteKey: string
}


export interface ApiUser {
    uid: string
    uuid: string
    is_admin: boolean
}

export interface YggProfile {
    name: string
    properties: {
        name: string
        value: string
    }[]
}

export interface ApiConfig {
    captcha: captcha
    AllowChangeName: boolean
    serverName: string
    NeedEmail: boolean
    AllowDomain: string[]
}

export interface UserInfo {
    uid: number
    uuid: string
    is_admin: boolean
    is_disable: boolean
    email: string
    reg_ip: string
    name: string
}

export interface EditUser {
    email?: string
    name?: string
    password?: string
    is_admin?: boolean
    is_disable?: boolean
    del_textures?: boolean
}