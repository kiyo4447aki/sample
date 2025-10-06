from uvc import UvcSource
from detector import Detector
from renderer import Renderer
from uploader import Uploader
from spooler import Spooler

from utils import http_client, thermal

import logging
import uuid
from requests.exceptions import HTTPError
import time
from datetime import datetime, timezone

try:
    from config import LogConfig, UvcConfig, SpoolConfig, RenderConfig, BackendConfig, DetectorConfig
except Exception as e:
    print(f"failed to load config: {e}")
    import sys
    sys.exit(1)

class Controller:
    logConfig = LogConfig()
    uvcConfig = UvcConfig()
    spoolConfig = SpoolConfig()
    renderConfig = RenderConfig()
    backendConfig = BackendConfig()
    detectorConfig = DetectorConfig()
    
    def __init__(self):
        self._logger = logging.getLogger("thermalwatcher")
        self.source = UvcSource(
            self.uvcConfig.width, 
            self.uvcConfig.height, 
            self.uvcConfig.fps
        )
        self.detector = Detector(
            self.detectorConfig.abs_thresh, 
            self.detectorConfig.delta_thresh,
            self.detectorConfig.min_blob_pix,
            self.detectorConfig.persist_frames,
            self.detectorConfig.roi_rect
        )
        self.renderer = Renderer(
            self.renderConfig.cmap_min,
            self.renderConfig.cmap_max,
            self.renderConfig.jpeg_quality,
        )
        self.uploader = Uploader(
            http_client.get_session(),
            self.backendConfig.auth_api_url,
            self.backendConfig.notification_api_url,
            self.backendConfig.auth_id,
            self.backendConfig.auth_password,
            self.backendConfig.uploads_api_url
        )
        self.spooler = Spooler(
            self.spoolConfig.spool_dir,
            self.spoolConfig.max_spool_bytes,
            self.uploader,
            self.spoolConfig.reconnect_max
        )
        
        self._cooldown_until = 0.0
        self._backoff = 1.0
        self._opened = False
        
    #TODO アラート詳細をデータクラスとして定義、引数の受け渡し方法を変更
    def _send_or_spool(self, jpeg: bytes, severity: str, temp_c: float) -> bool:
        event_id = str(uuid.uuid4())
        occurred_at = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        mime_type = "image/jpeg"
        
        try:
            if jpeg is None:
                raise RuntimeError("no jpeg input")
            #TODO widthとheightをroiに変更
            self.uploader.send_all(
                jpeg, occurred_at, event_id, mime_type, severity, temp_c, None, 
                self.uvcConfig.width, self.uvcConfig.height
            )
            self._logger.info("alert sent: sev=%s", severity)
            return True
        except HTTPError as he:
            self._logger.warning("send http error -> spool: %s", he)
            if jpeg is not None:
                self.spooler.enqueue(
                    event_id, jpeg, occurred_at, severity, temp_c,
                    self.uvcConfig.width, self.uvcConfig.height, None, mime_type,
                )
            return False
        except Exception as e:
            self._logger.warning("send error -> spool: %s", e)
            if jpeg is not None:
                self.spooler.enqueue(
                    event_id, jpeg, occurred_at, severity, temp_c,
                    self.uvcConfig.width, self.uvcConfig.height, None, mime_type,
                )
            return False
    
    def run_forever(self):
        self.spooler.start()
        while True:
            try:
                if not getattr(self, "_opened", False):
                    self.source.open()
                    self.source.start()
                    self._logger.info("streaming started (Y16 %dx%d@%dfps)", self.uvcConfig.width, self.uvcConfig.height, self.uvcConfig.fps)
                    self._backoff = 1.0
                    self._opened = True
                    
                raw = self.source.read(timeout=2.0)
                celsius = thermal.ktoc(raw)
                
                hot, maxC, bbox, _, amb = self.detector.infer(celsius)
                
                now = time.monotonic()
                if now < self._cooldown_until:
                    continue
                
                if hot:
                    severity = thermal.choose_severity(maxC, self.detectorConfig.abs_thresh)
                    
                    try:
                        jpeg = self.renderer.render_jpeg(celsius, bbox, maxC, amb)
                    except Exception as e:
                        self._logger.exception("failed to generate jpeg: %s", e)
                        jpeg = None
                    if jpeg is not None:
                        ok = self._send_or_spool(jpeg, severity, maxC)
                        self._cooldown_until = now + (self.detectorConfig.cooldown_sec if ok else min(self.detectorConfig.cooldown_sec, 10))
                        
                self._backoff = 1.0
                
            except KeyboardInterrupt:
                self._logger.info("KeyboardInterrupt")
                break
            except Exception as e:
                self._logger.warning("stream error: %s", e)
                try:
                    self.source.stop()
                    self.source.close()
                except Exception as e:
                    self._logger.warning("failed to stop uvc stream: %s", e)
                self._opened = False
                time.sleep(self._backoff)
                self._backoff = min(self.spoolConfig.reconnect_max, self._backoff * 2.0)
        try:
            self.source.stop()
            self.source.close()
        except Exception as e:
            self._logger.warning("failed to stop uvc stream: %s", e)
        self._opened = False
        self.spooler.stop()
