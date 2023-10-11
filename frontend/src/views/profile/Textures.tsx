import { useEffect, useState } from "react";
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import FormControl from "@mui/material/FormControl";
import FormLabel from "@mui/material/FormLabel";
import FormControlLabel from "@mui/material/FormControlLabel";
import Radio from "@mui/material/Radio";
import RadioGroup from "@mui/material/RadioGroup";
import Button from "@mui/material/Button";
import { CardHeader } from "@mui/material";
import useTitle from "@/hooks/useTitle";
import { MuiFileInput } from 'mui-file-input'
import Box from "@mui/material/Box";
import ReactSkinview3d from '@/components/Skinview3d'
import { useUnmount } from "ahooks";
import { useAtomValue, useSetAtom } from "jotai";
import { LayoutAlertErr, token } from "@/store/store";
import { upTextures } from "@/apis/apis";
import Loading from "@/components/Loading";
import Snackbar from "@mui/material/Snackbar";

const Textures = function Textures() {
    const [redioValue, setRedioValue] = useState("skin")
    useTitle("上传皮肤")
    const [file, setFile] = useState<File | null>(null)
    const setErr = useSetAtom(LayoutAlertErr)
    const [loading, setLoading] = useState(false)
    const nowToken = useAtomValue(token)
    const [ok, setOk] = useState(false)
    const [skinInfo, setSkinInfo] = useState({
        skin: "",
        cape: "",
        model: "default"
    })

    useUnmount(() => {
        skinInfo.skin && URL.revokeObjectURL(skinInfo.skin)
        skinInfo.cape && URL.revokeObjectURL(skinInfo.cape)
    })

    useEffect(() => {
        if (file) {
            setSkinInfo(v => {
                URL.revokeObjectURL(v.skin);
                URL.revokeObjectURL(v.cape);
                return { skin: "", cape: "", model: "" }
            })
            const nu = URL.createObjectURL(file)
            switch (redioValue) {
                case "skin":
                    setSkinInfo({ skin: nu, cape: "", model: "default" })
                    break
                case "slim":
                    setSkinInfo({ skin: nu, cape: "", model: "slim" })
                    break
                case "cape":
                    setSkinInfo({ skin: "", cape: nu, model: "slim" })
            }
        }
    }, [file, redioValue])


    const onRadioChange = (_a: React.ChangeEvent<HTMLInputElement>, value: string) => {
        setRedioValue(value)
    }
    const handleChange = (newFile: File | null) => {
        setFile(newFile)
    }

    const handleToUpload = () => {
        if (!file || loading) return
        setLoading(true)
        const textureType = redioValue == "cape" ? "cape" : "skin"
        const model = redioValue == "slim" ? "slim" : ""
        upTextures(nowToken, textureType, model, file).then(() => setOk(true)).catch(e => [setErr(String(e)), console.warn(e)]).
            finally(() => setLoading(false))
    }



    return (<>
        <Box sx={{
            display: "grid", gap: "1em", gridTemplateAreas: {
                lg: '"a b" ". b"',
                xs: '"a" "b"'
            }, gridTemplateColumns: { lg: "1fr 1fr" }
        }}>
            <Card sx={{ gridArea: "a" }}>
                <CardHeader title="设置皮肤" />
                <CardContent>
                    <FormControl>
                        <FormLabel>类型</FormLabel>
                        <RadioGroup
                            row
                            onChange={onRadioChange}
                            value={redioValue}
                        >
                            <FormControlLabel value="skin" control={<Radio />} label="Steve" />
                            <FormControlLabel value="slim" control={<Radio />} label="Alex" />
                            <FormControlLabel value="cape" control={<Radio />} label="披风" />
                        </RadioGroup>
                        <br />
                        <MuiFileInput label="选择文件" value={file} inputProps={{ accept: 'image/png' }} onChange={handleChange} />
                        <br />
                        <Button variant="contained" sx={{ maxWidth: "3em" }} onClick={handleToUpload}>上传</Button>
                    </FormControl>
                </CardContent>
            </Card>
            <Card sx={{ gridArea: "b" }}>
                <CardHeader title="预览" />
                <CardContent>
                    {file && <ReactSkinview3d
                        skinUrl={skinInfo.skin}
                        capeUrl={skinInfo.cape}
                        height="250"
                        width="250"
                        options={{ model: skinInfo.model as "default" | "slim" }}
                    />}
                </CardContent>
            </Card>
        </Box>
        <Snackbar
            open={ok}
            autoHideDuration={6000}
            anchorOrigin={{ vertical: "top", horizontal: "center" }}
            onClose={() => setOk(false)}
            message="成功"
        />
        {loading && <Loading />}
    </>)
}

export default Textures