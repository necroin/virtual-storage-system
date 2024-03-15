import os
import subprocess
import tenacity
import io
from modules.server import Server


class Project:
    url = "https://localhost:3301"

    def __init__(self, path: str = "tmp") -> None:
        self.path = path

        self.server: Server = Server(self.url)

        self.process = subprocess.Popen(
            "make.bat",
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
        )

    @tenacity.retry(wait=tenacity.wait_fixed(1), stop=tenacity.stop_after_delay(30))
    def wait_start(self):
        assert self.server.status() == "OK"

    def stop(self):
        self.process.terminate()
