
import requests
from typing import TypedDict

class AuthReq(TypedDict):
    username: str
    password: str


class ReserveUploadReq(TypedDict):
    mime_type: str
    width: int | None
    height: int | None
    sha256: str | None
    occurred_at: str
    event_id : str

class ReserveUploadResp(TypedDict):
    status: str
    object_key: str
    put_url: str
    expires_in: int

class FireAlertReq(TypedDict):
    severity: str
    temp_c: float | None
    occurred_at: str
    event_id: str
    mime_type: str
    object_key: str

class Uploader:
    def __init__(self, sess: requests.Session, auth_url: str, alerts_url: str, user_id: str, password: str, upload_url: str):
        self.sess = sess
        self.auth_url = auth_url
        self.alerts_url = alerts_url
        self.user_id = user_id
        self.password = password
        self.uploads_url = upload_url
        self.token: str | None = None
        
    def _validate_mime(self, mime_type: str) -> None:
        if mime_type not in ("image/jpeg", "image/png"):
            raise ValueError(f"unsupported mime_type: {mime_type}")
    
    def authorize(self) -> None:
        headers = {"Content-Type": "application/json"}
        body: AuthReq = {"username": self.user_id, "password": self.password}
        resp = self.sess.post(self.auth_url, headers=headers, json=body, timeout=10)
        resp.raise_for_status()
        j = resp.json()
        if "token" not in j:
            raise RuntimeError("failed to auth")
        self.token = j["token"]
    
    def reserve_upload(self, occurred_at: str, event_id: str, width: int|None, height: int|None, sha256: str|None, mime: str="image/jpeg") -> ReserveUploadResp:
        if self.token is None:
            raise RuntimeError("can't upload without authorization")
        self._validate_mime(mime)
        headers = {"Authorization": f"Bearer {self.token}", "Content-Type": "application/json"}
        #TODO 画像の詳細情報を渡す実装
        body: ReserveUploadReq = {
            "mime_type": mime,
            "width": width,
            "height": height,
            "sha256": sha256,
            "occurred_at": occurred_at,
            "event_id": event_id
            }
        resp = self.sess.post(self.uploads_url, headers=headers, json=body, timeout=10)
        resp.raise_for_status()
        j = resp.json()
        for k in ("status", "object_key", "put_url", "expires_in"):
            if k not in j:
                raise RuntimeError(f"bad reserve response: {j}")
        if j["status"] != "success":
            raise RuntimeError(f"reserve failed: {j}")
        r: ReserveUploadResp = {
            "status": j["status"],
            "object_key": j["object_key"], 
            "put_url": j["put_url"],
            "expires_in": j["expires_in"]
            }
        return r
    
    def put_image(self, put_url: str, jpeg_bytes: bytes, mime_type: str = "image/jpeg") -> None:
        self._validate_mime(mime_type)
        r = self.sess.put(put_url, headers={"Content-Type": mime_type}, data=jpeg_bytes, timeout=20)
        r.raise_for_status()
    
    def fire_alert(self, object_key: str, severity: str, event_id: str, mime_type: str, occurred_at: str, temp_c: float|None) -> None:
        if self.token is None:
            raise RuntimeError("can't upload without authorization")
        self._validate_mime(mime_type)
        headers = {
            "Authorization": f"Bearer {self.token}",
            "Content-Type": "application/json",
        }
        payload: FireAlertReq = {
            "severity": severity,
            "object_key": object_key,
            "event_id": event_id,
            "mime_type": mime_type,
            "occurred_at": occurred_at,
            "temp_c": temp_c,
            }
        r = self.sess.post(self.alerts_url, headers=headers, json=payload, timeout=10)
        r.raise_for_status()
        
    
    def send_all(self, jpeg_bytes: bytes,occurred_at: str, event_id: str, mime_type: str,
                severity: str, temp_c: float|None=None, sha256: str|None=None,
                width: int|None=None, height:int|None=None) -> None:
        #初回のみ認証、以降の認証はSpoolerクラスにて実行
        if self.token is None:
            self.authorize()
        r = self.reserve_upload(
            occurred_at, event_id, width, height, sha256, mime_type)
        self.put_image(r["put_url"], jpeg_bytes, mime_type)
        self.fire_alert(r["object_key"], severity, event_id, mime_type, occurred_at, temp_c)
    