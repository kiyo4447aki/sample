import ctypes as C 
import uvc.purethermal_uvc_types as uvc
import queue
import numpy as np
import numpy.typing as npt
from uvc_frame import _py_frame_callback_factory
import time

BUF_SIZE = 8

class UvcSource:
    def __init__(self, frame_w: int, frame_h: int, fps: int) -> None:
        self.frame_w = frame_w
        self.frame_h = frame_h
        self.fps = fps
        
        self.ctx = C.POINTER(uvc.uvc_context)()
        self.dev = C.POINTER(uvc.uvc_device)()
        self.devh = C.POINTER(uvc.uvc_device_handle)()
        self.ctrl = uvc.uvc_stream_ctrl()
        
        self.frame_q: "queue.Queue[npt.NDArray[np.uint16]]" = queue.Queue(BUF_SIZE)
        self._cb = _py_frame_callback_factory(self.frame_q)
        self._last_ffc_ts = 0.0
        
    def open(self):
        if uvc.libuvc.uvc_init(C.byref(self.ctx), 0) < 0:
            raise RuntimeError("uvc init failed")
        
        if uvc.libuvc.uvc_find_device(self.ctx, C.byref(self.dev), uvc.PT_USB_VID, uvc.PT_USB_PID, 0) < 0:
            raise RuntimeError("uvc_find_device failed")
        
        if uvc.libuvc.uvc_open(self.dev, C.byref(self.devh)) < 0:
            raise RuntimeError("uvc_open failed")
        
        if uvc.libuvc.uvc_get_stream_ctrl_format_size(
            self.devh, C.byref(self.ctrl),
            uvc.UVC_FRAME_FORMAT_Y16, self.frame_w, self.frame_h, self.fps
        ) < 0:
            raise RuntimeError("uvc_get_stream_ctrl_format_size failed")
        
    def start(self):
        if uvc.libuvc.uvc_start_streaming(self.devh, C.byref(self.ctrl), self._cb, None, 0) < 0:
            raise RuntimeError("uvc_start_streaming failed")
    
    def read(self, timeout: float = 2.0) -> npt.NDArray[np.uint16]:
        try:
            return self.frame_q.get(timeout=timeout)
        except queue.Empty:
            raise RuntimeError("no frames from UVC")
        
    def maybe_ffc(self, interval_sec: int):
        if interval_sec <= 0:
            return
        now = time.time()
        if(now - self._last_ffc_ts) >= interval_sec:
            self._try_ffc()
            self._last_ffc_ts = now
    
    def _try_ffc(self):
        #TODO 手動ffcの実装
        pass
    
    def stop(self):
        try: 
            uvc.libuvc.uvc_stop_streaming(self.devh)
        except Exception:
            pass
    
    def close(self):
        try:
            uvc.libuvc.uvc_unref_device(self.dev)
            uvc.libuvc.uvc_exit(self.ctx)
        except Exception:
            pass


