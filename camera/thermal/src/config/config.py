from dataclasses import dataclass
import os

from dotenv import load_dotenv
load_dotenv()

@dataclass
class UvcConfig:
    width:         int = int(os.getenv("FRAME_W", "160"))
    height:        int = int(os.getenv("FRAME_H", 120))
    fps:           int = int(os.getenv("FPS", "9"))
    ffc_interval:  int = int(os.getenv("FFC_INTERVAL_SEC", 0))

@dataclass
class DetectorConfig:
    abs_thresh:     float = float(os.getenv("ABS_THRESH_C", "65"))
    delta_thresh:   float = float(os.getenv("DELTA_THRESH_C", "18"))
    min_blob_pix:   int   = int(os.getenv("MIN_BLOB_PIX", "40"))
    persist_frames: int   = int(os.getenv("PERSIST_FRAMES", "3"))
    cooldown_sec:   int   = int(os.getenv("COOLDOWN_SEC", "30"))
    roi_rect:       str   = os.getenv("ROI_RECT", "") 

@dataclass
class RenderConfig:
    cmap_min:     float = float(os.getenv("CMAP_MIN_C", "0"))
    cmap_max:     float = float(os.getenv("CMAP_MAX_C", "120"))
    jpeg_quality: int   = int(os.getenv("JPEG_QUALITY", "90"))

@dataclass
class BackendConfig:
    device_id:            str = os.getenv("DEVICE_ID", "unknown")
    uploads_api_url:      str = os.getenv("UPLOADS_API_URL", "")
    notification_api_url: str = os.getenv("ALERTS_API_URL", "https://api.example.com/alerts")
    auth_api_url:         str = os.getenv("AUTH_API_URL", "https://api.example.com/alerts")
    auth_id:              str = os.getenv("AUTH_ID", "")
    auth_password:        str = os.getenv("AUTH_PASSWORD", "")

@dataclass
class SpoolConfig:
    spool_dir:       str = os.getenv("SPOOL_DIR", "/var/lib/thermalwatcher/spool")
    max_spool_bytes: int = int(os.getenv("MAX_SPOOL_BYTES", str(200*1024*1024)))
    reconnect_max: int = int(os.getenv("RECONNECT_MAX_SEC", 60))
    

@dataclass
class LogConfig:
    log_file:  str = os.getenv("LOG_FILE", "/var/log/thermalwatcher.log")
    log_level: str = os.getenv("LOG_LEVEL", "INFO").upper()



