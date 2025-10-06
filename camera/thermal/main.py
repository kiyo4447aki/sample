from src.log import setup_logging
import sys
import signal
from src.controller import Controller
from src.config import LogConfig

def main():
    loggerConfig = LogConfig()
    logger = setup_logging(loggerConfig)
    logger.info("thermalwatcher uploader process start")
    ctrl = Controller()
    def sigterm(*_):
        logger.info("SIGTERM")
        try:
            ctrl.spooler.stop()
        except Exception:
            logger.exception("failed to stop spooler")

        try:
            ctrl.source.stop()
            ctrl.source.close()
        except Exception:
            logger.exception("failed to stop/close UVC")
        
        sys.exit(0)
    signal.signal(signal.SIGTERM, sigterm)
    ctrl.run_forever()
    
if __name__ == "__main__":
    main()