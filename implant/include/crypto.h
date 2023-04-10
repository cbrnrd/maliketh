#ifndef CRYPTO_H
#define CRYPTO_H

//#include <windows.h>
#include "sodium/sodium.h"
#include <string>
#include <string.h>

/**
 * Create a new public/private key pair for use with LibSodium
 */
int createKeyPair(unsigned char *privateKeyOut, unsigned char *publicKeyOut);

/**
 * Base64 encode a string
 */
std::string base64Encode(const char *input);

/**
 * Base64 decode a string
 */
std::string base64Decode(std::string input);

#endif
