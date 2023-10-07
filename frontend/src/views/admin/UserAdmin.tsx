import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import useTitle from '@/hooks/useTitle';
import { useRequest } from 'ahooks';
import { ListUser } from '@/apis/apis';
import { useEffect, useState } from 'react';
import { useAtomValue } from 'jotai';
import { token } from '@/store/store';
import TablePagination from '@mui/material/TablePagination';
import Alert from '@mui/material/Alert';
import Snackbar from '@mui/material/Snackbar';
import Button from '@mui/material/Button';


export default function UserAdmin() {
    useTitle("用户管理")
    const [page, setPage] = useState(1)
    const nowtoken = useAtomValue(token)
    const [err, setErr] = useState("")

    const { data, run } = useRequest(ListUser, {
        cacheKey: "/api/v1/admin/users?page=" + String(page),
        onError: e => {
            setErr(String(e))
        }
    })

    useEffect(() => {
        run(page, nowtoken)
    }, [page, nowtoken, run])


    return (<>
        <Paper>
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
                                <TableCell><Button>编辑</Button></TableCell>
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

    </>);
}