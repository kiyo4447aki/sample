import numpy as np


def ktoc(val_u16: np.ndarray) -> np.ndarray:
    """
    サーモカメラからのデータ処理用
    Tliner仕様のセンチケルビン配列を摂氏配列に変換
    """
    return (val_u16.astype(np.float32) - 27315.0) / 100.0

def choose_severity(maxC: float, abs_thresh: float) -> str:
    if maxC >= abs_thresh + 20: return "critical"
    if maxC >= abs_thresh:      return "warning"
    return "info"
