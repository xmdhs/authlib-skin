import { Turnstile } from '@marsidev/react-turnstile'
import Button from '@mui/material/Button'
import { useRef, useState, memo, forwardRef, useImperativeHandle, useEffect } from 'react'
import type { TurnstileInstance } from '@marsidev/react-turnstile'
import Alert from '@mui/material/Alert';
import Skeleton from '@mui/material/Skeleton';
import { useRequest } from 'ahooks';
import { getConfig } from '@/apis/apis';

interface prop {
    onSuccess: ((token: string) => void)
}

export type refType = {
    reload: () => void
}


const CaptchaWidget = forwardRef<refType, prop>(({ onSuccess }, ref) => {
    const Turnstileref = useRef<TurnstileInstance>(null)
    const [key, setKey] = useState(1)
    const { data, error, loading } = useRequest(getConfig, {
        cacheKey: "/api/v1/config",
        staleTime: 600000,
        loadingDelay: 200
    })

    useImperativeHandle(ref, () => {
        return {
            reload: () => {
                setKey(key + 1)
            }
        }
    })
    useEffect(() => {
        if (data?.captcha?.type != "turnstile") {
            onSuccess("ok")
            return
        }
    }, [data?.captcha?.type, onSuccess])


    if (error) {
        console.warn(error)
        return <Alert severity="warning">{String(error)}</Alert>
    }
    if (!data && loading) {
        return <Skeleton variant="rectangular" width={300} height={65} />
    }

    if (data?.captcha.type == "") {
        return <></>
    }

    return (
        <>
            <Turnstile siteKey={data?.captcha?.siteKey ?? ""} key={key} onSuccess={onSuccess} ref={Turnstileref} scriptOptions={{ async: true }} />
            <Button onClick={() => setKey(key + 1)}>刷新验证码</Button>
        </>
    )
})

const CaptchaWidgetMemo = memo(CaptchaWidget)

export default CaptchaWidgetMemo