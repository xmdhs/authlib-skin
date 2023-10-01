import * as React from 'react';
import { styled, useTheme } from '@mui/material/styles';
import Box from '@mui/material/Box';
import Drawer from '@mui/material/Drawer';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import InboxIcon from '@mui/icons-material/MoveToInbox';
import MailIcon from '@mui/icons-material/Mail';
import AppBar from '@mui/material/AppBar';
import { Outlet } from 'react-router-dom';
import { AccountCircle } from '@mui/icons-material';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import { token, user } from '@/store/store';
import { useAtom, useAtomValue } from 'jotai';
import Button from '@mui/material/Button';
import { useNavigate } from "react-router-dom";
import { useRequest, useMemoizedFn } from 'ahooks';
import { serverInfo, userInfo } from '@/apis/apis'
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';
import { memo } from 'react';
import useMediaQuery from '@mui/material/useMediaQuery';
import useTilg from 'tilg'
import { ApiErr } from '@/apis/error';
import Typography from '@mui/material/Typography';


const drawerWidth = 240;

const DrawerHeader = styled('div')(({ theme }) => ({
    display: 'flex',
    alignItems: 'center',
    padding: theme.spacing(0, 1),
    // necessary for content to be below app bar
    ...theme.mixins.toolbar,
    justifyContent: 'flex-end',
}));


export default function Layout() {
    const theme = useTheme();
    const isLg = useMediaQuery(theme.breakpoints.up('lg'))
    const [open, setOpen] = React.useState(false);
    const nowToken = useAtomValue(token)
    const [err, setErr] = React.useState("");
    const navigate = useNavigate();


    const userinfo = useRequest(() => userInfo(nowToken), {
        refreshDeps: [nowToken],
        cacheKey: "/api/v1/user",
        cacheTime: 10000,
        onError: e => {
            if (e instanceof ApiErr && e.code == 5) {
                navigate("/login")
            }
            console.warn(e)
            setErr(String(e))
        }
    })

    useTilg(isLg, open)

    return (<>
        <Box sx={{ display: 'flex' }}>
            <AppBar position="fixed"
                sx={{
                    zIndex: { lg: theme.zIndex.drawer + 1 }
                }}
            >
                <MyToolbar setOpen={setOpen}></MyToolbar>
            </AppBar>
            {userinfo.data && (
                <Drawer
                    sx={{
                        width: drawerWidth,
                        flexShrink: 0,
                        '& .MuiDrawer-paper': {
                            width: drawerWidth,
                            boxSizing: 'border-box',
                        },
                    }}
                    variant={isLg ? "persistent" : "temporary"}
                    anchor="left"
                    open={open || isLg}
                    onClose={() => setOpen(false)}
                >
                    <DrawerHeader>
                        <IconButton onClick={() => setOpen(false)}>
                            {theme.direction === 'ltr' ? <ChevronLeftIcon /> : <ChevronRightIcon />}
                        </IconButton>
                    </DrawerHeader>
                    <Divider />
                    <List>
                        {['Inbox', 'Starred', 'Send email', 'Drafts'].map((text, index) => (
                            <ListItem key={text} disablePadding>
                                <ListItemButton>
                                    <ListItemIcon>
                                        {index % 2 === 0 ? <InboxIcon /> : <MailIcon />}
                                    </ListItemIcon>
                                    <ListItemText primary={text} />
                                </ListItemButton>
                            </ListItem>
                        ))}
                    </List>
                    {userinfo.data?.is_admin && (
                        <>
                            <Divider />
                            <List>
                                {['All mail', 'Trash', 'Spam'].map((text, index) => (
                                    <ListItem key={text} disablePadding>
                                        <ListItemButton>
                                            <ListItemIcon>
                                                {index % 2 === 0 ? <InboxIcon /> : <MailIcon />}
                                            </ListItemIcon>
                                            <ListItemText primary={text} />
                                        </ListItemButton>
                                    </ListItem>
                                ))}
                            </List>
                        </>)}
                </Drawer>
            )}
            <Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'center' }} open={err != ""} onClose={() => setErr("")}  >
                <Alert onClose={() => setErr("")} severity="error">{err}</Alert>
            </Snackbar>
            <Box
                component="main"
                sx={{
                    flexGrow: 1, bgcolor: 'background.default', p: 3
                }}
            >
                <Outlet />
            </Box>
        </Box>
    </>)
}

const MyToolbar = memo((p: { setOpen: (v: boolean) => void }) => {
    const [nowUser, setNowUser] = useAtom(user)
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
    const navigate = useNavigate();
    const [, setToken] = useAtom(token)

    const server = useRequest(serverInfo, {
        cacheKey: "/api/yggdrasil",
        cacheTime: 100000,
        onError: e => {
            console.warn(e)
        }
    })



    const handleLogOut = useMemoizedFn(() => {
        setAnchorEl(null);
        setNowUser({ name: "", uuid: "" })
        setToken("")
        navigate("/")
    })



    return (
        <>
            <Toolbar>
                {nowUser.name != "" && (<>
                    <IconButton
                        size="large"
                        edge="start"
                        color="inherit"
                        aria-label="menu"
                        sx={{ mr: 2, display: { lg: 'none' } }}
                        onClick={() => p.setOpen(true)}
                    >
                        <MenuIcon />
                    </IconButton >
                </>)
                }
                <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                    {server.data?.meta.serverName ?? "皮肤站"}
                </Typography>
                {nowUser.name != "" && (
                    <div>
                        <IconButton
                            size="large"
                            aria-label="account of current user"
                            aria-controls="menu-appbar"
                            aria-haspopup="true"
                            onClick={event => setAnchorEl(event.currentTarget)}
                            color="inherit"
                        >
                            <AccountCircle />
                        </IconButton>
                        <Menu
                            id="menu-appbar"
                            anchorEl={anchorEl}
                            anchorOrigin={{
                                vertical: 'top',
                                horizontal: 'right',
                            }}
                            keepMounted
                            transformOrigin={{
                                vertical: 'top',
                                horizontal: 'right',
                            }}
                            open={Boolean(anchorEl)}
                            onClose={() => setAnchorEl(null)}
                        >
                            <MenuItem onClick={handleLogOut}>登出</MenuItem>
                        </Menu>
                    </div>
                )}
                {nowUser.name == "" && (
                    <Button color="inherit" onClick={() => navigate("/login")} >登录</Button>
                )}
            </Toolbar>
        </>)
})
