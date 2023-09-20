interface tokenData {
    accessToken: string
    selectedProfile: {
        name: string
    }
}

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
