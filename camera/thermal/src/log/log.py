import logging
import os
from logging.handlers import RotatingFileHandler
import sys
from config import LogConfig

def setup_logging(cfg: LogConfig):
    logger = logging.getLogger("thermalwatcher")
    
    log_level = getattr(logging, cfg.log_level, logging.INFO)
    logger.setLevel(log_level)
    
    if logger.handlers:
        for h in list(logger.handlers):
            logger.removeHandler(h)
            
    logger.propagate = False
    
    fmt = logging.Formatter(("%(asctime)s %(levelname)s %(message)s"))

    os.makedirs(os.path.dirname(cfg.log_file),exist_ok=True)
    
    file_handler = RotatingFileHandler(
        cfg.log_file, 
        maxBytes=5*1024*1024, 
        backupCount=3,
        encoding="utf-8",
        delay=True
        )
    file_handler.setFormatter(fmt)
    logger.addHandler(file_handler)
    
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setFormatter(fmt)
    logger.addHandler(console_handler)
    
    return logger


