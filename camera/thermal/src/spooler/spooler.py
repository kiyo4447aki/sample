from uploader import Uploader
import logging
import threading
import os
from typing import Dict, Any
import json
import glob
import time
from uploader.uploader import ReserveUploadResp
import requests

class Spooler:
    META_VERSION = 1
    def __init__(
        self,
        dir_path: str,
        budget_bytes: int,
        uploader: Uploader,
        reconnect_max_sec: int = 60,
        isolation_keep_sec: int = 600,
        #15日ごとにトークンリフレッシュ
        auth_refresh_sec: int = 60*60*24*15
    ):
        self.log = logging.getLogger("thermalwatcher")
        self.dir = dir_path
        self.budget = budget_bytes
        self.uploader = uploader
        self.reconnect_max_sec = reconnect_max_sec
        self.isolation_keep_sec = isolation_keep_sec
        self.auth_refresh_sec = auth_refresh_sec
        
        self._stop_ev = threading.Event()
        self._th: threading.Thread | None = None
        self._last_auth_ts: float =0.0
        
        os.makedirs(self.dir, exist_ok=True)
    
    def enqueue(
        self, 
        event_id: str,
        jpeg_bytes: bytes,
        occurred_at_utc: str,
        severity: str,
        temp_c: float | None,
        width: int | None,
        height: int | None,
        sha256: str | None,
        mime_type: str = "image/jpeg",
    ) -> None:
        self._ensure_budget()
        
        jpg_path = self._jpg_path(event_id)
        meta_path = self._json_path(event_id)
        
        #TODO メタデータの型定義
        meta: Dict[str, Any] = {
            "v": self.META_VERSION,
            "event_id": event_id,
            "occurred_at": occurred_at_utc,
            "mime_type": mime_type,
            "severity": severity,
            "temp_c": temp_c,
            "width": width,
            "height": height,
            "sha256": sha256,
            "object_key": None,
            "uploaded": False,
            "created_at": int(time.time()),
            "last_error": None,
        }
        self._write_bytes(jpg_path, jpeg_bytes)
        self._write_json(meta_path, meta)
        
        self.log.info("spooled %s", event_id)
        
        self._ensure_budget()
        
    
    def start(self) -> None:
        if self._th and self._th.is_alive():
            return
        self._th = threading.Thread(target=self._loop, daemon=True)
        self._th.start()
    
    def stop(self) -> None:
        self._stop_ev.set()
        if self._th:
            self._th.join(timeout=5)
    
    def _json_path(self, event_id: str) -> str:
        return os.path.join(self.dir, f"{event_id}.json")
    
    def _jpg_path(self, event_id: str) -> str:
        return os.path.join(self.dir, f"{event_id}.jpg")
    
    def _getsize(self, p: str) -> int:
        try:
            return os.path.getsize(p)
        except OSError:
            return 0
    
    def _write_bytes(self, path: str, data: bytes) -> None:
        tmp = path + ".tmp"
        with open(tmp, "wb") as f:
            f.write(data)
        os.replace(tmp, path)
    
    def _write_json(self, path: str, obj: Dict[str, Any]) -> None:
        tmp = path + ".tmp"
        with open(tmp, "w", encoding="utf-8") as f:
            json.dump(obj, f, ensure_ascii=False)
        os.replace(tmp, path)
    
    def _load_json(self, path: str) -> Dict[str, Any] | None:
        try:
            with open(path, "r", encoding="utf-8") as f:
                return json.load(f)
        except Exception as e:
            self.log.warning("bad json %s: %s ", path, e)
            try:
                os.remove(path)
            except:
                self.log.warning("faild to remove json %s: %s ", path, e)
            return None
        
    def _ensure_budget(self) -> None:
        files = sorted(glob.glob(os.path.join(self.dir, "*.json")), key=os.path.getatime)
        def pair_size(json_path: str) -> int:
            event_id = os.path.splitext(os.path.basename(json_path))[0]
            return self._getsize(json_path) + self._getsize(self._jpg_path(event_id))
        
        #ファイルが容量を超えたとき、古いファイルを削除
        total = sum(pair_size(p) for p in files)
        while total > self.budget and files:
            old_json = files.pop(0)
            event_id = os.path.splitext(os.path.basename(old_json))[0]
            old_jpg = self._jpg_path(event_id)
            try:
                if os.path.exists(old_json):
                    os.remove(old_json)
                if os.path.exists(old_jpg):
                    os.remove(old_jpg)
                self.log.warning("spool budget: remove old %s", event_id)
            except Exception as e:
                self.log.error("spool budget removal error: %s", e)
            total = sum(pair_size(p) for p in files)
            
        #孤立したjpegを削除
        now = time.time()
        for jpg in glob.glob(os.path.join(self.dir, "*.jpg")):
            event_id = os.path.splitext(os.path.basename(jpg))[0]
            jpath = self._json_path(event_id)
            if not os.path.exists(jpath) and (now - os.path.getmtime(jpg) > self.isolation_keep_sec):
                try:
                    os.remove(jpg)
                    self.log.warning("removed isolated jpg %s", jpg)
                except Exception as e:
                    self.log.error("remove isolated jpg failed: %s", e)
                    
        #孤立したjsonを削除
        for jpath in glob.glob(os.path.join(self.dir, "*.json")):
            event_id = os.path.splitext(os.path.basename(jpath))[0]
            jpg = self._jpg_path(event_id)
            if not os.path.exists(jpg) and (now - os.path.getmtime(jpath) > self.isolation_keep_sec):
                try:
                    os.remove(jpath)
                    self.log.warning("removed isolated json %s", jpath)
                except Exception as e:
                    self.log.error("remove isolated json failed: %s", e)
    
    def _auth_maybe_refresh(self, force: bool = False) -> None:
        now = time.time()
        need = force or (self.uploader.token is None) or (now - float(self._last_auth_ts) > self.auth_refresh_sec)
        if not need:
            return
        try:
            self.uploader.authorize()
            self._last_auth_ts = now
            self.log.debug("authorized/refreshed token")
        except Exception as e:
            self.log.warning("authorize failed: %s",e)
    
    def _loop(self) -> None:
        backoff = 5
        while not self._stop_ev.is_set():
            try:
                progressed = False
                
                #認証が必要かチェック、必要なら認証
                self._auth_maybe_refresh()
                
                for meta_path in sorted(glob.glob(os.path.join(self.dir, "*.json")), key=os.path.getmtime):
                    if self._stop_ev.is_set():
                        break
                    meta = self._load_json(meta_path)
                    if meta is None:
                        progressed = True
                        continue
                    
                    event_id = meta.get("event_id")
                    if not event_id:
                        os.remove(meta_path)
                        self.log.warning("drop invalid meta(no event_id): %s", meta_path)
                        progressed = True
                        continue
                    
                    jpg = self._jpg_path(event_id)
                    if not os.path.exists(jpg):
                        age = time.time() - os.path.getmtime(meta_path)
                        if age > self.isolation_keep_sec:
                            os.remove(meta_path)
                            self.log.warning("isolated json(no jpeg) dropped: %s", meta_path)
                            progressed = True
                        continue
                    
                    done, updated = self._proccess_one(meta, jpg)
                    if updated is not None:
                        self._write_json(meta_path, updated)
                    if done:
                        try:
                            os.remove(meta_path)
                        except Exception as e:
                            self.log.error("remove completed json failed: %s", e)
                        try:
                            os.remove(jpg)
                        except Exception as e:
                            self.log.error("remove completed jpg failed: %s", e)
                        progressed = True
                
                if progressed:
                    #タスク成功時のみバックオフをリセット
                    #その後バックオフ秒wait
                    backoff = 5
                self._stop_ev.wait(backoff)
                backoff = min(self.reconnect_max_sec, int(backoff * 1.5)) if backoff > 0 else 5
            except Exception as e:
                self.log.exception("spool worker fatal Exception: %s", e)
                self._stop_ev.wait(10)
                
    
    def _proccess_one(self, meta: Dict[str, Any], jpg_path: str) -> tuple[bool, Dict[str, Any] | None]:
        """
        returns:
            done: True=完了（メタ削除）、False=保留（再試行）
            new_meta: 更新がある場合はDict、ない場合はNone
        """
        new_meta = dict(meta)
        new_meta["last_error"] = None
        self._auth_maybe_refresh()
        
        mime = meta.get("mime_type", "image/jpeg")
        #アラート予約
        try:
            rsv = self._reserve_with_reauth(
                meta["occurred_at"], meta["event_id"], mime,
                meta.get("width", None), meta.get("height", None), meta.get("sha256", None)
                )
        except requests.HTTPError as he:
            code = getattr(he.response, "status_code", None)
            """
            4xx系エラー（認証以外）は際試行せずに破棄
            認証系エラーとその他エラーは再試行
            """
            if code in (400, 404, 422):
                new_meta["last_error"] = f"reserve http {code}"
                self.log.warning("reserve non-retry http %s for %s -> drop", code, meta["event_id"])
                return True, new_meta
            self.log.warning("reserve retryable error(%s): %s", code, he)
            new_meta["last_error"] = f"reserve http {code}"
            return False, new_meta
        except Exception as e:
            self.log.warning("reserve error: %s", e)
            new_meta["last_error"] = "reserve error"
            return False, new_meta
        
        object_key = rsv.get("object_key")
        if object_key and meta.get( "object_key") != object_key:
            self.log.debug("object_key updated %s -> %s (%s)", meta.get("object_key"), object_key, meta["event_id"])
            new_meta["object_key"] = object_key
        
        put_url = rsv["put_url"]
        
        #過去に画像アップロートを行っていない場合のみアップロード
        if not meta.get("uploaded", False):
            try:
                with open(jpg_path, "rb") as jf:
                    self.uploader.put_image(put_url, jf.read())
                new_meta["uploaded"] = True
            except requests.HTTPError as he:
                code = getattr(he.response, "status_code", None)
                if code in (400, 404, 409, 412, 422):
                    self.log.warning("put non-retry http %s for %s -> drop", code, meta["event_id"])
                    new_meta["last_error"] = f"put http {code}"
                    return True, new_meta
                self.log.warning("put retryable error(%s): %s", code, he)
                new_meta["last_error"] = f"put http {code}"
                return False, new_meta
            except Exception as e:
                self.log.warning("put error: %s", e)
                new_meta["last_error"] = "put error"
                return False, new_meta
        
        try:
            self._fire_with_reauth(
                new_meta["object_key"],
                meta.get("severity", "info"),
                meta["event_id"],
                mime,
                meta["occurred_at"],
                meta.get("temp_c")
            )
            return True, new_meta
        except requests.HTTPError as he:
            code = getattr(he.response, "status_code", None)
            if code in (400, 404, 409, 422):
                self.log.warning("fire non-retry http %s for %s -> drop", code, meta["eid"])
                new_meta["last_error"] = f"fire http {code}"
                return True, new_meta
            self.log.warning("fire retryable error(%s): %s", code, he)
            new_meta["last_error"] = f"fire http {code}"
            return False, new_meta
        except Exception as e:
            self.log.warning("fire error: %s", e)
            new_meta["last_error"] = "fire error"
            return False, new_meta
                


    #アラート予約APIを叩き、権限エラーの場合のみ再認証→再試行
    def _reserve_with_reauth(self, occurred_at: str, event_id: str, mime: str, width: int|None, height: int|None, sha256: str|None) -> ReserveUploadResp:
        try:
            if self.uploader.token is None:
                self.uploader.authorize()
                self._last_auth_ts = time.time()
            return self.uploader.reserve_upload(
                occurred_at, event_id, width, height, sha256, mime
            )
        except requests.HTTPError as he:
            if getattr(he.response, "status_code") in (401, 403):
                self.log.info("reserve got %s -> reauth", he.response.status_code)
                self.uploader.authorize()
                self._last_auth_ts = time.time()
                return self.uploader.reserve_upload(
                    occurred_at, event_id, width, height, sha256, mime
                )
            #TODO エラー定義
            raise Exception
    #アラート確定APiを叩き、権限エラーの場合のみ再認証→再試行
    def _fire_with_reauth(self, object_key: str, severity: str, event_id: str, mime_type: str, occurred_at: str, temp_c: float|None) -> None:
        try:
            if self.uploader.token is None:
                self.uploader.authorize()
                self._last_auth_ts = time.time()
            return self.uploader.fire_alert(object_key, severity, event_id, mime_type, occurred_at, temp_c)
        except requests.HTTPError as he:
            if getattr(he.response, "status_code") in (401, 403):
                self.log.info("reserve got %s -> reauth", he.response.status_code)
                self.uploader.authorize()
                self._last_auth_ts = time.time()
                return self.uploader.fire_alert(object_key, severity, event_id, mime_type, occurred_at, temp_c)    
            #TODO エラー定義
            raise Exception
