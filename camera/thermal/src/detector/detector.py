import numpy as np
import logging
try:
    import cv2
except Exception:
    
    cv2 = None
from typing import TypeAlias

RoiRect: TypeAlias = tuple[int, int, int, int]
DetectResult: TypeAlias = tuple[bool, float, RoiRect | None, int, float]


class Detector:
    def __init__(self, abs_th: float, delta_th: float, min_blob: int, persist_frames: int, roi_rect: str):
        self._logger = logging.getLogger("thermalwatcher")
        self.abs_th = abs_th    #絶対温度との閾値
        self.delta_th = delta_th       #周辺温度との差分閾値
        self.min_blob = min_blob       #異常発熱と判断する際の最低ピクセル数
        self.persist_frames = persist_frames      #以上発熱と判断する際のフレーム数
        self.roi_rect = roi_rect        #"x,y,w,h"形式でroi、空文字なら全面
        self._ok_count = 0
        self._kernel = np.ones((3,3), np.uint8) if cv2 is not None else None

        
    def _roi(self, img: np.ndarray) -> RoiRect:
        #roi_rectがある場合、"x,y,w,h"文字列をパース
        if not self.roi_rect:
            return (0, 0, img.shape[1], img.shape[0])
        parts = [p.strip() for p in self.roi_rect.split(",")]
        if len(parts) != 4:
            self._logger.warning(
                "Invalid ROI_RECT %r (value must be x,y,w,h format). Use full frame.", self.roi_rect, 
            )
            return (0, 0, img.shape[1], img.shape[0])
        try:
            x, y, w, h = map(int, parts)
        except Exception as e:
            self._logger.warning(
                "Invalid ROI_RECT %r (%s). Use full frame.", self.roi_rect, e
            )
            return (0, 0, img.shape[1], img.shape[0])
        x = max(0, min(x, img.shape[1] - 1))
        y = max(0, min(y, img.shape[0] - 1))
        w = max(1, min(w, img.shape[1] - x))
        h = max(1, min(h, img.shape[0] - y))
        
        return (x, y, w, h)
    
    def infer(self, celsius: np.ndarray) -> DetectResult:
        """
        returns: 
            (
                hot_for_alert: bool,   発熱の検知結果
                maxC: float,        ROI内の最高温度
                bbox: tuple[int, int, int, int] | None,     # ホット領域の外接矩形 (x, y, w, h)
                area: int,      # ホット領域の面積（ピクセル数）
                ambient: float      # ROI内の周囲温度（中央値）
            )
        """
        x, y, w, h = self._roi(celsius)
        roi = celsius[y:y+h, x:x+w]
        ambient = float(np.nanmedian(roi))
        hotmask = (roi >= self.abs_th) | ((roi - ambient) >= self.delta_th)
        
        if (cv2 is not None) and (self._kernel is not None):
            mask = (hotmask.astype(np.uint8) * 255)
            mask = cv2.morphologyEx(mask, cv2.MORPH_OPEN, self._kernel, iterations=1)
            found = cv2.findContours(mask, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
            cnts = found[0] if len(found) == 2 else found[1]
            if not cnts:
                hot = False; area = 0; bbox = None
            else:
                c = max(cnts, key=cv2.contourArea)
                area = int(cv2.contourArea(c))
                if area < self.min_blob:
                    hot = False; bbox = None
                else:
                    bx, by, bw, bh = cv2.boundingRect(c)
                    hot = True; bbox = (x+bx, y+by, bw, bh)
        else:
            area = int(np.sum(hotmask))
            hot = area >= self.min_blob
            bbox = None
            
        maxC = float(np.nanmax(roi))
        self._ok_count = (self._ok_count + 1) if hot else 0
        hot_for_alert = hot and (self._ok_count >= self.persist_frames)
        return hot_for_alert, maxC, bbox, area, ambient