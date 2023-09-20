import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import { useState } from 'react';
import { checkEmail } from '@/utils/email';
import { Link as RouterLink } from "react-router-dom";



export default function SignUp() {
    const [emailErr, setEmailErr] = useState("");
    const [passErr, setpassErr] = useState("");



    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        setEmailErr("")
        setpassErr("")
        event.preventDefault();
        const data = new FormData(event.currentTarget);
        const d = {
            email: data.get('email')?.toString(),
            password: data.get('password')?.toString(),
            password1: data.get('password1')?.toString(),
            username: data.get("username")?.toString()
        }
        if (!checkEmail(d.email ?? "")) {
            setEmailErr("需要为邮箱")
            return
        }
        if (d.password != d.password1) {
            setpassErr("密码不一致")
            return
        }
        if (d.password == "") {
            setpassErr("密码为空")
            return
        }
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
                            <TextField
                                error={emailErr != ""}
                                helperText={emailErr}
                                required
                                fullWidth
                                name="email"
                                label="邮箱"
                                autoComplete="email"
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <TextField
                                required
                                fullWidth
                                name="username"
                                label="角色名"
                                autoComplete="email"
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <TextField
                                error={passErr != ""}
                                helperText={passErr}
                                required
                                fullWidth
                                label="密码"
                                type="password"
                                name="password"
                                autoComplete="new-password"
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <TextField
                                required
                                fullWidth
                                label="确认密码"
                                type="password"
                                name="password1"
                                autoComplete="new-password"
                            />
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
        </Container>
    );
}