import { Routes, Route, Outlet } from "react-router-dom";
import Login from '@/views/Login'

export function PageRoute() {
    return (
        <>
            <Routes>
                <Route path="/" element={<Layout />}></Route>
                <Route path="/login" element={<Login />} />
            </Routes>
        </>
    )
}


function Layout() {
    return (
        <>
            <Outlet />
        </>
    )
}
