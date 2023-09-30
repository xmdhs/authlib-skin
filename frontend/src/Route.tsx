import { Routes, Route, createBrowserRouter, RouterProvider } from "react-router-dom";
import { ScrollRestoration } from "react-router-dom";
import Login from '@/views/Login'
import Register from '@/views/Register'
import User from '@/views/User'
import Layout from '@/views/Layout'

const router = createBrowserRouter([
    { path: "*", Component: Root },
])

function Root() {
    return (
        <>
            <Routes>
                <Route path="/" element={<Layout />}>
                    <Route path="/login" element={<Login />} />
                    <Route path="/register" element={<Register />} />
                    <Route path="/profile" element={<User />}>
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


