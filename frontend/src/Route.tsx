import { Routes, Route, Outlet, createBrowserRouter, RouterProvider } from "react-router-dom";
import { ScrollRestoration } from "react-router-dom";
import Login from '@/views/Login'
import Register from '@/views/Register'
import User from '@/views/User'

const router = createBrowserRouter([
    { path: "*", Component: Root },
]);

function Root() {
    return (
        <>
            <Routes>
                <Route path="/" element={<Layout />}>
                    <Route path="/user" element={<User />} />

                </Route>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
            </Routes>
            <ScrollRestoration />
        </>
    )
}


export function PageRoute() {
    return (
        <>
            <RouterProvider router={router} />;
        </>
    )
}


function Layout() {
    return (
        <>
            <p>我是头部</p>
            <Outlet />
        </>
    )
}
