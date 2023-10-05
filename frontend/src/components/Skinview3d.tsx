import { memo, useEffect, useRef } from "react";
import { SkinViewer, SkinViewerOptions } from "skinview3d";

// https://github.com/Hacksore/react-skinview3d/blob/master/src/index.tsx

/**
 * This is the interface that describes the parameter in `onReady`
 */
export interface ViewerReadyCallbackOptions {
    /**
     * The instance of the skinview3d
     */
    viewer: SkinViewer;
    /**
     * The ref to the canvas element
     */
    canvasRef: HTMLCanvasElement;
}

export interface ReactSkinview3dOptions {
    /**
     * The class names to apply to the canvas
     */
    className?: string;
    /**
     * The width of the canvas
     */
    width: number | string;
    /**
     * The height of the canvas
     */
    height: number | string;
    /**
     * The skin to load in the canvas
     */
    skinUrl: string;
    /**
     * The cape to load in the canvas
     */
    capeUrl?: string;
    /**
     * A function that is called when the skin viewer is ready
     * @param {SkinViewer} instance callback function to execute when the viewer is loaded {@link SkinViewer}
     * @example
     * onReady((instance) => {
     *  console.log(instance)
     * })
     */
    onReady?: ({ viewer, canvasRef }: ViewerReadyCallbackOptions) => void;
    /**
     * Parameters passed to the skinview3d constructor allowing you to override or add extra features
     * @notes please take a look at the upstream repo for more info
     * [bs-community/skinview3d](https://bs-community.github.io/skinview3d/)
     */
    options?: SkinViewerOptions;
}

/**
 * A skinview3d component
 */
const ReactSkinview3d = memo(function ReactSkinview3d({
    className,
    width,
    height,
    skinUrl,
    capeUrl,
    onReady,
    options,
}: ReactSkinview3dOptions) {
    const canvasRef = useRef<HTMLCanvasElement | null>(null);
    const skinviewRef = useRef<SkinViewer>();

    useEffect(() => {
        if (!canvasRef.current) return

        const viewer = new SkinViewer({
            canvas: canvasRef.current,
            width: Number(width),
            height: Number(height),
            ...options,
        });

        // handle cape/skin load initially
        skinUrl && viewer.loadSkin(skinUrl, { model: options?.model ?? "auto-detect" });
        capeUrl && viewer.loadCape(capeUrl);

        skinviewRef.current = viewer;

        // call onReady with the viewer instance
        if (onReady) {
            onReady({ viewer: skinviewRef.current, canvasRef: canvasRef.current });
        }

        return () => viewer.dispose()
    }, [capeUrl, height, onReady, options, skinUrl, width]);


    return <canvas className={className} ref={canvasRef} />;
})

export default ReactSkinview3d;