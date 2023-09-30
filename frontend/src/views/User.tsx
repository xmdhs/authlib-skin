


export default function User(){
    return (<>
    </>)

}

// import { token, username } from "@/store/store"
// import { useRequest } from "ahooks"
// import { useAtomValue } from "jotai"
// import { userInfo } from '@/apis/apis'


// export default function User() {
//     const nowToken = useAtomValue(token)
//     const nowUsername = useAtomValue(username)

//     const { data, error } = useRequest(() => userInfo(nowToken), {
//         refreshDeps: [nowToken],
//         cacheKey: "/api/v1/user/reg",
//         cacheTime: 10000
//     })


//     return (
//         <>
//             <p>你好: {nowUsername}</p>
//             <p>token: {nowToken} </p>
//             {error && String(error)}
//             {!error && <>
//                 <p>uid: {data?.data.uid}</p>
//                 <p>uuid: {data?.data.uuid}</p>
//             </>}
//         </>
//     )
// }