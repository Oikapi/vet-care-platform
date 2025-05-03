package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vet-app-appointments/internal/models"
)

type UserServiceClient interface {
	GetDoctorByID(id uint) (*models.Doctor, error)
	GetClinicByID(id uint) (*models.Clinic, error)
	GetClientByID(id uint) (*models.User, error)
}

type HTTPUserServiceClient struct {
	BaseURL string
}

func NewHTTPUserServiceClient(baseURL string) UserServiceClient {
	return &HTTPUserServiceClient{BaseURL: baseURL}
}

func (c *HTTPUserServiceClient) GetDoctorByID(id uint) (*models.Doctor, error) {
	resp, err := http.Get(fmt.Sprintf("%s/doctors/%d", c.BaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("doctor service returned status %d", resp.StatusCode)
	}

	var doctor models.Doctor
	if err := json.NewDecoder(resp.Body).Decode(&doctor); err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (c *HTTPUserServiceClient) GetClinicByID(id uint) (*models.Clinic, error) {
	url := fmt.Sprintf("%s/clinics/%d", c.BaseURL, id)
	fmt.Println("Requesting:", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("clinic service returned status %d", resp.StatusCode)
	}

	var clinic models.Clinic
	if err := json.NewDecoder(resp.Body).Decode(&clinic); err != nil {
		return nil, err
	}
	return &clinic, nil
}

func (c *HTTPUserServiceClient) GetClientByID(id uint) (*models.User, error) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%d", c.BaseURL, id))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client service returned status %d", resp.StatusCode)
	}

	var client models.User
	if err := json.NewDecoder(resp.Body).Decode(&client); err != nil {
		return nil, err
	}
	return &client, nil
}
