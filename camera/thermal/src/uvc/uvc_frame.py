import queue
import numpy as np
import numpy.typing as npt
import ctypes as C
import logging
import purethermal_uvc_types as uvc
from typing import Any

def _py_frame_callback_factory(
    frame_q: "queue.Queue[npt.NDArray[np.uint16]]",
):
    logger = logging.getLogger("thermalwatcher")

    # キャッシュ用（解像度が変わらない前提で一度だけ作る）
    cached: dict[str, Any] = {
        "w": None,          # type: Optional[int]
        "h": None,          # type: Optional[int]
        "ArrayType": None,  # type: Optional[type]
    }

    def _cb(frame_ptr: C.c_void_p, userptr: C.c_void_p) -> None:
        try:
            fptr = C.cast(frame_ptr, C.POINTER(uvc.uvc_frame))
            f = fptr.contents

            w = int(f.width); h = int(f.height)
            if int(f.data_bytes) != 2 * w * h:
                return

            # まず満杯なら即捨てる（後続の配列生成コストを回避）
            if frame_q.full():
                return

            # 配列型を一度だけ生成してキャッシュ
            if cached["w"] != w or cached["h"] != h or cached["ArrayType"] is None:
                cached["w"] = w
                cached["h"] = h
                cached["ArrayType"] = C.c_uint16 * (w * h)

            ArrayType = cached["ArrayType"]
            ap = C.cast(f.data, C.POINTER(ArrayType))
            # buffer 明示でゼロコピー
            data_view = np.frombuffer(memoryview(ap.contents), dtype=np.uint16)
            data_view = data_view.reshape(h, w)

            # バッファ寿命を独立させるためにコピーして投入（元実装と同じ）
            frame_q.put_nowait(data_view.copy())

        except Exception:
            logger.exception("frame callback error")

    return C.CFUNCTYPE(None, C.POINTER(uvc.uvc_frame), C.c_void_p)(_cb)
