import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';
import { useRequest, useTitle } from 'ahooks';
import { getConfig } from '@/apis/apis';
import { useEffect, useRef, useState } from 'react';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';
import CaptchaWidget from '@/components/CaptchaWidget';
import type { refType as CaptchaWidgetRef } from '@/components/CaptchaWidget'
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';
import DialogActions from '@mui/material/DialogActions';
import { useNavigate } from "react-router-dom";
import { ApiErr } from '@/apis/error';
import Loading from '@/components/Loading';

export default function SendEmail({ title, anyEmail = false, sendService }: { title: string, anyEmail?: boolean, sendService: (email: string, captchaToken: string) => Promise<unknown> }) {
    const [err, setErr] = useState("");
    const [domain, setDomain] = useState("");
    const [email, setEmail] = useState("")
    const captchaRef = useRef<CaptchaWidgetRef>(null)
    const [captchaToken, setCaptchaToken] = useState("");
    const [open, setOpen] = useState(false);
    useTitle(title)
    const navigate = useNavigate();
    const [helperText, setHelperText] = useState("")
    const [loading, setLoading] = useState(false);

    const server = useRequest(getConfig, {
        cacheKey: "/api/v1/config",
        staleTime: 60000,
        onError: e => {
            console.warn(e)
            setErr(String(e))
        }
    })

    useEffect(() => {
        if (server.data?.AllowDomain.length != 0) {
            setDomain(server.data?.AllowDomain[0] ?? "")
        }
    }, [server.data?.AllowDomain])

    const emailonChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setEmail(e.target.value)
        if (e.target.value == "") {
            setHelperText("邮箱不得为空")
        }
        setHelperText("")
    }

    const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        if (email == "") {
            setHelperText("邮箱不得为空")
        }
        const sendEmail = (() => {
            if (domain != "") {
                return `${email}@${domain}`
            }
            return email
        })()

        if (!/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(sendEmail)) {
            setHelperText("邮箱格式错误")
            return
        }

        if (server.data?.captcha.type != "" && captchaToken == "") {
            return
        }
        setLoading(true)
        sendService(sendEmail, captchaToken).then(() => setOpen(true)).catch(e => {
            captchaRef.current?.reload()
            console.warn(e)
            if (e instanceof ApiErr) {
                switch (e.code) {
                    case 10:
                        setErr("验证码错误")
                        return
                    case 11:
                        setErr("暂时无法对此邮箱发送邮件")
                        return
                }
            }
            setErr(String(e))
        }).finally(() => setLoading(false))

    }

    const handleClose = () => {
        navigate("/")
    }


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
                    {title}
                </Typography>
                <Box component="form" noValidate onSubmit={onSubmit} sx={{ mt: 3 }}>
                    <Grid container spacing={2}>
                        <Grid item xs={12} sx={{ display: 'grid', columnGap: '3px', gridTemplateColumns: "1fr auto" }}>
                            <TextField fullWidth
                                required
                                name="email"
                                label="邮箱"
                                value={email}
                                helperText={helperText}
                                error={helperText != ""}
                                onChange={emailonChange}
                            />
                            {
                                server.data?.AllowDomain.length != 0 && !anyEmail &&
                                <FormControl>
                                    <InputLabel>域名</InputLabel>
                                    <Select label="域名" value={domain} onChange={v => setDomain(v.target.value)}>
                                        {server.data?.AllowDomain.map(v => <MenuItem value={v}>@{v}</MenuItem>)}
                                    </Select>
                                </FormControl>
                            }
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
                        继续
                    </Button>
                </Box>
            </Box>
            <Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'center' }} open={err !== ""}>
                <Alert onClose={() => setErr("")} severity="error">{err}</Alert>
            </Snackbar>
            <Dialog open={open}>
                <DialogTitle>邮件已发送</DialogTitle>
                <DialogContent>
                    <Typography>请到收件箱（或垃圾箱）点击验证链接以继续。</Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose}>返回首页</Button>
                </DialogActions>
            </Dialog>
            {loading && <Loading />}
        </Container>
    )
}