import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import useTitle from '@/hooks/useTitle';
import { useRequest } from 'ahooks';
import { ListUser, editUser } from '@/apis/apis';
import { useEffect, useRef, useState } from 'react';
import { useAtomValue } from 'jotai';
import { token } from '@/store/store';
import TablePagination from '@mui/material/TablePagination';
import Alert from '@mui/material/Alert';
import Snackbar from '@mui/material/Snackbar';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Box from '@mui/material/Box';
import Chip from '@mui/material/Chip';
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';
import DialogActions from '@mui/material/DialogActions';
import { EditUser, UserInfo } from '@/apis/model';
import { produce } from 'immer'
import Checkbox from '@mui/material/Checkbox';
import FormGroup from '@mui/material/FormGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import SkinViewUUID from '@/components/SkinViewUUID';
import Loading from '@/components/Loading';

export default function UserAdmin() {
    useTitle("用户管理")
    const [page, setPage] = useState(1)
    const nowtoken = useAtomValue(token)
    const [err, setErr] = useState("")
    const [email, setEmail] = useState("")
    const [name, setName] = useState("")
    const [open, setOpen] = useState(false);
    const [row, setRow] = useState<UserInfo | null>(null)


    const handleOpen = (row: UserInfo) => {
        setRow(row)
        setOpen(true)
    }

    const uq = new URLSearchParams("/api/v1/admin/users")
    uq.set("page", String(page))
    uq.set("email", email)
    uq.set("name", name)

    const { data, run } = useRequest(ListUser, {
        cacheKey: uq.toString(),
        debounceWait: 300,
        onError: e => {
            setErr(String(e))
        }
    })
    useEffect(() => {
        run(page, nowtoken, email, name)
    }, [page, nowtoken, run, email, name])


    return (<>
        <Paper>
            <Box sx={{ p: "1em", display: "flex", gap: "1em", alignItems: "flex-end" }}>
                <Chip label="前缀筛选" />
                <TextField onChange={v => setEmail(v.target.value)} label="邮箱" variant="standard" />
                <TextField onChange={v => setName(v.target.value)} label="用户名" variant="standard" />
            </Box>
            <TableContainer >
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>邮箱</TableCell>
                            <TableCell>用户名</TableCell>
                            <TableCell>注册 ip</TableCell>
                            <TableCell>uuid</TableCell>
                            <TableCell></TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {data?.list.map((row) => (
                            <TableRow key={row.uid}>
                                <TableCell sx={{ maxWidth: 'min-content' }}>{row.email}</TableCell>
                                <TableCell>{row.name}</TableCell>
                                <TableCell>{row.reg_ip}</TableCell>
                                <TableCell>{row.uuid}</TableCell>
                                <TableCell><Button onClick={() => handleOpen(row)}>编辑</Button></TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
            <TablePagination
                rowsPerPageOptions={[20]}
                component="div"
                count={data?.total ?? 0}
                rowsPerPage={20}
                page={page - 1}
                onPageChange={(_, page) => setPage(page + 1)}
            />
        </Paper >
        <Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'center' }} open={err != ""} >
            <Alert onClose={() => setErr("")} severity="error">{err}</Alert>
        </Snackbar>

        <MyDialog open={open} setOpen={setOpen} row={row} onUpdate={() => run(page, nowtoken, email, name)} />
    </>);
}

interface MyDialogProp {
    open: boolean
    setOpen: (b: boolean) => void
    row: UserInfo | null
    onUpdate: () => void
}

function MyDialog({ open, row, setOpen, onUpdate }: MyDialogProp) {
    const handleClose = () => {
        setOpen(false)
    }
    const [erow, setErow] = useState<EditUser>({
        email: "",
        name: "",
        password: "",
        is_admin: false,
        is_disable: false,
        del_textures: false,
    })
    const [load, setLoad] = useState(false)
    const nowToken = useAtomValue(token)
    const [err, setErr] = useState("")
    const editValue = useRef<EditUser>({});

    useEffect(() => {
        if (!row) return
        setErow({
            email: row.email,
            name: row.name,
            password: "",
            is_admin: row.is_admin,
            is_disable: row.is_disable,
            del_textures: false,
        })
        editValue.current = {}
    }, [row, open])

    const handleOpen = () => {
        if (load) return
        setLoad(true)
        editUser(editValue.current, nowToken, String(row?.uid)).then(() => [setOpen(false), onUpdate(), editValue.current = {}]).finally(() => setLoad(false)).
            catch(e => setErr(String(e)))
    }

    type StringKeys<T> = {
        [K in keyof T]: T[K] extends string ? K : never;
    }[keyof T];

    function handleSetValue(key: StringKeys<Required<EditUser>>, value: string) {
        setErow(produce(v => {
            v[key] = value
            editValue.current[key] = value
        }))
    }

    type BoolKeys<T> = {
        [K in keyof T]: T[K] extends boolean ? K : never;
    }[keyof T];

    function handleSetChecked(key: BoolKeys<Required<EditUser>>, value: boolean) {
        setErow(produce(v => {
            v[key] = value
            editValue.current[key] = value
        }))
    }


    return (<>
        <Dialog open={open}>
            <DialogTitle>修改用户信息</DialogTitle>
            <DialogContent sx={{
                display: "grid", gap: '1em',
                gridTemplateAreas: {
                    md: "'a c' 'b b'",
                    xs: "'a' 'c' 'b'"
                }
            }}>
                <Box sx={{ display: "flex", flexDirection: 'column', gap: '0.5em', gridArea: "a" }}>
                    <TextField
                        margin="dense"
                        label="邮箱"
                        type="email"
                        variant="standard"
                        value={erow?.email}
                        onChange={e => handleSetValue('email', e.target.value)}
                    />
                    <TextField
                        margin="dense"
                        label="用户名"
                        type="text"
                        variant="standard"
                        value={erow?.name}
                        onChange={e => handleSetValue('name', e.target.value)}
                    />
                    <TextField
                        margin="dense"
                        label="密码"
                        type="text"
                        placeholder='（未更改）'
                        variant="standard"
                        value={erow?.password}
                        onChange={e => handleSetValue('password', e.target.value)}
                    />
                </Box>
                <FormGroup row sx={{ gridArea: "b" }}>
                    <FormControlLabel control={<Checkbox checked={erow?.is_admin} onChange={e => handleSetChecked('is_admin', e.target.checked)} />} label="管理权限" />
                    <FormControlLabel control={<Checkbox checked={erow?.is_disable} onChange={e => handleSetChecked('is_disable', e.target.checked)} />} label="禁用" />
                    <FormControlLabel control={<Checkbox checked={erow?.del_textures} onChange={e => handleSetChecked('del_textures', e.target.checked)} />} label="清空材质" />
                </FormGroup>
                <Box sx={{ gridArea: "c" }}>
                    <SkinViewUUID uuid={row?.uuid ?? ""} width={175} height={175} />
                </Box>
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose}>取消</Button>
                <Button onClick={handleOpen}>确认</Button>
            </DialogActions>
        </Dialog>
        {load && <Loading />}
        <Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'center' }} open={err != ""} >
            <Alert onClose={() => setErr("")} severity="error">{err}</Alert>
        </Snackbar>
    </>)
}