import { useEffect, useRef, useState } from "react";
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import FormControl from "@mui/material/FormControl";
import FormLabel from "@mui/material/FormLabel";
import FormControlLabel from "@mui/material/FormControlLabel";
import Radio from "@mui/material/Radio";
import RadioGroup from "@mui/material/RadioGroup";
import Button from "@mui/material/Button";
import { CardHeader } from "@mui/material";
import useTilg from "tilg";
import useTitle from "@/hooks/useTitle";
import { MuiFileInput } from 'mui-file-input'
import Box from "@mui/material/Box";
import ReactSkinview3d from "react-skinview3d";
import { useUnmount } from "ahooks";
import { SkinViewer } from "skinview3d";

const Textures = function Textures() {
    const [redioValue, setRedioValue] = useState("skin")
    useTitle("上传皮肤")
    const [file, setFile] = useState<File | null>(null)
    const skin = useRef({
        skinUrl: "",
        capeUrl: "",
    })
    const skinview3dView = useRef<SkinViewer | null>(null);

    useUnmount(() => {
        skin.current.skinUrl && URL.revokeObjectURL(skin.current.skinUrl)
        skin.current.capeUrl && URL.revokeObjectURL(skin.current.capeUrl)
        skinview3dView.current?.dispose()
    })

    useEffect(() => {
        if (file) {
            const nu = URL.createObjectURL(file)
            skin.current.skinUrl && URL.revokeObjectURL(skin.current.skinUrl)
            skin.current.capeUrl && URL.revokeObjectURL(skin.current.capeUrl)
            skinview3dView.current?.loadSkin(null)
            skinview3dView.current?.loadCape(null)
            switch (redioValue) {
                case "skin":
                    skin.current.skinUrl = nu
                    skinview3dView.current?.loadSkin(nu, { model: "default" }).then(() =>
                        skinview3dView.current?.loadSkin(nu, { model: "default" })
                    )
                    break
                case "slim":
                    skin.current.skinUrl = nu
                    skinview3dView.current?.loadSkin(nu, { model: "slim" }).then(() =>
                        skinview3dView.current?.loadSkin(nu, { model: "slim" })
                    )
                    break
                case "cape":
                    skin.current.capeUrl = nu
                    skinview3dView.current?.loadCape(nu).then(() => {
                        skinview3dView.current?.loadCape(nu)
                    })
            }
        }
    }, [file, redioValue])


    const onRadioChange = (_a: React.ChangeEvent<HTMLInputElement>, value: string) => {
        setRedioValue(value)
    }
    const handleChange = (newFile: File | null) => {
        setFile(newFile)
    }

    useTilg()

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
                    </FormControl>
                </CardContent>
                <CardActions>
                    <Button variant="contained" sx={{ maxWidth: "3em" }}>上传</Button>
                </CardActions>
            </Card>
            <Card sx={{ gridArea: "b" }}>
                <CardHeader title="预览" />
                <CardContent>
                    {file && <ReactSkinview3d
                        skinUrl={""}
                        capeUrl={""}
                        height="250"
                        width="250"
                        onReady={v => skinview3dView.current = v.viewer}
                    />}
                </CardContent>
            </Card>
        </Box>
    </>)
}

export default Textures