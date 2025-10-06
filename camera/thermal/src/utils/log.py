def format_bytes(n):
    for u in ["B", "KB", "MB", "GB"]:
        if n < 1024 :
            return f"{n:.1f}{u}"
        n /= 1024
    return f"{n:.1f}TB"