import { ApiErr } from "./error"

export async function apiGet<T>(v: Response) {
    type api = { data: T, msg: string, code: number }
    const data = await v.json() as api
    if (!v.ok) {
        throw new ApiErr(data.code, data.msg)
    }
    return data.data
}
