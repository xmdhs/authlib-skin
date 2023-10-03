import { serverInfo } from '@/apis/apis'
import { useTitle as auseTitle, useRequest } from 'ahooks'
import { useEffect } from 'react'

export default function useTitle(title: string) {
    const { data, error } = useRequest(serverInfo, {
        cacheKey: "/api/yggdrasil",
        staleTime: 60000,
    })
    useEffect(() => {
        error && console.warn(error)
    }, [error])
    auseTitle(title + " - " + data?.meta.serverName ?? "", {
        restoreOnUnmount: true
    })
}