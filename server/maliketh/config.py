"""
Functions to parse configuration files for the listeners
"""

import os
import yaml

CONFIG_DIR = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "config"))
SERVER_PUB_KEY_PATH = os.path.join(CONFIG_DIR, "admin", "certs", "server_pub")
SERVER_PRIV_KEY_PATH = os.path.join(CONFIG_DIR, "admin", "certs", "server_priv")

ROUTES = yaml.safe_load(open(os.path.join(CONFIG_DIR, "admin", "routes.yaml")))  # Parse once, use many times

def get_config(config_file: str):
    with open(os.path.join(CONFIG_DIR, config_file)) as f:
        return yaml.safe_load(f)

def get_admin_creds():
    config = get_config(os.path.join("admin", "creds.yaml"))
    return config["username"], config["auth_token"]
