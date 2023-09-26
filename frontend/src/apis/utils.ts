export async function apiWrapGet<T>(v: Response) {
    const data = await v.json()
    if (!v.ok) {
        throw data.msg
    }
    return data as T
}
