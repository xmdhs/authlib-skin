import TextField from '@mui/material/TextField';
import { useState, useImperativeHandle, forwardRef } from 'react';
import type { TextFieldProps } from '@mui/material/TextField';
import { useControllableValue } from 'ahooks';

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
    const [value, setValue] = useControllableValue<string>(textFied);


    const check = (value: string) => {
        if (required && (!value || value == "")) {
            setErr("此项必填")
            return false
        }
        for (const v of checkList) {
            if (!v.reg.test(value)) {
                setErr(v.errMsg)
                return false
            }
        }
        setErr("")
        return true
    }

    const verify = () => {
        return check(value)
    }

    useImperativeHandle(ref, () => {
        return {
            verify
        }
    })

    const onChange = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const value = event.target.value
        setValue(value)
        check(value)
    }



    return <TextField
        error={err != ""}
        onChange={onChange}
        helperText={err}
        required={required}
        value={value}
        {...textFied}
    />
})

export default CheckInput