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


interface captcha {
    type: string
    siteKey: string
}


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

export interface ApiConfig {
    captcha: captcha
    AllowChangeName: boolean
}