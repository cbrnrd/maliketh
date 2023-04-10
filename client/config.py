from dataclasses import dataclass, asdict
import json
from typing import Optional
from cli.logging import LogLevel
from crypto import base64_encrypt_and_sign_str

# Globals
log_level: LogLevel = LogLevel.INFO


@dataclass
class OperatorConfig:
    name: str
    c2: str
    c2_port: int
    login_secret: str
    secret: str
    public: str
    signing_key: str
    verify_key: str
    server_pub: str
    rmq_queue: str
    auth_token: Optional[str] = None

    def to_dict(self) -> dict:
        return asdict(self)

    def to_json(self) -> str:
        return json.dumps(self.to_dict())

    @staticmethod
    def from_json(json_str: str) -> "OperatorConfig":
        obj = json.loads(json_str)
        return OperatorConfig.from_dict(obj)

    @staticmethod
    def from_dict(args: dict) -> "OperatorConfig":
        return OperatorConfig(
            name=args["name"],
            c2=args["c2"],
            c2_port=args["c2_port"],
            login_secret=args["login_secret"],
            secret=args["secret"],
            public=args["public"],
            signing_key=args["signing_key"],
            verify_key=args["verify_key"],
            server_pub=args["server_pub"],
            rmq_queue=args["rmq_queue"],
        )

    def enc_and_sign_secret(self) -> str:
        return base64_encrypt_and_sign_str(
            encryption_key_b64=self.secret,
            recv_pub_key_b64=self.server_pub,
            signing_key_b64=self.signing_key,
            data=self.login_secret,
        )
