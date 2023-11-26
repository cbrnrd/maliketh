import logging
import structlog
import os

class GunicornLogger(object):
    """
    A stripped down version of https://github.com/benoitc/gunicorn/blob/master/gunicorn/glogging.py to provide structlog logging in gunicorn
    Modified from http://stevetarver.github.io/2017/05/10/python-falcon-logging.html
    """

    def __init__(self, cfg):
        self._error_logger = structlog.get_logger("gunicorn.error")
        self._error_logger.setLevel(logging.INFO)
        self._access_logger = structlog.get_logger("gunicorn.access")
        self._access_logger.setLevel(logging.INFO)
        self.cfg = cfg

    def critical(self, msg, *args, **kwargs) -> None:
        self._error_logger.error(msg, *args, **kwargs)

    def error(self, msg, *args, **kwargs) -> None:
        self._error_logger.error(msg, *args, **kwargs)

    def warning(self, msg, *args, **kwargs) -> None:
        self._error_logger.warning(msg, *args, **kwargs)

    def info(self, msg, *args, **kwargs) -> None:
        self._error_logger.info(msg, *args, **kwargs)

    def debug(self, msg, *args, **kwargs) -> None:
        self._error_logger.debug(msg, *args, **kwargs)

    def exception(self, msg, *args, **kwargs) -> None:
        self._error_logger.exception(msg, *args, **kwargs)

    def log(self, lvl, msg, *args, **kwargs) -> None:
        self._error_logger.log(lvl, msg, *args, **kwargs)

    def access(self, resp, req, environ, request_time) -> None:
        status = resp.status
        if isinstance(status, str):
            status = status.split(None, 1)[0]

        self._access_logger.info(
            "request",
            method=environ["REQUEST_METHOD"],
            request_uri=environ["RAW_URI"],
            status=status,
            response_length=getattr(resp, "sent", None),
            request_time_seconds="%d.%06d" % (request_time.seconds, request_time.microseconds),
            pid="<%s>" % os.getpid(),
        )

    def reopen_files(self) -> None:
        pass  # we don't support files

    def close_on_exec(self) -> None:
        pass  # we don't support files