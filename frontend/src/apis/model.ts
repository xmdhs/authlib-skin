export interface tokenData {
    accessToken: string
    selectedProfile: {
        name: string
        id: string
    }
}

export interface Api<T> {
    code: number
    msg: string
    data: T
}

export type ApiErr = Api<unknown>

interface captcha {
    type: string
    siteKey: string
}

export type ApiCaptcha = Api<captcha>

export interface ApiUser {
    uid: string
    uuid: string
    is_admin: boolean
}

export interface ApiServerInfo {
    meta: {
        serverName: string
    }
}

export interface YggProfile {
    name: string
    properties: {
        name: string
        value: string
    }[]
}