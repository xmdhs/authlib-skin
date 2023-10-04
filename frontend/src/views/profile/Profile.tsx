import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import CardHeader from '@mui/material/CardHeader';
import { useHover, useRequest, useUnmount } from 'ahooks';
import { LayoutAlertErr, user } from '@/store/store';
import { useAtomValue, useSetAtom } from 'jotai';
import { yggProfile } from '@/apis/apis';
import { useNavigate } from 'react-router-dom';
import Box from '@mui/material/Box';
import { useEffect, useRef, useState } from 'react';
import { decodeSkin } from '@/utils/skin';
import ReactSkinview3d from "react-skinview3d"
import type { ReactSkinview3dOptions } from "react-skinview3d"
import { WalkingAnimation } from "skinview3d"
import type { SkinViewer } from "skinview3d"
import Skeleton from '@mui/material/Skeleton';
import useTitle from '@/hooks/useTitle';

const Profile = function Profile() {
    const navigate = useNavigate();
    const setErr = useSetAtom(LayoutAlertErr)
    const userinfo = useAtomValue(user)

    const [textures, setTextures] = useState({ skin: "", cape: "", model: "default" })
    useTitle("个人信息")


    const SkinInfo = useRequest(() => yggProfile(userinfo.uuid ?? ""), {
        cacheKey: "/api/yggdrasil/sessionserver/session/minecraft/profile/" + userinfo?.uuid,
        onError: e => {
            console.warn(e)
            setErr(String(e))
        },
        refreshDeps: [userinfo?.uuid],
    })

    useEffect(() => {
        if (!SkinInfo.data) return
        const [skin, cape, model] = decodeSkin(SkinInfo.data)
        setTextures({ cape: cape, skin: skin, model: model })
    }, [SkinInfo.data])


    

    return (
        <>
            <Box sx={{
                display: "grid", gap: "1em", gridTemplateAreas: {
                    lg: '"a b d" "c b d"',
                    xs: '"a" "b" "c" "d"'
                }, gridTemplateColumns: { lg: "1fr 1fr auto" }
            }}>
                <Card sx={{ gridArea: "a" }}>
                    <CardHeader title="信息" />
                    <CardContent sx={{ display: "grid", gridTemplateColumns: "4em auto" }}>
                        <Typography>name</Typography>
                        <Typography>{userinfo.name}</Typography>
                        <Typography>uuid</Typography>
                        <Typography sx={{ wordBreak: 'break-all' }}>{userinfo.uuid}</Typography>
                    </CardContent>
                    {/* <CardActions>
                    <Button size="small">更改</Button>
                </CardActions> */}
                </Card>
                <Card sx={{ gridArea: "b" }}>
                    <CardHeader title="皮肤" />
                    <CardContent sx={{ display: "flex", justifyContent: 'center' }}>
                        {
                            (SkinInfo.loading && !SkinInfo.data) ? <Skeleton variant="rectangular" width={250} height={250} />
                                : (textures.skin != "" || textures.cape != "") && (
                                    <MySkin
                                        skinUrl={textures.skin}
                                        capeUrl={textures.cape}
                                        height="250"
                                        width="250"
                                        options={{ model: textures.model as "default" | "slim" }}
                                    />)
                        }
                    </CardContent>
                    <CardActions>
                        <Button onClick={() => navigate('/textures')} size="small">更改</Button>
                    </CardActions>
                </Card>
                <Card sx={{ gridArea: "c" }}>
                    <CardHeader title="启动器设置" />
                    <CardContent>
                        <Typography>本站 Yggdrasil API 地址</Typography>
                        <code>{getYggRoot()}</code>
                    </CardContent>
                </Card>
                <Box sx={{ gridArea: "d" }}></Box>
            </Box >
        </>
    )
}


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

    return <div ref={refSkinview3d}>
        <ReactSkinview3d
            {...p}
            onReady={v => [v.viewer.animation = new WalkingAnimation(), v.viewer.autoRotate = true, skinview3dView.current = v.viewer]}
        />
    </div>
}

function getYggRoot() {
    const u = new URL((import.meta.env.VITE_APIADDR ?? location.origin) + "/api/yggdrasil")
    return u.toString()
}

export default Profile