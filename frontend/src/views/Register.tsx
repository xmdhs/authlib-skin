import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import { Link as RouterLink } from "react-router-dom";
import { register } from '@/apis/apis'
import CheckInput, { refType } from '@/components/CheckInput'
import { useRef, useState } from 'react';
import Alert from '@mui/material/Alert';
import Snackbar from '@mui/material/Snackbar';
import Loading from '@/components/Loading'
import { useNavigate } from "react-router-dom";
import CaptchaWidget from '@/components/CaptchaWidget';
import type { refType as CaptchaWidgetRef } from '@/components/CaptchaWidget'
import useTitle from '@/hooks/useTitle';
import { ApiErr } from '@/apis/error';
import { useSetAtom } from 'jotai';
import { token, user } from '@/store/store';

export default function SignUp() {
    const [regErr, setRegErr] = useState("");
    const navigate = useNavigate();
    const [captchaToken, setCaptchaToken] = useState("");
    const captchaRef = useRef<CaptchaWidgetRef>(null)
    const [loading, setLoading] = useState(false);
    useTitle("注册")
    const setToken = useSetAtom(token)
    const setUserInfo = useSetAtom(user)


    const checkList = React.useRef<Map<string, refType>>(new Map<string, refType>())


    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (loading) return
        const data = new FormData(event.currentTarget);
        const d = {
            email: data.get('email')?.toString(),
            password: data.get('password')?.toString(),
            username: data.get("username")?.toString()
        }
        if (!Array.from(checkList.current.values()).map(v => v.verify()).reduce((p, v) => (p == true) && (v == true))) {
            return
        }
        if (captchaToken == "") {
            setRegErr("验证码无效")
            return
        }
        setLoading(true)
        register(d.email ?? "", d.username ?? "", d.password ?? "", captchaToken).
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
                console.warn(v)
                if (v instanceof ApiErr) {
                    switch (v.code) {
                        case 10:
                            setRegErr("验证码错误")
                            return    
                        case 3:
                            setRegErr("邮箱已存在")
                            return
                        case 7:
                            setRegErr("用户名已存在")
                            return
                    }
                }
                setRegErr(String(v))
            }).
            finally(() => setLoading(false))
    };

    return (
        <Container component="main" maxWidth="xs">
            <CssBaseline />
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
                    注册
                </Typography>
                <Box component="form" noValidate onSubmit={handleSubmit} sx={{ mt: 3 }}>
                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            <CheckInput
                                checkList={[
                                    {
                                        errMsg: "需为邮箱",
                                        reg: /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
                                    }
                                ]}
                                required
                                fullWidth
                                name="email"
                                label="邮箱"
                                autoComplete="email"
                                ref={(dom) => {
                                    dom && checkList.current.set("1", dom)
                                }}
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <CheckInput
                                ref={(dom) => {
                                    dom && checkList.current.set("2", dom)
                                }}
                                checkList={[
                                    {
                                        errMsg: "长度在 3-16 之间",
                                        reg: /^.{3,16}$/
                                    }
                                ]}
                                required
                                fullWidth
                                name="username"
                                label="角色名"
                                autoComplete="username"
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <CheckInput
                                ref={(dom) => {
                                    dom && checkList.current.set("3", dom)
                                }}
                                checkList={[
                                    {
                                        errMsg: "长度在 6-50 之间",
                                        reg: /^.{6,50}$/
                                    }
                                ]}
                                required
                                fullWidth
                                label="密码"
                                type="password"
                                name="password"
                                autoComplete="new-password"
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <CaptchaWidget ref={captchaRef} onSuccess={setCaptchaToken} />
                        </Grid>
                    </Grid>
                    <Button
                        type="submit"
                        fullWidth
                        variant="contained"
                        sx={{ mt: 3, mb: 2 }}
                    >
                        注册
                    </Button>
                    <Grid container justifyContent="flex-end">
                        <Grid item>
                            <Link component={RouterLink} to={"/login"} variant="body2">
                                登录
                            </Link>
                        </Grid>
                    </Grid>
                </Box>
            </Box>
            <Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'center' }} open={regErr !== ""}>
                <Alert onClose={() => setRegErr("")} severity="error">{regErr}</Alert>
            </Snackbar>
            {loading && <Loading />}
        </Container>
    );
}