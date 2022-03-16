package scc

import (
	"bytes"
	"context"
	"errors"
	"strconv"
)

type HAService service

// SCC API docs: https://help.sap.com/viewer/cca91383641e40ffbe03bdc78f00f681/Cloud/en-US/8aed644c6bce4bcabfd31d1c74e8fec4.html

// GetHASettings gets common properties of Cloud Connector.
func (s *HAService) GetHASettings(ctx context.Context) (string, *Response, error) {
	req, err := s.client.NewRequest("GET", "api/v1/configuration/connector/haRole", nil)
	if err != nil {
		return "", nil, err
	}

	var b bytes.Buffer
	resp, err := s.client.Do(ctx, req, &b)
	if err != nil {
		return "", resp, err
	}

	return b.String(), resp, nil
}

// SetHASettings sets the role of a fresh installation.
// As of version 2.12.0, this API also allows to switch the roles if a shadow instance is connected to the master.
// Role mast be "master" or "shadow".
func (s *HAService) SetHASettings(ctx context.Context, role string) (*Response, error) {
	req, err := s.client.NewRequest("POST", "api/v1/configuration/connector/haRole", role)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

type MasterConfiguration struct {
	HAEnabled         bool   `json:"haEnabled"`
	AllowedShadowHost string `json:"allowedShadowHost"`
}

// GetMasterConfiguration gets master instance configuration:
//  ●haEnabled: a Boolean value that indicates whether or not a shadow system is allowed to connect.
//  ●allowedShadowHost: the name of the shadow host (a string) that is allowed to connect; an empty string signifies that any host is allowed to connect as shadow.
func (s *HAService) GetMasterConfiguration(ctx context.Context) (*MasterConfiguration, *Response, error) {
	req, err := s.client.NewRequest("GET", "api/v1/configuration/connector/ha/master/config", nil)
	if err != nil {
		return nil, nil, err
	}

	masterConfig := new(MasterConfiguration)
	resp, err := s.client.Do(ctx, req, masterConfig)
	if err != nil {
		return nil, resp, err
	}

	return masterConfig, resp, nil
}

// SetMasterConfiguration sets master instance configuration:
//  ●haEnabled: Boolean value that indicates whether or not a shadow system is allowed to connect.
//  ●allowedShadowHost: Name of the shadow host (a string) that is allowed to connect. An empty string means that any host is allowed to connect as shadow.
func (s *HAService) SetMasterConfiguration(ctx context.Context, haEnabled bool, allowedShadowHost string) (*MasterConfiguration, *Response, error) {
	req, err := s.client.NewRequest("PUT", "api/v1/configuration/connector/ha/master/config", MasterConfiguration{
		HAEnabled:         haEnabled,
		AllowedShadowHost: allowedShadowHost,
	})
	if err != nil {
		return nil, nil, err
	}

	masterConfig := new(MasterConfiguration)
	resp, err := s.client.Do(ctx, req, masterConfig)
	if err != nil {
		return nil, resp, err
	}

	return masterConfig, resp, nil
}

type MasterState struct {
	State      string `json:"state"`
	ShadowHost string `json:"shadowHost"`
}

// GetMasterState gets state of master instance:
//  ●state: One of the following strings: ALONE, BINDING, CONNECTED or BROKEN.
//  ●shadowHost: Connected shadow host (a string).
func (s *HAService) GetMasterState(ctx context.Context) (*MasterState, *Response, error) {
	req, err := s.client.NewRequest("GET", "api/v1/configuration/connector/ha/master/state", nil)
	if err != nil {
		return nil, nil, err
	}

	masterState := new(MasterState)
	resp, err := s.client.Do(ctx, req, masterState)
	if err != nil {
		return nil, resp, err
	}

	return masterState, resp, nil
}

// SetMasterState sets state of master instance where 'op' is one of the following strings:
//  ●SWITCH: Switch roles with shadow.
//  ●FORCE_SWITCH: Take over the shadow role, even if shadow instance does not respond.
func (s *HAService) SetMasterState(ctx context.Context, op string) (*MasterState, *Response, error) {
	req, err := s.client.NewRequest("POST", "api/v1/configuration/connector/ha/master/state", struct {
		Op string `json:"op"`
	}{Op: op})
	if err != nil {
		return nil, nil, err
	}

	masterState := new(MasterState)
	resp, err := s.client.Do(ctx, req, masterState)
	if err != nil {
		return nil, resp, err
	}

	return masterState, resp, nil
}

// ResetMaster restores default values for all settings related to high availability on the master side.
// Do not perform this call if the shadow is connected to a master.
func (s *HAService) ResetMaster(ctx context.Context) (*Response, error) {
	req, err := s.client.NewRequest("DELETE", "api/v1/configuration/connector/ha/master/state", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != 204 {
		return resp, errors.New("master reset failed with status code " + strconv.Itoa(resp.StatusCode))
	}

	return resp, nil
}

type ShadowConfiguration struct {
	MasterHost             string `json:"masterHost"`
	MasterPort             string `json:"masterPort"`
	CheckIntervalInSeconds int    `json:"checkIntervalInSeconds"`
	TakeoverDelayInSeconds int    `json:"takeoverDelayInSeconds"`
	OwnHost                string `json:"ownHost"`
	ConnectTimeoutInMillis int    `json:"connectTimeoutInMillis"`
	RequestTimeoutInMillis int    `json:"requestTimeoutInMillis"`
	Links                  struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		State struct {
			Href string `json:"href"`
		} `json:"state"`
	} `json:"_links,omitempty"`
}

// GetShadowConfiguration gets configuration settings for Cloud Connector shadow instance.
// The APIs below are only permitted on a Cloud Connector shadow instance. The master instance will reject the requests with error code 403 – FORBIDDEN_REQUEST.
func (s *HAService) GetShadowConfiguration(ctx context.Context) (*ShadowConfiguration, *Response, error) {
	req, err := s.client.NewShadowRequest("GET", "api/v1/configuration/connector/ha/shadow/config", nil)
	if err != nil {
		return nil, nil, err
	}

	shadowConfig := new(ShadowConfiguration)
	resp, err := s.client.Do(ctx, req, shadowConfig)
	if err != nil {
		return nil, resp, err
	}

	return shadowConfig, resp, nil
}

// SetShadowConfiguration sets configuration settings for Cloud Connector shadow instance.
// The APIs below are only permitted on a Cloud Connector shadow instance. The master instance will reject the requests with error code 403 – FORBIDDEN_REQUEST.
func (s *HAService) SetShadowConfiguration(ctx context.Context, masterHost, masterPort, ownHost string, checkIntervalInSeconds, takeoverDelayInSeconds, connectTimeoutInMillis, requestTimeoutInMillis int) (*ShadowConfiguration, *Response, error) {
	req, err := s.client.NewShadowRequest("PUT", "api/v1/configuration/connector/ha/shadow/config", ShadowConfiguration{
		MasterHost:             masterHost,
		MasterPort:             masterPort,
		CheckIntervalInSeconds: checkIntervalInSeconds,
		TakeoverDelayInSeconds: takeoverDelayInSeconds,
		OwnHost:                ownHost,
		ConnectTimeoutInMillis: connectTimeoutInMillis,
		RequestTimeoutInMillis: requestTimeoutInMillis,
	})
	if err != nil {
		return nil, nil, err
	}

	shadowConfigResp := new(ShadowConfiguration)
	resp, err := s.client.Do(ctx, req, shadowConfigResp)
	if err != nil {
		return nil, resp, err
	}

	return shadowConfigResp, resp, nil
}

type ShadowState struct {
	State          string `json:"state"`
	OwnHosts       string `json:"ownHosts"`
	StateMessage   string `json:"stateMessage"`
	MasterVersions string `json:"masterVersions"`
}

// GetShawodState gets state of shadow instance.
func (s *HAService) GetShawodState(ctx context.Context, description string) (*ShadowState, *Response, error) {
	req, err := s.client.NewShadowRequest("GET", "api/v1/configuration/connector/ha/shadow/state", nil)
	if err != nil {
		return nil, nil, err
	}

	shadowState := new(ShadowState)
	resp, err := s.client.Do(ctx, req, shadowState)
	if err != nil {
		return nil, resp, err
	}

	return shadowState, resp, nil
}

// ChangeShadowState changes state of shadow instance:
//  ●op: String value representing the state change operation. Possible values are CONNECT or DISCONNECT.
//  ●user: User for logon to the master instance.
//  ●password: Password for logon to the master instance.
func (s *HAService) ChangeShadowState(ctx context.Context, op, user, password string) (*Response, error) {
	req, err := s.client.NewShadowRequest("POST", "api/v1/configuration/connector/ha/shadow/state", struct {
		Op       string `json:"op"`
		User     string `json:"user"`
		Password string `json:"password"`
	}{Op: op, User: user, Password: password})
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != 204 {
		return resp, errors.New("shadow reset failed with status code " + strconv.Itoa(resp.StatusCode))
	}

	return resp, nil
}

// ResetShadow deletes master host and port, and restores default values for all other settings related to a connection to the master.
// Do not perform this call if the shadow is connected to a master.
// Available as of version 2.13.0.
func (s *HAService) ResetShadow(ctx context.Context) (*Response, error) {
	req, err := s.client.NewShadowRequest("DELETE", "api/v1/configuration/connector/ha/shadow/state", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != 204 {
		return resp, errors.New("shadow reset failed with status code " + strconv.Itoa(resp.StatusCode))
	}

	return resp, nil
}
