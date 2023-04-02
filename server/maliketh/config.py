"""
Functions to parse configuration files for the listeners
"""

import os
from typing import List, Optional
import yaml

CONFIG_DIR = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "config"))
SERVER_PUB_KEY_PATH = os.path.join(CONFIG_DIR, "admin", "certs", "server_pub")
SERVER_PRIV_KEY_PATH = os.path.join(CONFIG_DIR, "admin", "certs", "server_priv")
DEFAULT_OP_ROUTES_CONFIG = os.path.join(CONFIG_DIR, "admin", "routes.yaml")
DEFAULT_C2_PROFILE = os.path.join(CONFIG_DIR, "profiles", "default.yaml")

OP_ROUTES = yaml.safe_load(open(os.path.join(CONFIG_DIR, "admin", "routes.yaml")))  # Operator listener Flask routes
C2_PROFILE = yaml.safe_load(open(DEFAULT_C2_PROFILE))  # C2 listener Flask routes
C2_BASE_PATH = C2_PROFILE['server']['routes']['base_path']

def set_c2_profile(path: str):
    global C2_PROFILE
    C2_PROFILE = yaml.safe_load(open(path))

def get_c2_route(route: str) -> str:
    return C2_PROFILE['server']['routes'].get(route, {}).get('path', '')

def get_c2_route_methods(route: str) -> List[str]:
    return C2_PROFILE['server']['routes'].get(route, {}).get('methods', [])

def get_c2_server_option(option: str) -> str:
    return C2_PROFILE['server'].get(option, '')


def get_config(config_file: str):
    with open(os.path.join(CONFIG_DIR, config_file)) as f:
        return yaml.safe_load(f)

def get_admin_creds():
    config = get_config(os.path.join("admin", "creds.yaml"))
    return config["username"], config["auth_token"]
