package scc

import (
	"context"
)

type CommonService service

type CommonProperties struct {
	Ha struct {
		Role string `json:"role"`
	} `json:"ha"`
	Description string `json:"description"`
}

type Version struct {
	Version string `json:"version"`
}

// SCC API docs: https://help.sap.com/viewer/cca91383641e40ffbe03bdc78f00f681/Cloud/en-US/8aed644c6bce4bcabfd31d1c74e8fec4.html

// GetCommonProperties gets common properties of Cloud Connector
func (s *CommonService) GetCommonProperties(ctx context.Context) (*CommonProperties, *Response, error) {
	req, err := s.client.NewRequest("GET", "api/v1/configuration/connector", nil)
	if err != nil {
		return nil, nil, err
	}

	commonProperties := new(CommonProperties)
	resp, err := s.client.Do(ctx, req, commonProperties)
	if err != nil {
		return nil, resp, err
	}

	return commonProperties, resp, nil
}

// GetVersion gets the version of Cloud Connector
func (s *CommonService) GetVersion(ctx context.Context) (*Version, *Response, error) {
	req, err := s.client.NewRequest("GET", "api/v1/connector/version", nil)
	if err != nil {
		return nil, nil, err
	}

	version := new(Version)
	resp, err := s.client.Do(ctx, req, version)
	if err != nil {
		return nil, resp, err
	}

	return version, resp, nil
}

// SetDescription sets the description of Cloud Connector
func (s *CommonService) SetDescription(ctx context.Context, description string) (*CommonProperties, *Response, error) {
	req, err := s.client.NewRequest("PUT", "api/v1/configuration/connector", struct {
		Description string `json:"description"`
	}{Description: description})
	if err != nil {
		return nil, nil, err
	}

	commonProperties := new(CommonProperties)
	resp, err := s.client.Do(ctx, req, commonProperties)
	if err != nil {
		return nil, resp, err
	}

	return commonProperties, resp, nil
}
