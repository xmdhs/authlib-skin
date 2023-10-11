import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import CardHeader from '@mui/material/CardHeader';
import { user } from '@/store/store';
import { useAtomValue } from 'jotai';
import { useNavigate } from 'react-router-dom';
import Box from '@mui/material/Box';
import useTitle from '@/hooks/useTitle';
import SkinViewUUID from '@/components/SkinViewUUID';
import root from '@/utils/root';

const Profile = function Profile() {
    const navigate = useNavigate();
    const userinfo = useAtomValue(user)

    useTitle("个人信息")


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
                        <SkinViewUUID uuid={userinfo?.uuid ?? ""} width={250} height={250} />
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

function getYggRoot() {
    const u = new URL(root() + "/api/yggdrasil")
    return u.toString()
}

export default Profile