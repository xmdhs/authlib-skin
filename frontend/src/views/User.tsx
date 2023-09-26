import { token, username } from "@/store/store"
import { useAtomValue } from "jotai"

export default function User() {
    const nowToken = useAtomValue(token)
    const nowUsername = useAtomValue(username)


    return (
        <>
            <p>你好: {nowUsername}</p>
            <p>token: {nowToken} </p>
        </>
    )
}