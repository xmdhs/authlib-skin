import { Turnstile } from '@marsidev/react-turnstile'
import Button from '@mui/material/Button'
import { useRef, useState, memo, forwardRef, useImperativeHandle } from 'react'
import type { TurnstileInstance } from '@marsidev/react-turnstile'
import { ApiCaptcha } from '@/apis/model';
import Alert from '@mui/material/Alert';
import Skeleton from '@mui/material/Skeleton';
import { useRequest } from 'ahooks';

interface prop {
    onSuccess: ((token: string) => void)
}

export type refType = {
    reload: () => void
}


const CaptchaWidget = forwardRef<refType, prop>(({ onSuccess }, ref) => {
    const Turnstileref = useRef<TurnstileInstance>(null)
    const [key, setKey] = useState(1)
    const { data, error, loading } = useRequest(() => fetch(import.meta.env.VITE_APIADDR + '/api/v1/captcha').then(v => v.json() as Promise<ApiCaptcha>), {
        loadingDelay: 500
    })

    useImperativeHandle(ref, () => {
        return {
            reload: () => {
                setKey(key + 1)
            }
        }
    })


    if (error) {
        console.warn(error)
        return <Alert severity="warning">{String(error)}</Alert>
    }
    if (loading) {
        return <Skeleton variant="rectangular" width={300} height={65} />
    }
    if (data?.code != 0) {
        console.warn(error)
        return <Alert severity="warning">{String(data?.msg)}</Alert>
    }
    if (data.data.type != "turnstile") {
        onSuccess("ok")
        return <></>
    }

    return (
        <>
            <Turnstile siteKey={data?.data.siteKey ?? ""} key={key} onSuccess={onSuccess} ref={Turnstileref} scriptOptions={{ async: true }} />
            <Button onClick={() => setKey(key + 1)}>刷新验证码</Button>
        </>
    )
})

const CaptchaWidgetMemo = memo(CaptchaWidget)

export default CaptchaWidgetMemo