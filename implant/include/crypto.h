#ifndef CRYPTO_H
#define CRYPTO_H

// #include <windows.h>
#include "sodium/sodium.h"
#include <windows.h>
#include <wincrypt.h>
#include <string>
#include <vector>

/**
 * Create a new public/private key pair for use with LibSodium
 */
int createKeyPair(unsigned char *privateKeyOut, unsigned char *publicKeyOut);

/**
 * Create a new public/private key pair for use with LibSodium and base64 encode them
 */
int createBase64KeyPair(std::string *privateKeyOut, std::string *publicKeyOut);

/**
 * Base64 encode a string
 */
std::string base64Encode(const std::vector<unsigned char> &data);

/**
 * Base64 decode a string
 */
std::vector<unsigned char> base64Decode(const std::string &str);

#endif
