import { YggProfile } from "@/apis/model";

export function decodeSkin(y: YggProfile) {
    if (y.properties.length == 0) {
        return ["", "", ""]
    }
    const p = y.properties.find(v => v.name == "textures")
    if (!p?.value || p?.value == "") {
        return ["", "", ""]
    }
    const textures = JSON.parse(atob(p.value))

    const skin = textures?.textures?.SKIN?.url as string ?? ""
    const cape = textures?.textures?.CAPE?.url as string ?? ""
    const model = textures?.textures?.SKIN?.metadata?.model as string ?? "default"
    return [skin, cape, model]
}