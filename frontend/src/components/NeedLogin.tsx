import { userInfo } from "@/apis/apis";
import { ApiErr } from "@/apis/error";
import { token } from "@/store/store";
import { useRequest } from "ahooks";
import { useAtomValue } from "jotai";
import { useEffect } from "react";
import { useNavigate, Navigate } from "react-router-dom";


export default function NeedLogin({ children, needAdmin = false }: { children: JSX.Element, needAdmin?: boolean }) {
    const t = useAtomValue(token)
    const navigate = useNavigate();
    const u = useRequest(() => userInfo(t), {
        refreshDeps: [t],
        cacheKey: "/api/v1/user" + t,
        staleTime: 60000,
        onError: e => {
            if (e instanceof ApiErr && e.code == 5) {
                navigate("/login")
            }
            console.warn(e)
        }
    })

    useEffect(() => {
        if (!u.data) return
        if (!u.data.is_admin && needAdmin) {
            navigate("/login")
        }
        if (u.data.uuid == "") {
            navigate("/login")
        }
    }, [navigate, needAdmin, u.data])

    if (!localStorage.getItem("token") || localStorage.getItem("token") == '""') {
        return <Navigate to="/login" />
    }


    return <> {children}</>
}