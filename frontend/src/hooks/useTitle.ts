import { getConfig } from '@/apis/apis'
import { useTitle as auseTitle, useRequest } from 'ahooks'
import { useEffect } from 'react'

export default function useTitle(title: string) {
    const { data, error } = useRequest(getConfig, {
        cacheKey: "/api/v1/config",
        staleTime: 60000,
    })
    useEffect(() => {
        error && console.warn(error)
    }, [error])
    auseTitle(title + " - " + data?.serverName ?? "", {
        restoreOnUnmount: true
    })
}