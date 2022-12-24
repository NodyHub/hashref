package hashref

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NodyHub/hashref/pkg/util"
)

type HashrefClient struct {
	config Config
}

func NewClient(config Config) HashrefClient {
	return HashrefClient{config: config}
}

func (hc *HashrefClient) GetRemoteData(inputType HashType, input, hashValue string) (bool, map[string]interface{}) {
	log.Printf("Request data for %v %v\n", Lookup[inputType], hashValue)
	remoteData := make(map[string]interface{})

	// prepare get request
	requestUri := fmt.Sprintf("%v/api/hash/%v", hc.config.HashrefServer, hashValue)
	log.Printf("Request-uri: %v\n", requestUri)
	req, err := http.NewRequest(http.MethodGet, requestUri, nil)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return false, map[string]interface{}{
			"error": err.Error(),
		}
	}
	req.Header.Add("Authorization", fmt.Sprintf("%v", hc.config.Publisher))

	// Create Client
	client := &http.Client{}

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return false, map[string]interface{}{
			"error": err.Error(),
		}
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log.Printf("Response Status %v: %v\n", resp.StatusCode, resp.Status)
		remoteData["status"] = resp.Status
		return false, remoteData
	}
	if err := json.Unmarshal(body, &remoteData); err != nil {
		remoteData["status"] = err.Error()
		return false, remoteData
	}
	return true, remoteData
}

func (hc *HashrefClient) GetRemoteDataFromPublisher(inputType HashType, input, hashValue, publisher string) (bool, map[string]interface{}) {
	log.Printf("Request data for t:%v h:%v p:%v\n", Lookup[inputType], hashValue, publisher)

	// prepare get request
	requestUri := fmt.Sprintf("%v/api/hash/%v/publisher/%v", hc.config.HashrefServer, hashValue, publisher)
	log.Printf("Request-uri: %v\n", requestUri)
	req, err := http.NewRequest(http.MethodGet, requestUri, nil)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return false, map[string]interface{}{
			"error": err.Error(),
		}
	}
	req.Header.Add("Authorization", fmt.Sprintf("%v", hc.config.Publisher))

	// Create Client
	client := &http.Client{}

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return false, map[string]interface{}{
			"error": err.Error(),
		}
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log.Printf("Response Status %v: %v\n", resp.StatusCode, resp.Status)
		return false, map[string]interface{}{
			"status": resp.Status,
			"code":   resp.StatusCode,
		}
	}

	remoteData := make(map[string]interface{})
	if err := json.Unmarshal(body, &remoteData); err != nil {
		return false, map[string]interface{}{
			"error": err.Error(),
		}
	}
	return true, remoteData
}

// RemoveHash deletes the metadata remotly to the provided hash.
func (hc *HashrefClient) RemoveHash(force bool, input, calculatedHash string) bool {
	if !force && !util.YesOrNoQuestion(fmt.Sprintf("Should %v really be removed from hashrev?", input)) {
		return false
	}
	log.Printf("Delete metadata for %v\n", calculatedHash)

	// prepare delete request
	requestUri := fmt.Sprintf("%v/api/hash/%v", hc.config.HashrefServer, calculatedHash)
	req, err := http.NewRequest(http.MethodDelete, requestUri, nil)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return false
	}

	// Create Client
	client := &http.Client{}

	// Set request meta data
	req.Header.Set("Authorization", fmt.Sprintf("%v", hc.config.Publisher))

	// Perform request
	resp, err := client.Do(req)
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log.Printf("Response Status %v: %v\n", resp.StatusCode, resp.Status)
		return false
	}
	return true
}

func (hc *HashrefClient) SetRemoteData(inputType HashType, input string, calculatedHash string, metadata map[string]interface{}) bool {
	log.Printf("Set data for hash %v\n", calculatedHash)
	jsonData, err := json.Marshal(metadata)

	// prepare get request
	targetApi := "hash"
	if inputType == Publisher {
		targetApi = "publisher"
	}
	requestUri := fmt.Sprintf("%v/api/%v/%v", hc.config.HashrefServer, targetApi, calculatedHash)
	req, err := http.NewRequest(http.MethodPost, requestUri, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return false
	}

	// Set request meta data
	req.Header.Set("Authorization", fmt.Sprintf("%v", hc.config.Publisher))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return false
	}

	// Create Client
	client := &http.Client{}

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return false
	}

	// Get response
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return false
	}
	if resp.StatusCode >= 400 {
		log.Printf("Response Status %v: %v\n", resp.StatusCode, resp.Status)
		return false
	}
	return true
}

// SetSelf sets the metadata to the hash of the identity remoely
func (hc *HashrefClient) SetSelf(metadata map[string]interface{}) bool {
	log.Printf("Set data for yourself %v\n", hc.config.Publisher)

	// Transform json data
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return false
	}

	log.Printf("%s\n", jsonData)

	// Create Client
	client := &http.Client{}

	// Prepare request obj
	requestUri := fmt.Sprintf("%v/api/self", hc.config.HashrefServer)
	req, err := http.NewRequest(http.MethodPost, requestUri, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return false
	}
	req.Header.Set("Authorization", fmt.Sprintf("%v", hc.config.Publisher))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return false
	}

	// Post process response
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return false
	}
	log.Printf("Response Status %v: %v\n", resp.StatusCode, resp.Status)
	return resp.StatusCode < 400
}

// GetSelf performs a request to the server and collects the metadata
// that is stored remotly to the publisher
func (hc *HashrefClient) GetSelf() map[string]interface{} {

	// prepare get request
	requestUri := fmt.Sprintf("%v/api/self", hc.config.HashrefServer)
	log.Printf("Request-uri: %v\n", requestUri)
	req, err := http.NewRequest(http.MethodGet, requestUri, nil)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return map[string]interface{}{}
	}
	req.Header.Add("Authorization", fmt.Sprintf("%v", hc.config.Publisher))

	// Create Client
	client := &http.Client{}

	// Request own data
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{"Error": err.Error()}
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		log.Printf("Response Status %v: %v\n", resp.StatusCode, resp.Status)
		return map[string]interface{}{
			"status": resp.Status,
			"code":   resp.StatusCode,
		}
	}

	// Transform remote data into map
	remoteData := make(map[string]interface{})
	if err := json.Unmarshal(body, &remoteData); err != nil {
		return map[string]interface{}{
			"Error": err.Error(),
		}
	}
	return remoteData
}

// CollectLocalMetadata collects based on the provided input
// metadata and returns it as a map with metadata
func (hc *HashrefClient) CollectLocalMetadata(inputType HashType, input, hash string) map[string]interface{} {

	// Default values
	retMap := map[string]interface{}{
		"input":          input,
		"type":           Lookup[inputType],
		"last_published": time.Now().String(),
	}

	// Collect further data depended on type
	log.Printf("Collect metadata for %v (%v)\n", input, Lookup[inputType])
	switch inputType {

	// Get meta to text
	case Text:
		retMap["length"] = fmt.Sprint(len(input))

	// Get meta to file
	case File:
		fInfo, err := os.Stat(input)
		if err != nil {
			log.Printf("ERROR:\n%v", err)
			return make(map[string]interface{})
		}
		retMap["permission"] = fInfo.Mode().Perm().String()
		retMap["size"] = fmt.Sprintf("%v", fInfo.Size())
	}
	return retMap
}
