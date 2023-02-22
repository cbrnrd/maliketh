import os
import sys
from maliketh.buildapp import build_app, init_db
from maliketh.crypto.x509 import generate_ecc_keypair, check_ecc_keypair
from maliketh.logging.standard_logger import StandardLogger, LogLevel

#init_db()

app  = build_app()

def init_simple_logger(level: LogLevel=LogLevel.INFO) -> StandardLogger:
  return StandardLogger(sys.stdout, sys.stderr, level)
  

def entry():
  logger = init_simple_logger()
  logger.info("Starting up...")
  logger.debug("Checking for existing keypair...")
  if check_ecc_keypair(os.path.join(os.getcwd(), "config", "admin", "certs")):
    logger.debug("Keypair found!")
  else:
    logger.debug("Keypair not found, generating...")
    generate_ecc_keypair(os.path.join(os.getcwd(), "config", "admin", "certs"))
    logger.debug("Keypair generated!")

def run_app():
  app.run()

if __name__ == "__main__":
  entry()
  run_app()
