import TextField from '@mui/material/TextField';
import { useState, useImperativeHandle, forwardRef } from 'react';
import type { TextFieldProps } from '@mui/material/TextField';

export type refType = {
    verify: () => boolean
}

type prop = {
    checkList: {
        errMsg: string
        reg: RegExp
    }[]
} & Omit<Omit<TextFieldProps, 'error'>, 'helperText'>

export const CheckInput = forwardRef<refType, prop>(({ required, checkList, ...textFied }, ref) => {
    const [err, setErr] = useState("");

    const verify = () => {
        return err == ""
    }

    useImperativeHandle(ref, () => {
        return {
            verify
        }
    })

    const onChange = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const value = event.target.value
        if (required && (!value || value == "")) {
            setErr("此项必填")
            return
        }
        for (const v of checkList) {
            if (!v.reg.test(value)) {
                setErr(v.errMsg)
                return
            }
        }
        setErr("")
    }



    return <TextField
        error={err != ""}
        onChange={onChange}
        helperText={err}
        required={required}
        {...textFied}
    />
})

export default CheckInput