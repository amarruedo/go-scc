package scc

import (
	"context"
	"errors"
	"mime"
	"os"
	"path/filepath"
	"strconv"
)

type BackupService service

// SCC API docs https://help.sap.com/viewer/cca91383641e40ffbe03bdc78f00f681/Cloud/en-US/d94b9db4320c4392bcb15accef64d369.html

// CreateBackup creates a backup configuration with 'password' as the password used for encrypting sensible data.
// Only sensitive data in the backup are encrypted with an arbitrary password of your choice. The password is required for the restore operation. The returned ZIP archive itself is not password-protected
func (s *BackupService) CreateBackup(ctx context.Context, password string, file *os.File) (*Response, error) {
	req, err := s.client.NewRequest("POST", "api/v1/configuration/backup", struct {
		Password string `json:"password"`
	}{Password: password})
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, file)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RestoreBackup restores a backup configuration
func (s *BackupService) RestoreBackup(ctx context.Context, password string, file *os.File) (*Response, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, errors.New("the asset to upload can't be a directory")
	}

	mediaType := mime.TypeByExtension(filepath.Ext(file.Name()))
	req, err := s.client.NewUploadRequest("PUT", "api/v1/configuration/backup", file, stat.Size(), mediaType)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != 204 {
		return resp, errors.New("backup restore failed with status code " + strconv.Itoa(resp.StatusCode))
	}

	return resp, nil
}
