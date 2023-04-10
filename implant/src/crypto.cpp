#include "crypto.h"

int createKeyPair(unsigned char *privateKeyOut, unsigned char *publicKeyOut)
{
  return crypto_box_keypair(publicKeyOut, privateKeyOut);
}

int createBase64KeyPair(std::string *privateKeyOut, std::string *publicKeyOut)
{
  unsigned char privateKey[crypto_box_SECRETKEYBYTES];
  unsigned char publicKey[crypto_box_PUBLICKEYBYTES];
  int status = createKeyPair(privateKey, publicKey);
  if (status != 0)
  {
    return status;
  }
  *privateKeyOut = base64Encode(std::vector<unsigned char>(privateKey, privateKey + crypto_box_SECRETKEYBYTES));
  *publicKeyOut = base64Encode(std::vector<unsigned char>(publicKey, publicKey + crypto_box_PUBLICKEYBYTES));
  return 0;
}

std::string base64Encode(const std::vector<unsigned char> &data)
{
  DWORD outlen = 0;
  CryptBinaryToStringA(&data[0], static_cast<DWORD>(data.size()), CRYPT_STRING_BASE64 | CRYPT_STRING_NOCRLF, nullptr, &outlen);
  std::vector<char> out(outlen);
  CryptBinaryToStringA(&data[0], static_cast<DWORD>(data.size()), CRYPT_STRING_BASE64 | CRYPT_STRING_NOCRLF, &out[0], &outlen);
  return std::string(&out[0], outlen - 1); // -1 to remove the null terminator
}

std::vector<unsigned char> base64Decode(const std::string &str)
{
  DWORD outlen = 0;
  CryptStringToBinaryA(str.c_str(), static_cast<DWORD>(str.size()), CRYPT_STRING_BASE64, nullptr, &outlen, nullptr, nullptr);
  std::vector<unsigned char> out(outlen);
  CryptStringToBinaryA(str.c_str(), static_cast<DWORD>(str.size()), CRYPT_STRING_BASE64, &out[0], &outlen, nullptr, nullptr);
  return out;
}
