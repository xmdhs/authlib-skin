import { Routes, Route, createBrowserRouter, RouterProvider, useNavigate, Outlet, Navigate } from "react-router-dom";
import { ScrollRestoration } from "react-router-dom";
import Login from '@/views/Login'
import Register from '@/views/Register'
import Profile from '@/views/profile/Profile'
import Textures from '@/views/profile/Textures'
import Security from '@/views/profile/Security'
import Layout from '@/views/Layout'
import { useAtomValue } from "jotai";
import { token } from "@/store/store";
import { ApiErr } from "@/apis/error";
import { userInfo } from "@/apis/apis";
import { useRequest } from "ahooks";
import UserAdmin from "@/views/admin/UserAdmin";

const router = createBrowserRouter([
    { path: "*", Component: Root },
])

function Root() {
    return (
        <>
            <Routes>
                <Route path="/" element={<Layout />}>
                    <Route index />
                    <Route path="/*" element={<p>404</p>} />
                    <Route path="/login" element={<Login />} />
                    <Route path="/register" element={<Register />} />

                    <Route element={<NeedLogin><Outlet /></NeedLogin>}>
                        <Route path="/profile" element={<Profile />} />
                        <Route path="/textures" element={<Textures />} />
                        <Route path="/security" element={<Security />} />
                    </Route>

                    <Route path="admin" element={<NeedLogin needAdmin><Outlet /></NeedLogin>}>
                        <Route path="user" element={<UserAdmin />} />
                    </Route>

                </Route>
            </Routes>
            <ScrollRestoration />
        </>
    )
}


export function PageRoute() {
    return (
        <>
            <RouterProvider router={router} />
        </>
    )
}


function NeedLogin({ children, needAdmin = false }: { children: JSX.Element, needAdmin?: boolean }) {
    const t = useAtomValue(token)
    const navigate = useNavigate();
    useRequest(() => userInfo(t), {
        refreshDeps: [t],
        cacheKey: "/api/v1/user" + t,
        staleTime: 60000,
        onError: e => {
            if (e instanceof ApiErr && e.code == 5) {
                navigate("/login")
            }
            console.warn(e)
        },
        onSuccess: u => {
            if (!u) return
            if (!u.is_admin && needAdmin) {
                navigate("/login")
            }
            if (u.uuid == "") {
                navigate("/login")
            }
        }
    })

    if (!localStorage.getItem("token") || localStorage.getItem("token") == '""') {
        return <Navigate to="/login" />
    }


    return <> {children}</>
}