import Profile from "@/views/profile/Profile"
import Login from "@/views/Login"

export default function Index() {
    if (localStorage.getItem("token") && localStorage.getItem("token") != '""') {
        return <Profile />
    }
    return <Login />
}