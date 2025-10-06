import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

_http_session: requests.Session | None = None

def get_session() -> requests.Session:
    global _http_session
    if _http_session is None:
        s = requests.Session()
        retries = Retry(
            total=5,
            backoff_factor=1.0,
            status_forcelist=[429, 500, 502, 503, 504],
            allowed_methods=["GET", "POST", "PUT"],
        )
        s.mount("http://", HTTPAdapter(max_retries=retries))
        s.mount("https://", HTTPAdapter(max_retries=retries))
        _http_session = s
    return _http_session