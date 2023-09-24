import { Turnstile } from '@marsidev/react-turnstile'
import Button from '@mui/material/Button'
import { useRef, useState, memo } from 'react'
import type { TurnstileInstance } from '@marsidev/react-turnstile'
import useSWR from "swr";
import { ApiCaptcha } from '@/apis/model';
import Alert from '@mui/material/Alert';
import Skeleton from '@mui/material/Skeleton';

interface prop {
    onSuccess: ((token: string) => void)
}

function CaptchaWidget({ onSuccess }: prop) {
    const ref = useRef<TurnstileInstance>(null)
    const [key, setKey] = useState(1)
    const { data, error, isLoading } = useSWR<ApiCaptcha>(import.meta.env.VITE_APIADDR + '/api/v1/captcha')

    if (error) {
        console.warn(error)
        return <Alert severity="warning">{String(error)}</Alert>
    }
    if (isLoading) {
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
            <Turnstile siteKey={data?.data.siteKey ?? ""} key={key} onSuccess={onSuccess} ref={ref} scriptOptions={{ async: true }} />
            <Button onClick={() => setKey(key + 1)}>刷新验证码</Button>
        </>
    )
}

const CaptchaWidgetMemo = memo(CaptchaWidget)

export default CaptchaWidgetMemo