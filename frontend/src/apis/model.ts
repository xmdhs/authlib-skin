export interface tokenData {
    accessToken: string
    selectedProfile: {
        name: string
        uuid: string
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

interface user {
    uid: string
    uuid: string
}

export type ApiUser = Api<user>
