import { Routes, Route, createBrowserRouter, RouterProvider, useNavigate } from "react-router-dom";
import { ScrollRestoration } from "react-router-dom";
import Login from '@/views/Login'
import Register from '@/views/Register'
import Profile from '@/views/profile/Profile'
import Textures from '@/views/profile/Textures'
import Security from '@/views/profile/Security'
import Layout from '@/views/Layout'
import { useAtomValue } from "jotai";
import { token } from "@/store/store";

const router = createBrowserRouter([
    { path: "*", Component: Root },
])

function Root() {
    return (
        <>
            <Routes>
                <Route path="/" element={<Layout />}>
                    <Route index element={<p>123</p>} />
                    <Route path="/*" element={<p>404</p>} />
                    <Route path="/login" element={<Login />} />
                    <Route path="/register" element={<Register />} />

                    <Route path="/profile" element={<NeedLogin><Profile /></NeedLogin>} />
                    <Route path="/textures" element={<NeedLogin><Textures /></NeedLogin>} />
                    <Route path="/security" element={<NeedLogin><Security /></NeedLogin>} />
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


function NeedLogin({ children }: { children: JSX.Element }) {
    const t = useAtomValue(token)
    const navigate = useNavigate();
    if (t == "") {
        navigate("/login")
        return
    }
    return <> {children}</>
}