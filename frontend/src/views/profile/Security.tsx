import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CardHeader from "@mui/material/CardHeader";
import TextField from "@mui/material/TextField";
import { useEffect, useState } from "react";
import { produce } from 'immer'
import { changeName, changePasswd, getConfig } from "@/apis/apis";
import { useAtom, useAtomValue, useSetAtom } from "jotai";
import { LayoutAlertErr, token, user } from "@/store/store";
import Loading from "@/components/Loading";
import { ApiErr } from "@/apis/error";
import { useNavigate } from "react-router-dom";
import useTitle from "@/hooks/useTitle";
import Box from "@mui/material/Box";
import { useRequest } from "ahooks";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogActions from "@mui/material/DialogActions";

export default function Security() {
    useTitle("账号设置")
    const setLayoutErr = useSetAtom(LayoutAlertErr)

    const { data } = useRequest(getConfig, {
        cacheKey: "/api/v1/config",
        staleTime: 600000,
        onError: e => {
            setLayoutErr(String(e))
        }
    })

    return (<>
        <Box sx={{
            display: "grid", gap: "1em",
            gridTemplateColumns: {
                lg: "1fr 1fr"
            }
        }}>
            <ChangePasswd />
            {data?.AllowChangeName && <ChangeName />}
        </Box>
    </>)
}

function ChangePasswd() {
    const [pass, setPass] = useState({
        old: "",
        pass1: "",
        pass2: "",
    })
    const [err, setErr] = useState("")
    const [oldPassErr, setOldPassErr] = useState(false)
    const [nowToken, setToken] = useAtom(token)
    const [load, setLoad] = useState(false)
    const setLayoutErr = useSetAtom(LayoutAlertErr)
    const setUser = useSetAtom(user)
    const navigate = useNavigate();

    useEffect(() => {
        if (pass.pass1 != pass.pass2 && pass.pass2 != "") {
            setErr("密码不相等")
            return
        }
        setErr("")
    }, [pass.pass1, pass.pass2])

    const handelClick = () => {
        if (pass.pass1 != pass.pass2) return
        if (load) return
        setLoad(true)
        changePasswd(pass.old, pass.pass1, nowToken).catch(e => {
            if (e instanceof ApiErr && e.code == 6) {
                setOldPassErr(true)
                return
            }
            setLayoutErr(String(e))
        }).finally(() => setLoad(false)).then(() => [navigate("/login"), setToken(""), setUser({ name: "", uuid: "" })])
    }


    return (<>
        <Card sx={{ maxWidth: "30em" }}>
            <CardHeader title="更改密码" />
            <CardContent>
                <TextField
                    margin='dense'
                    fullWidth
                    label="旧密码"
                    type="password"
                    required
                    error={oldPassErr}
                    helperText={oldPassErr ? "旧密码错误" : ""}
                    onChange={p => setPass(produce(v => { v.old = p.target.value; return v }))}
                    autoComplete="current-password"
                />
                <TextField
                    margin='dense'
                    fullWidth
                    label="新密码"
                    type="password"
                    required
                    onChange={p => setPass(produce(v => { v.pass1 = p.target.value; return v }))}
                    autoComplete="new-password"
                />
                <TextField
                    margin='dense'
                    fullWidth
                    label="确认新密码"
                    type="password"
                    required
                    error={err != ""}
                    helperText={err}
                    onChange={p => setPass(produce(v => { v.pass2 = p.target.value; return v }))}
                    autoComplete="new-password"
                />
                <Button sx={{ marginTop: "1em" }} onClick={handelClick} variant='contained'>提交</Button>
            </CardContent>
        </Card>
        {load && <Loading />}

    </>)
}

function ChangeName() {
    const [err, setErr] = useState("")
    const [name, setName] = useState("")
    const [open, setOpen] = useState(false)
    const [load, setLoad] = useState(false)
    const nowToken = useAtomValue(token)
    const setUser = useSetAtom(user)

    const handelClick = () => {
        if (name == "") return
        setOpen(true)
    }

    const handleClose = () => {
        setOpen(false)
    }

    const handleSubmit = () => {
        if (load) return
        setLoad(true)
        changeName(name, nowToken).then(() => {
            setName("")
            setUser(v => { return { name: name, uuid: v.uuid } })
        }).catch(e => {
            if (e instanceof ApiErr && e.code == 7) {
                setErr("用户名已存在")
                return
            }
            setErr(String(e))
            console.warn(e)
        }).finally(() => [setLoad(false), setOpen(false)])
    }

    return (<>
        <Card sx={{ height: "min-content" }}>
            <CardHeader title="更改用户名" />
            <CardContent>
                <TextField
                    margin='dense'
                    fullWidth
                    label="新用户名"
                    type='text'
                    required
                    error={err != ""}
                    helperText={err}
                    value={name}
                    onChange={v => setName(v.target.value)}
                    autoComplete="username"
                />
                <Button sx={{ marginTop: "1em" }} onClick={handelClick} variant='contained'>提交</Button>
            </CardContent>
        </Card>
        <Dialog
            open={open}
            onClose={handleClose}
        >
            <DialogTitle>
                确认修改后的用户名
            </DialogTitle>
            <DialogContent>
                <DialogContentText>
                    {`用户名改为`}  {<code> {name} </code>} {`？`}
                </DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose}>取消</Button>
                <Button onClick={handleSubmit} autoFocus>
                    好
                </Button>
            </DialogActions>
        </Dialog>
        {load && <Loading />}
    </>)
}