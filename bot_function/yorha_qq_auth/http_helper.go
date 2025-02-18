package yorha_qq_auth

import (
	yorha_defines "EndlessEmbrace/bot_function/yorha_qq_auth/defines"
	yorha_rsa "EndlessEmbrace/depends/YoRHa/rsa"
	"bytes"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ...
func UnmarshalClientRequest[T any](body io.ReadCloser, key *rsa.PrivateKey) (result T, err error) {
	// read content from body
	hexBody, err := io.ReadAll(body)
	if err != nil {
		err = fmt.Errorf("UnmarshalClientRequest: %v", err)
		return
	}
	// decode to byte
	bodyEncryptedBytes, err := hex.DecodeString(string(hexBody))
	if err != nil {
		err = fmt.Errorf("UnmarshalClientRequest: %v", err)
		return
	}
	// decrypt body
	bodyBytes, err := yorha_rsa.DecryptPKCS1v15(key, bodyEncryptedBytes)
	if err != nil {
		err = fmt.Errorf("UnmarshalClientRequest: %v", err)
		return
	}
	// unmarshal
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		err = fmt.Errorf("UnmarshalClientRequest: %v", err)
		return
	}
	// return
	return
}

// ...
func PostJSON[T any](url string, key *rsa.PublicKey, object any) (result T, err error) {
	// marshal body
	originBody, err := json.Marshal(object)
	if err != nil {
		err = fmt.Errorf("PostJSON: %v", err)
		return
	}
	// encrypt and do hex
	encryptedBody, err := yorha_rsa.EncryptPKCS1v15(key, originBody)
	if err != nil {
		err = fmt.Errorf("PostJSON: %v", err)
		return
	}
	body := hex.EncodeToString(encryptedBody)
	// post and get response
	resp, err := http.Post(url, "qq-auth", bytes.NewBuffer([]byte(body)))
	if err != nil {
		err = fmt.Errorf("PostJSON: %v", err)
		return
	}
	defer resp.Body.Close()
	// read resp body
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("PostJSON: %v", err)
		return
	}
	// unmarshal resp body
	err = json.Unmarshal(respBodyBytes, &result)
	if err != nil {
		err = fmt.Errorf("PostJSON: %v", err)
		return
	}
	// return
	return
}

// ...
func WriteResponse(c *gin.Context, key *rsa.PublicKey, response yorha_defines.ServerResponse) {
	packetBytes, _ := json.Marshal(response)
	encryptedBody, _ := yorha_rsa.EncryptPKCS1v15(key, packetBytes)
	c.Writer.Write([]byte(encryptedBody))
}
