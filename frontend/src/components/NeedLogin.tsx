import { userInfo } from "@/apis/apis";
import { ApiErr } from "@/apis/error";
import { token } from "@/store/store";
import { useRequest } from "ahooks";
import { useAtomValue } from "jotai";
import { useNavigate, Navigate } from "react-router-dom";


export default function NeedLogin({ children, needAdmin = false }: { children: JSX.Element, needAdmin?: boolean }) {
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