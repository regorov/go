package clipperz

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type LoginInfo struct{}

func ParseLoginInfo(input interface{}) (loginInfo *LoginInfo) {
	return nil
}

type Backend struct {
	URL            string
	ShouldPayTolls bool
	Tolls          map[string][]*Toll
	SRP            *SRPConnection
}

func NewBackend(urlStr string, shouldPayTolls bool, srp *SRPConnection) (backend *Backend) {
	return &Backend{
		URL:            urlStr,
		ShouldPayTolls: shouldPayTolls,
		Tolls:          make(map[string][]*Toll),
		SRP:            srp,
	}
}

func (backend *Backend) PayToll(requestType string, input interface{}) (output interface{}, err error) {
	if backend.ShouldPayTolls {
		tolls, ok := backend.Tolls[requestType]
		if !ok {
			output, err = backend.SendMessage("knock", map[string]interface{}{"requestType": requestType})
			if err != nil {
				return nil, err
			}

			_, err = backend.SetToll(output.(map[string]interface{}))
			if err != nil {
				return nil, err
			}

			tolls = backend.Tolls[requestType]
		}

		toll := tolls[0]
		backend.Tolls[requestType] = tolls[1:]

		err = toll.Pay()
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"parameters": input,
			"toll":       toll.RequestData(),
		}, nil
	}

	return map[string]interface{}{
		"parameters": input,
	}, nil
}

func (backend *Backend) SendMessage(method string, input interface{}) (output map[string]interface{}, err error) {
	fmt.Printf(">>> %s %v\n", method, input)

	dataEncoded, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	params := make(url.Values)
	params.Set("method", method)
	params.Set("parameters", string(dataEncoded))
	paramsEncoded := params.Encode()

	req, err := http.NewRequest("POST", backend.URL, strings.NewReader(paramsEncoded))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.FormatInt(int64(len(paramsEncoded)), 10))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request returned %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return nil, err
	}

	result, ok := output["result"].(string)
	if ok && result == "EXCEPTION" {
		return nil, fmt.Errorf(output["message"].(string))
	}

	fmt.Printf("<<< %v\n", result)

	return output, nil
}

func (backend *Backend) SetToll(input map[string]interface{}) (output interface{}, err error) {
	tollData, ok := input["toll"]
	if ok {
		tollMap := tollData.(map[string]interface{})
		requestType := tollMap["requestType"].(string)
		targetValueHex := tollMap["targetValue"].(string)
		cost := int(tollMap["cost"].(float64))

		targetValue, err := hex.DecodeString(targetValueHex)
		if err != nil {
			return nil, err
		}

		toll := NewToll(requestType, targetValue, cost)
		tolls := backend.Tolls[requestType]
		backend.Tolls[requestType] = append(tolls, toll)
	}

	return input["result"], nil
}

func (backend *Backend) ProcessMessage(method string, parameters interface{}, requestType string) (response interface{}, err error) {
	parameters, err = backend.PayToll(requestType, parameters)
	if err != nil {
		return nil, err
	}

	responseMap, err := backend.SendMessage(method, parameters)
	if err != nil {
		return nil, err
	}

	return backend.SetToll(responseMap)
}

func (backend *Backend) Registration(parameters interface{}) (response interface{}, err error) {
	return backend.ProcessMessage("registration", parameters, "REGISTER")
}

func (backend *Backend) Handshake(parameters interface{}) (response interface{}, err error) {
	return backend.ProcessMessage("handshake", parameters, "CONNECT")
}

func (backend *Backend) Message(parameters interface{}) (response interface{}, err error) {
	return backend.ProcessMessage("message", parameters, "MESSAGE")
}

func (backend *Backend) Logout(parameters interface{}) (response interface{}, err error) {
	return backend.ProcessMessage("logout", parameters, "MESSAGE")
}

func (backend *Backend) Login() (connectionID string, loginInfo *LoginInfo, offlineCopyNeeded bool, lock string, err error) {
	args := map[string]interface{}{
		"message": "connect",
		"version": "2.0",
		"parameters": map[string]interface{}{
			"C": hex.EncodeToString(backend.SRP.C),
			"A": fmt.Sprintf("%X", backend.SRP.A),
		},
	}

	response, err := backend.Handshake(args)
	if err != nil {
		return "", nil, false, "", err
	}

	responseMap := response.(map[string]interface{})
	s, ok := big.NewInt(0).SetString(responseMap["s"].(string), 16)
	if !ok {
		return "", nil, false, "", fmt.Errorf("Could not initialise SRP s")
	}

	B, ok := big.NewInt(0).SetString(responseMap["B"].(string), 16)
	if !ok {
		return "", nil, false, "", fmt.Errorf("Could not initialise SRP B")
	}

	backend.SRP.SetResponseData(s, B)

	// -------------------------------------------------------------------------

	args = map[string]interface{}{
		"message": "credentialCheck",
		"version": "2.0",
		"parameters": map[string]interface{}{
			"M1": backend.SRP.M1,
		},
	}

	response, err = backend.Handshake(args)
	if err != nil {
		return "", nil, false, "", err
	}

	responseMap = response.(map[string]interface{})
	M2 := responseMap["M2"].(string)

	if M2 != backend.SRP.M2 {
		return "", nil, false, "", fmt.Errorf("Checksum mismatch")
	}

	// -------------------------------------------------------------------------

	connectionID = responseMap["connectionId"].(string)
	loginInfo = ParseLoginInfo(responseMap["loginInfo"])
	offlineCopyNeeded = responseMap["offlineCopyNeeded"].(bool)
	lock = responseMap["lock"].(string)

	return connectionID, loginInfo, offlineCopyNeeded, lock, nil
}
