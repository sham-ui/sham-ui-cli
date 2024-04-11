package asset

import (
	"cms/internal/controller/grpc/proto"
	"cms/internal/model"
	"cms/pkg/asserts"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestHandler_Asset(t *testing.T) {
	filePath := "/test.txt"
	content := []byte("test")
	testCases := []struct {
		name             string
		ctx              context.Context
		request          *proto.AssetRequest
		assetService     func(t mockConstructorTestingTNewMockAssetService) *MockAssetService
		expectedResponse *proto.AssetResponse
		expectedError    error
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			request: &proto.AssetRequest{Path: filePath},
			assetService: func(t mockConstructorTestingTNewMockAssetService) *MockAssetService {
				m := NewMockAssetService(t)
				m.
					On("ReadFile", mock.Anything, filePath).
					Return(content, nil).
					Once()
				return m
			},
			expectedResponse: &proto.AssetResponse{
				Response: &proto.AssetResponse_File{
					File: content,
				},
			},
		},
		{
			name: "clean file path",
			ctx:  context.Background(),
			request: &proto.AssetRequest{
				Path: "\\TEST.TXT",
			},
			assetService: func(t mockConstructorTestingTNewMockAssetService) *MockAssetService {
				m := NewMockAssetService(t)
				m.
					On("ReadFile", mock.Anything, filePath).
					Return(content, nil).
					Once()
				return m
			},
			expectedResponse: &proto.AssetResponse{
				Response: &proto.AssetResponse_File{
					File: content,
				},
			},
		},
		{
			name:    "file not found",
			ctx:     context.Background(),
			request: &proto.AssetRequest{Path: filePath},
			assetService: func(t mockConstructorTestingTNewMockAssetService) *MockAssetService {
				m := NewMockAssetService(t)
				m.
					On("ReadFile", mock.Anything, filePath).
					Return(nil, model.NewAssetNotFoundError(filePath)).
					Once()
				return m
			},
			expectedResponse: &proto.AssetResponse{
				Response: &proto.AssetResponse_NotFound{
					NotFound: &proto.NotFound{},
				},
			},
		},
		{
			name:    "fail read file",
			ctx:     context.Background(),
			request: &proto.AssetRequest{Path: filePath},
			assetService: func(t mockConstructorTestingTNewMockAssetService) *MockAssetService {
				m := NewMockAssetService(t)
				m.
					On("ReadFile", mock.Anything, filePath).
					Return(nil, errors.New("test")).
					Once()
				return m
			},
			expectedError: errors.New("read file: test"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			h := NewHandler(test.assetService(t))

			// Action
			res, err := h.Asset(test.ctx, test.request)

			// Assert
			asserts.Equals(t, test.expectedResponse, res)
			asserts.ErrorsEqual(t, test.expectedError, err)
		})
	}
}
