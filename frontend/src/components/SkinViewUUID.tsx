import { yggProfile } from "@/apis/apis";
import { decodeSkin } from "@/utils/skin";
import Skeleton from "@mui/material/Skeleton";
import { useHover, useMemoizedFn, useRequest, useUnmount } from "ahooks";
import { memo, useEffect, useRef, useState } from "react";
import ReactSkinview3d, { ReactSkinview3dOptions } from "@/components/Skinview3d";
import { SkinViewer, WalkingAnimation } from "skinview3d";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";

interface prop {
    uuid: string
    width: number
    height: number
}

const SkinViewUUID = memo(function SkinViewUUID({ uuid, width, height }: prop) {
    const [textures, setTextures] = useState({ skin: "", cape: "", model: "default" })
    const [err, setErr] = useState("")

    const SkinInfo = useRequest(() => yggProfile(uuid), {
        cacheKey: "/api/yggdrasil/sessionserver/session/minecraft/profile/" + uuid,
        onError: e => {
            console.warn(e)
            setErr(String(e))
        },
        refreshDeps: [uuid],
    })

    useEffect(() => {
        if (!SkinInfo.data) return
        const [skin, cape, model] = decodeSkin(SkinInfo.data)
        setTextures({ cape: cape, skin: skin, model: model })
    }, [SkinInfo.data])

    if (err != "") {
        return <Typography color={"error"}>{err}</Typography>
    }
    return (<>
        {
            (SkinInfo.loading && !SkinInfo.data) ? <Skeleton variant="rectangular" width={width} height={height} />
                : (textures.skin != "" || textures.cape != "") ? (
                    <MySkin
                        skinUrl={textures.skin}
                        capeUrl={textures.cape}
                        height={width}
                        width={height}
                        options={{ model: textures.model as "default" | "slim" }}
                    />) : <Box sx={{ minHeight: height + "px" }}>
                    <Typography>还没有设置皮肤</Typography>
                </Box>
        }
    </>)
})


const MySkin = function MySkin(p: ReactSkinview3dOptions) {
    const refSkinview3d = useRef(null);
    const skinisHovering = useHover(refSkinview3d);
    const skinview3dView = useRef<SkinViewer | null>(null);

    useEffect(() => {
        if (skinview3dView.current) {
            skinview3dView.current.autoRotate = !skinisHovering
        }
        if (skinview3dView.current?.animation) {
            skinview3dView.current.animation.paused = skinisHovering
        }
    }, [skinisHovering])

    useUnmount(() => {
        skinview3dView.current?.dispose()
    })

    const handelOnReady = useMemoizedFn(v => {
        v.viewer.animation = new WalkingAnimation()
        v.viewer.autoRotate = true
        skinview3dView.current = v.viewer
    })

    return <div ref={refSkinview3d}>
        <ReactSkinview3d
            {...p}
            onReady={handelOnReady}
        />
    </div>
}

export default SkinViewUUID