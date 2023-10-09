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
import { useSetAtom } from 'jotai';
import { token, user } from '@/store/store'
import { login } from '@/apis/apis'
import { Link as RouterLink, useNavigate } from "react-router-dom";
import Loading from '@/components/Loading'
import CheckInput, { refType } from '@/components/CheckInput'
import useTitle from '@/hooks/useTitle';
import CaptchaWidget from '@/components/CaptchaWidget';
import type { refType as CaptchaWidgetRef } from '@/components/CaptchaWidget'
import { ApiErr } from '@/apis/error';



export default function SignIn() {
    const [err, setErr] = useState("");
    const [loading, setLoading] = useState(false);
    const setToken = useSetAtom(token)
    const setUserInfo = useSetAtom(user)
    const checkList = React.useRef<Map<string, refType>>(new Map<string, refType>())
    const navigate = useNavigate();
    useTitle("登录")
    const captchaRef = React.useRef<CaptchaWidgetRef>(null)
    const [captchaToken, setCaptchaToken] = useState("");

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        const data = new FormData(event.currentTarget);
        const postData = {
            email: data.get('email')?.toString(),
            password: data.get('password')?.toString(),
        }
        if (!Array.from(checkList.current.values()).map(v => v.verify()).reduce((p, v) => (p == true) && (v == true))) {
            return
        }

        if (loading) return
        setLoading(true)
        login(postData.email!, postData.password ?? "", captchaToken).
            then(v => {
                if (!v) return
                setToken(v.token)
                setUserInfo({
                    uuid: v.uuid,
                    name: v.name,
                })
                navigate("/profile")
            }).
            catch(v => {
                captchaRef.current?.reload()

                if (v instanceof ApiErr) {
                    switch (v.code) {
                        case 6:
                            setErr("错误的密码")
                            break
                        case 9:
                            setErr("用户已被禁用")
                            break
                    }
                    return
                }
                setErr(String(v))
                console.warn(v)
            }).
            finally(() => setLoading(false))

    };


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
                    <CheckInput
                        ref={(dom) => {
                            dom && checkList.current.set("1", dom)
                        }}
                        checkList={[
                            {
                                errMsg: "需为邮箱",
                                reg: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
                            }
                        ]}
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
                    <CaptchaWidget ref={captchaRef} onSuccess={setCaptchaToken} />
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
                            {/* <Link href="#" variant="body2">
                                忘记密码？
                            </Link> */}
                        </Grid>
                        <Grid item>
                            <Link component={RouterLink} to="/register" variant="body2">
                                {"注册"}
                            </Link>
                        </Grid>
                    </Grid>
                </Box>
            </Box>
            <Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'center' }} open={err !== ""}>
                <Alert onClose={() => setErr("")} severity="error">{err}</Alert>
            </Snackbar>
            {loading && <Loading />}
        </Container>
    );
}