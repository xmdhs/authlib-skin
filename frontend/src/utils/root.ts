export default function root() {
    if (import.meta.env.VITE_APIADDR != "") {
        return import.meta.env.VITE_APIADDR
    }
    return location.origin
}