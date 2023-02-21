"""
Functions to parse configuration files for the listeners
"""

import os
import yaml

CONFIG_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "config")

def get_config(config_file: str):
    with open(os.path.join(CONFIG_DIR, config_file)) as f:
        return yaml.safe_load(f)

def get_admin_creds():
    config = get_config(os.path.join("admin", "admin.yml"))
    return config["username"], config["auth_token"]
