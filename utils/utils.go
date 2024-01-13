package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

func GenerateID(size int) (string, error) {
	id := make([]byte, 4)
	_, err := rand.Read(id)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(id), nil

}

func CreateHTTPRequest(method, uri, port, endPoint string, payload interface{}) (*http.Request, error) {

	url := fmt.Sprintf("%s:%s/%s", uri, port, endPoint)

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(payload)

	fmt.Printf("Utils: Sending request to %s\n", url)
	fmt.Printf("Utils: Request Data\n%+v\n", reqBodyBytes)
	return http.NewRequest(method, url, bytes.NewBuffer(reqBodyBytes.Bytes()))
}

func DispatchRequest(request *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(request)
}
