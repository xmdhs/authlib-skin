import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import useTitle from '@/hooks/useTitle';
import { useMemoizedFn, useRequest } from 'ahooks';
import { ListUser } from '@/apis/apis';
import { useEffect, useState } from 'react';
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
import { UserInfo } from '@/apis/model';
import { produce } from 'immer'

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

        <MyDialog open={open} setOpen={setOpen} row={row} />
    </>);
}

interface MyDialogProp {
    open: boolean
    setOpen: (b: boolean) => void
    row: UserInfo | null
}

function MyDialog({ open, row, setOpen }: MyDialogProp) {
    const handleClose = useMemoizedFn(() => {
        setOpen(false)
    })
    const [nrow, setNrow] = useState(row)

    useEffect(() => {
        setNrow(row)
    }, [row])


    return (
        <Dialog open={open}>
            <DialogTitle>修改用户信息</DialogTitle>
            <DialogContent>
                <TextField
                    margin="dense"
                    label="邮箱"
                    type="email"
                    fullWidth
                    variant="standard"
                    value={nrow?.email}
                    onChange={e => setNrow(produce(v => {
                        if (!v) return
                        v.email = e.target.value
                        return
                    }))}
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose}>取消</Button>
                <Button onClick={handleClose}>确认</Button>
            </DialogActions>
        </Dialog>
    )
}