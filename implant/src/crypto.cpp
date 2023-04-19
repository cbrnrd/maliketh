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

std::vector<BYTE> encrypt(std::vector<BYTE> pubKey, std::vector<BYTE> privKey, std::vector<BYTE> message)
{
	std::vector<BYTE> cipherText(message.size() + crypto_box_MACBYTES);
	std::vector<BYTE> nonce(crypto_box_NONCEBYTES);
	randombytes_buf(nonce.data(), crypto_box_NONCEBYTES);
	crypto_box_easy(cipherText.data(), message.data(), message.size(), nonce.data(), pubKey.data(), privKey.data());
	cipherText.insert(cipherText.begin(), nonce.begin(), nonce.end());
	return cipherText;
}

std::string encryptB64String(std::string pubKey, std::string privKey, std::string message)
{
	std::vector<BYTE> pubKeyBytes = base64Decode(pubKey);
	std::vector<BYTE> privKeyBytes = base64Decode(privKey);
	std::vector<BYTE> messageBytes(message.begin(), message.end());
	std::vector<BYTE> cipherText = encrypt(pubKeyBytes, privKeyBytes, messageBytes);
	return base64Encode(cipherText);
}

std::string encryptB64SecretBox(std::string key, std::string message)
{
	std::vector<BYTE> keyBytes = base64Decode(key);
	std::vector<BYTE> messageBytes(message.begin(), message.end());
	std::vector<BYTE> cipherText = encryptSecretBox(keyBytes, messageBytes);
	return base64Encode(cipherText);
}

std::vector<BYTE> encryptSecretBox(std::vector<BYTE> key, std::vector<BYTE> message)
{
	std::vector<BYTE> cipherText(message.size() + crypto_secretbox_MACBYTES);
	std::vector<BYTE> nonce(crypto_secretbox_NONCEBYTES);
	randombytes_buf(nonce.data(), crypto_secretbox_NONCEBYTES);
	crypto_secretbox_easy(cipherText.data(), message.data(), message.size(), nonce.data(), key.data());
	cipherText.insert(cipherText.begin(), nonce.begin(), nonce.end());
	return cipherText;
}

std::string decrypt(std::vector<BYTE> pubKey, std::vector<BYTE> privKey, std::vector<BYTE> cipherText)
{
	std::vector<BYTE> message(cipherText.size() - crypto_box_MACBYTES);
	// Get nonce from ciphertext
	std::vector<BYTE> nonce(cipherText.begin(), cipherText.begin() + crypto_box_NONCEBYTES);
	crypto_box_open_easy(message.data(), cipherText.data() + crypto_box_NONCEBYTES, cipherText.size() - crypto_box_NONCEBYTES, nonce.data(), pubKey.data(), privKey.data());
	return std::string(message.begin(), message.end());
}

/**
 * Decrypt a base64 encoded string and return a string.
 */
std::string decryptB64String(std::string pubKey, std::string privKey, std::string cipherText)
{
	std::vector<BYTE> pubKeyBytes = base64Decode(pubKey);
	std::vector<BYTE> privKeyBytes = base64Decode(privKey);
	std::vector<BYTE> cipherTextBytes = base64Decode(cipherText);
	return decrypt(pubKeyBytes, privKeyBytes, cipherTextBytes);
}

std::string base64Encode(const std::vector<BYTE> &data)
{
	DWORD length = 0;
	CryptBinaryToStringA(data.data(), static_cast<DWORD>(data.size()), CRYPT_STRING_BASE64 | CRYPT_STRING_NOCRLF, NULL, &length);

	std::string encoded(length, 0);
	CryptBinaryToStringA(data.data(), static_cast<DWORD>(data.size()), CRYPT_STRING_BASE64 | CRYPT_STRING_NOCRLF, &encoded[0], &length);

	return encoded;
}

std::string base64EncodeString(const std::string &data)
{
	DWORD length = 0;
	CryptBinaryToStringA(reinterpret_cast<const BYTE *>(data.c_str()), static_cast<DWORD>(data.size()), CRYPT_STRING_BASE64 | CRYPT_STRING_NOCRLF, NULL, &length);

	std::string encoded(length, 0);
	CryptBinaryToStringA(reinterpret_cast<const BYTE *>(data.c_str()), static_cast<DWORD>(data.size()), CRYPT_STRING_BASE64 | CRYPT_STRING_NOCRLF, &encoded[0], &length);

	return encoded;
}

std::vector<BYTE> base64Decode(const std::string &encoded)
{
	DWORD length = 0;
	CryptStringToBinaryA(encoded.c_str(), static_cast<DWORD>(encoded.length()), CRYPT_STRING_BASE64, NULL, &length, NULL, NULL);

	std::vector<BYTE> decoded(length, 0);
	CryptStringToBinaryA(encoded.c_str(), static_cast<DWORD>(encoded.length()), CRYPT_STRING_BASE64, decoded.data(), &length, NULL, NULL);

	return decoded;
}

std::string base64DecodeToString(const std::string &base64String)
{
	DWORD dwDecodedLength = 0;
	CryptStringToBinaryA(base64String.c_str(), base64String.length(), CRYPT_STRING_BASE64, nullptr, &dwDecodedLength, nullptr, nullptr);

	std::vector<BYTE> buffer(dwDecodedLength);
	CryptStringToBinaryA(base64String.c_str(), base64String.length(), CRYPT_STRING_BASE64, buffer.data(), &dwDecodedLength, nullptr, nullptr);

	return std::string(buffer.begin(), buffer.end());
}
