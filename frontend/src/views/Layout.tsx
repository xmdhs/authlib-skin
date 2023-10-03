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
import AppBar from '@mui/material/AppBar';
import { Outlet } from 'react-router-dom';
import { AccountCircle } from '@mui/icons-material';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import { LayoutAlertErr, token, user } from '@/store/store';
import { atom, useAtom, useAtomValue, useSetAtom } from 'jotai';
import Button from '@mui/material/Button';
import { useNavigate } from "react-router-dom";
import { useRequest, useMemoizedFn } from 'ahooks';
import { serverInfo, userInfo } from '@/apis/apis'
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';
import { memo } from 'react';
import useMediaQuery from '@mui/material/useMediaQuery';
import { ApiErr } from '@/apis/error';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import PersonIcon from '@mui/icons-material/Person';
import SecurityIcon from '@mui/icons-material/Security';
import SettingsIcon from '@mui/icons-material/Settings';
import useTilg from 'tilg'

const drawerWidth = 240;
const DrawerOpen = atom(false)

const DrawerHeader = styled('div')(({ theme }) => ({
    display: 'flex',
    alignItems: 'center',
    padding: theme.spacing(0, 1),
    // necessary for content to be below app bar
    ...theme.mixins.toolbar,
    justifyContent: 'flex-end',
}));

interface ListItem {
    icon: JSX.Element
    title: string
    link: string
}

const Layout = memo(function Layout() {
    const theme = useTheme();
    const [err, setErr] = useAtom(LayoutAlertErr)

    useTilg()

    return (<>
        <Box sx={{ display: 'flex' }}>
            <AppBar position="fixed"
                sx={{
                    zIndex: { lg: theme.zIndex.drawer + 1 }
                }}
            >
                <MyToolbar />
            </AppBar>
            <MyDrawer />
            <Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'center' }} open={err != ""} onClose={() => setErr("")}  >
                <Alert onClose={() => setErr("")} severity="error">{err}</Alert>
            </Snackbar>
            <Box
                component="main"
                sx={{
                    flexGrow: 1, bgcolor: 'background.default', p: 3
                }}
            >
                <Toolbar />
                <Container maxWidth="lg">
                    <Outlet />
                </Container>
            </Box>
        </Box>
    </>)
})

const MyToolbar = memo(function MyToolbar() {
    const [nowUser, setNowUser] = useAtom(user)
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
    const navigate = useNavigate();
    const [, setToken] = useAtom(token)
    const setErr = useSetAtom(LayoutAlertErr)
    const setOpen = useSetAtom(DrawerOpen)

    const server = useRequest(serverInfo, {
        cacheKey: "/api/yggdrasil",
        staleTime: 60000,
        onError: e => {
            console.warn(e)
            setErr(String(e))
        }
    })


    const handleLogOut = useMemoizedFn(() => {
        setAnchorEl(null);
        setNowUser({ name: "", uuid: "" })
        setToken("")
        navigate("/login")
    })

    useTilg()

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
                        onClick={() => setOpen(true)}
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

const MyList = memo(function MyList(p: { list: ListItem[] }) {
    useTilg()

    return (
        <>
            <List>
                {p.list.map(item =>
                    <MyListItem {...item} key={item.title} />
                )}
            </List>
        </>
    )
})

const MyListItem = function MyListItem(p: ListItem) {
    const navigate = useNavigate();

    const handleClick = () => {
        navigate(p.link)
    }

    return (
        <ListItem disablePadding>
            <ListItemButton onClick={handleClick}>
                <ListItemIcon>
                    {p.icon}
                </ListItemIcon>
                <ListItemText primary={p.title} />
            </ListItemButton>
        </ListItem>
    )
}

const MyDrawer = function MyDrawer() {
    const nowToken = useAtomValue(token)
    const setErr = useSetAtom(LayoutAlertErr)
    const navigate = useNavigate();
    const theme = useTheme();
    const isLg = useMediaQuery(theme.breakpoints.up('lg'))
    const [open, setOpen] = useAtom(DrawerOpen)

    const userinfo = useRequest(() => userInfo(nowToken), {
        refreshDeps: [nowToken],
        cacheKey: "/api/v1/user" + nowToken,
        staleTime: 60000,
        onError: e => {
            if (e instanceof ApiErr && e.code == 5) {
                navigate("/login")
            }
            console.warn(e)
            setErr(String(e))
        },
    })

    const userDrawerList = React.useMemo(() => [
        {
            icon: <PersonIcon />,
            title: '个人信息',
            link: '/profile'
        },
        {
            icon: <SettingsIcon />,
            title: '皮肤设置',
            link: '/textures'
        },
        {
            icon: <SecurityIcon />,
            title: '安全设置',
            link: '/security'
        }
    ] as ListItem[], [])

    const adminDrawerList = React.useMemo(() => [
        {
            icon: <PersonIcon />,
            title: 'test'
        }
    ] as ListItem[], [])

    useTilg()

    return (<>
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
                <MyList list={userDrawerList} />
                {userinfo.data?.is_admin && (
                    <>
                        <Divider />
                        <MyList list={adminDrawerList} />
                    </>)}
            </Drawer>
        )}
    </>)
}

export default Layout