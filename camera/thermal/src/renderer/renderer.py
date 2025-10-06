import numpy as np
import cv2
from datetime import datetime, timezone

class Renderer:
    def __init__ (self, vmin: float, vmax: float, jpeg_quality: int):
        self.vmin = vmin
        self.vmax = vmax
        if 0 <= jpeg_quality <= 100:
            self.jpeg_quality = jpeg_quality
        else:
            raise ValueError("jpeg_quality must be 0 to 100 (got："+ str(jpeg_quality) + ")" )
        
    def render_jpeg(self, celsius: np.ndarray, bbox: tuple[int, int, int, int] | None, maxC: float, ambient: float) -> bytes:
        if celsius.ndim != 2:
            raise ValueError(f"celsius must be 2D array, got shape={celsius.shape}")
        arr = celsius.astype(np.float32, copy=False)
        vmin, vmax = self.vmin, self.vmax
        vis = np.clip((arr - vmin) * (255.0 / max(1e-6, vmax - vmin)), 0, 255).astype(np.uint8)
        color = cv2.applyColorMap(vis, cv2.COLORMAP_INFERNO)
        now = datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%SZ")
        cv2.putText(color, f"max={maxC:.1f}C amb={ambient:.1f}C", (4, 14), 
                    cv2.FONT_HERSHEY_SIMPLEX, 0.4, (255,255,255), 1, cv2.LINE_AA)
        cv2.putText(color, f"time={now}", (4,30),
                    cv2.FONT_HERSHEY_SIMPLEX, 0.4, (255,255,255), 1, cv2.LINE_AA)
        if bbox is not None:
            x,y,w,h = bbox
            cv2.rectangle(color, (x,y), (x+w, y+h), (255,255,255), 1)
        color = np.ascontiguousarray(color)
        ok, buf = cv2.imencode(".jpg", color, [cv2.IMWRITE_JPEG_QUALITY, self.jpeg_quality])
        if not ok:
            raise RuntimeError("imencode failed")
        return buf.tobytes()
            