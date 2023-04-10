#include "crypto.h"

int createKeyPair(unsigned char *privateKeyOut, unsigned char *publicKeyOut)
{
  return crypto_box_keypair(publicKeyOut, privateKeyOut);
}

std::string base64Encode(const char *input) {
  size_t inputSize = strlen(input);
  std::string encoded;
  int variant = sodium_base64_VARIANT_ORIGINAL;
  size_t maxLen = sodium_base64_ENCODED_LEN(inputSize, variant);
  encoded.resize(maxLen);
  sodium_bin2base64(&encoded[0], maxLen, (const unsigned char *)input, inputSize, variant);
  return encoded;
}
