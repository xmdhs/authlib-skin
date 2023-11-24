import { Routes, Route, createBrowserRouter, RouterProvider, Outlet } from "react-router-dom";
import { ScrollRestoration } from "react-router-dom";
import Login from '@/views/Login'
import Register from '@/views/Register'
import Profile from '@/views/profile/Profile'
import Textures from '@/views/profile/Textures'
import Security from '@/views/profile/Security'
import Layout from '@/views/Layout'
import UserAdmin from "@/views/admin/UserAdmin";
import NeedLogin from "@/components/NeedLogin";
import Index from "@/views/Index";
import SignUpEmail from "@/views/SignUpEmail";

const router = createBrowserRouter([
    { path: "*", Component: Root },
])

function Root() {
    return (
        <>
            <Routes>
                <Route path="/" element={<Layout />}>
                    <Route index element={<Index />} />
                    <Route path="/*" element={<p>404</p>} />
                    <Route path="/login" element={<Login />} />
                    <Route path="/register" element={<Register />} />
                    <Route path="/register_email" element={<SignUpEmail />} />

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

