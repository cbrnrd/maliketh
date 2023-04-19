#pragma once
#ifndef CRYPTO_H
#define CRYPTO_H

// #include <windows.h>
//#include "sodium/sodium.h"
#include <sodium.h>
#include <windows.h>
#include <wincrypt.h>
#include <string>
#include <vector>

#pragma comment(lib, "crypt32.lib")

/**
 * Create a new public/private key pair for use with LibSodium
 */
int createKeyPair(unsigned char* privateKeyOut, unsigned char* publicKeyOut);

/**
 * Create a new public/private key pair for use with LibSodium and base64 encode them
 */
int createBase64KeyPair(std::string* privateKeyOut, std::string* publicKeyOut);

std::vector<BYTE> encrypt(std::vector<BYTE> pubKey, std::vector<BYTE> privKey, std::vector<BYTE> message);

std::string encryptB64SecretBox(std::string key, std::string message);

std::vector<BYTE> encryptSecretBox(std::vector<BYTE> key, std::vector<BYTE> message);

/**
 * Encrypt a string and return its base64 encoded representation
*/
std::string encryptB64String(std::string pubKey, std::string privKey, std::string message);

std::string decrypt(std::vector<BYTE> pubKey, std::vector<BYTE> privKey, std::vector<BYTE> cipherText);

std::string decryptB64String(std::string pubKey, std::string privKey, std::string cipherText);

/**
 * Base64 encode a string
 */
std::string base64Encode(const std::vector<BYTE>& data);

std::string base64EncodeString(const std::string& data);

/**
 * Base64 decode a string
 */
std::vector<BYTE> base64Decode(const std::string& encoded);

std::string base64DecodeToString(const std::string& base64String);

#endif
