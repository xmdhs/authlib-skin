import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import TextField from '@mui/material/TextField';
import { useTitle } from 'ahooks';
import { useEffect, useState } from 'react';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';
import Loading from '@/components/Loading';
import { produce } from 'immer';
import { forgotPassWord } from '@/apis/apis';
import { useNavigate } from 'react-router-dom';

export default function Forgot() {
    const [err, setErr] = useState("")
    useTitle("重设密码")
    const [passerr, setPasserr] = useState("")
    const [pass, setPass] = useState({
        pass1: "",
        pass2: "",
    })
    const [load, setLoad] = useState(false)
    const [email, setEmail] = useState("")
    const [code, setCode] = useState("")
    const navigate = useNavigate();

    useEffect(() => {
        if (pass.pass1 != pass.pass2 && pass.pass2 != "") {
            setPasserr("密码不相等")
            return
        }
        setPasserr("")
    }, [pass.pass1, pass.pass2])

    const u = new URL(location.href)

    useEffect(() => {
        setEmail(u.searchParams.get("email") ?? "")
        setCode(u.searchParams.get("code") ?? "")
    }, [u.searchParams])


    const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        setLoad(true)
        forgotPassWord(email, code, pass.pass1).then(() => {
            navigate("/")
        }).catch(e => {
            setErr(String(e))
        }).finally(() => { setLoad(false) })
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
                    找回密码
                </Typography>
                <Box component="form" noValidate onSubmit={onSubmit} sx={{ mt: 3 }}>
                    <Grid container spacing={2}>
                        <Grid item xs={12}>
                            <TextField
                                margin='dense'
                                fullWidth
                                label="新密码"
                                type="password"
                                required
                                onChange={p => setPass(produce(v => { v.pass1 = p.target.value }))}
                                autoComplete="new-password"
                            />
                        </Grid>
                        <Grid item xs={12}>
                            <TextField
                                margin='dense'
                                fullWidth
                                label="确认新密码"
                                type="password"
                                required
                                error={passerr != ""}
                                helperText={passerr}
                                onChange={p => setPass(produce(v => { v.pass2 = p.target.value }))}
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
                        提交
                    </Button>
                </Box>
            </Box>
            <Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'center' }} open={err !== ""}>
                <Alert onClose={() => setErr("")} severity="error">{err}</Alert>
            </Snackbar>

            {load && <Loading />}
        </Container>
    )
}