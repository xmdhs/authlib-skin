import * as React from 'react';
import { useState } from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';
import { loadable } from "jotai/utils"
import { atom, useAtom } from "jotai"
import Backdrop from '@mui/material/Backdrop';
import CircularProgress from '@mui/material/CircularProgress';

const loginData = atom({ email: "", password: "" })
const loginErr = atom("")

const fetchReg = loadable(atom(async (get) => {
    const ld = get(loginData)
    const v = await fetch(import.meta.env.VITE_APIADDR + "/api/yggdrasil/authserver/authenticate", {
        method: "POST",
        body: JSON.stringify({
            "username": ld.email,
            "password": ld.password,
        })
    })
    const rData = await v.json()
    return rData
}))


const ToLogin = React.memo(() => {
    const [, setErr] = useAtom(loginErr);
    const [value] = useAtom(fetchReg)
    if (value.state === 'hasError') {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        setErr(String(value.error as any))
        console.warn(value.error)
        return <></>
    }
    if (value.state === 'loading') {
        return (
            <>
                <Backdrop
                    sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
                    open={true}
                >
                    <CircularProgress color="inherit" />
                </Backdrop>
            </>
        )
    }
    return (
        <>
            <p>{JSON.stringify(value.data)}</p>
        </>
    )
})



export default function SignIn() {
    const [emailErr, setEmailErr] = useState("");
    const [err, setErr] = useAtom(loginErr);
    const [login, setLogin] = useState(false);
    const [, setloginData] = useAtom(loginData);

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const data = new FormData(event.currentTarget);
        const postData = {
            email: data.get('email')?.toString(),
            password: data.get('password')?.toString(),
        }
        if (!postData.email?.includes("@")) {
            setEmailErr("需要为邮箱")
            return
        }
        setloginData({
            email: postData.email,
            password: postData.password ?? ""
        })
        setLogin(true)
    };

    const closeAlert = () => {
        setLogin(false)
        setErr("")
    }

    return (
        <Container component="main" maxWidth="xs">
            <Box
                sx={{
                    marginTop: 8,
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                }}
            >
                <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
                    <LockOutlinedIcon />
                </Avatar>
                <Typography component="h1" variant="h5">
                    登录
                </Typography>
                <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
                    <TextField
                        error={emailErr != ""}
                        helperText={emailErr}
                        margin="normal"
                        fullWidth
                        id="email"
                        label="邮箱"
                        name="email"
                        autoComplete="email"
                        autoFocus
                    />
                    <TextField
                        margin="normal"
                        fullWidth
                        name="password"
                        label="密码"
                        type="password"
                        id="password"
                        autoComplete="current-password"
                    />
                    <Button
                        type="submit"
                        fullWidth
                        variant="contained"
                        sx={{ mt: 2, mb: 2 }}
                    >
                        登录
                    </Button>
                    <Grid container>
                        <Grid item xs>
                            <Link href="#" variant="body2">
                                忘记密码？
                            </Link>
                        </Grid>
                        <Grid item>
                            <Link href="#" variant="body2">
                                {"注册"}
                            </Link>
                        </Grid>
                    </Grid>
                </Box>
            </Box>
            <Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'center' }} open={err !== ""} onClose={closeAlert}  >
                <Alert onClose={closeAlert} severity="error">{err}</Alert>
            </Snackbar>
            {login && <ToLogin></ToLogin>}
        </Container>
    );
}