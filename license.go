package bytebuilders

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func VerifyLicense(token string) (*License, error) {
	url := "https://c7c341b6.ngrok.io/api/v1/user/licenses/verify"
	licenseToken := LicenseVerificationParams{
		Raw: token,
	}
	jsonBytes, err := json.Marshal(licenseToken)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "JWT "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	license := new(License)
	if err = json.NewDecoder(resp.Body).Decode(license); err != nil {
		return nil, err
	}

	return license, nil
}
