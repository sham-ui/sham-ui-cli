package asset

import (
	"cms/internal/model"
	"cms/pkg/tracing"
	"context"
	"crypto/md5" //nolint:gosec
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const scopeName = "service.asset"

type service struct {
	uploadPath string
}

func (srv *service) Save(ctx context.Context, filename string, content []byte) (*model.Asset, error) {
	const op = "Save"

	_, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	destFilename := srv.buildFileName(ctx, filename, content)
	destPath := filepath.Join(srv.uploadPath, destFilename)
	fileExists, err := srv.checkFileExists(ctx, destPath)
	switch {
	case err != nil:
		return nil, fmt.Errorf("check file exists: %w", err)
	case fileExists:
		return &model.Asset{
			Path:    destFilename,
			Content: content,
		}, nil
	}
	fo, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer fo.Close()
	if _, err := fo.Write(content); err != nil {
		return nil, fmt.Errorf("write file: %w", err)
	}
	return &model.Asset{
		Path:    destFilename,
		Content: content,
	}, nil
}

func (srv *service) ReadFile(ctx context.Context, path string) ([]byte, error) {
	const op = "ReadFile"

	_, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	destPath := filepath.Join(srv.uploadPath, path)

	exists, err := srv.checkFileExists(ctx, destPath)
	switch {
	case err != nil:
		return nil, fmt.Errorf("check file exists: %w", err)
	case !exists:
		return nil, model.NewAssetNotFoundError(destPath)
	}
	content, err := os.ReadFile(destPath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return content, nil
}

func (srv *service) buildFileName(ctx context.Context, filename string, content []byte) string {
	const op = "buildFileName"

	_, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	hash := md5.Sum(content) //nolint:gosec
	return hex.EncodeToString(hash[:]) + strings.ToLower(filepath.Ext(filename))
}

func (srv *service) checkFileExists(ctx context.Context, destPath string) (bool, error) {
	const op = "checkFileExists"

	_, span := tracing.StartSpan(ctx, scopeName, op)
	defer span.End()

	_, err := os.Stat(destPath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("stat: %w", err)
	default:
		return true, nil
	}
}

func New(uploadPath string) *service {
	return &service{
		uploadPath: uploadPath,
	}
}
