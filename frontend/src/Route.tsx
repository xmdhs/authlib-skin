import { Routes, Route, Outlet, createBrowserRouter, RouterProvider } from "react-router-dom";
import Login from '@/views/Login'
import Register from '@/views/Register'
import { ScrollRestoration } from "react-router-dom";

const router = createBrowserRouter([
    { path: "*", Component: Root },
]);

function Root() {
    return (
        <>
            <Routes>
                <Route path="/" element={<Layout />}></Route>
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
            <Outlet />
        </>
    )
}
