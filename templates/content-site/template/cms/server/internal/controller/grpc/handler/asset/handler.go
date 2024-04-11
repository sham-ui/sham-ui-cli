package asset

import (
	"cms/internal/controller/grpc/proto"
	"cms/internal/model"
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
)

type handler struct {
	assetService AssetService
}

func (h *handler) Asset(ctx context.Context, request *proto.AssetRequest) (*proto.AssetResponse, error) {
	filePath := strings.ToLower(request.Path)
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	filePath = path.Clean(filePath)
	content, err := h.assetService.ReadFile(ctx, filePath)
	switch {
	case errors.Is(err, model.AssetNotFoundError{}):
		return &proto.AssetResponse{
			Response: &proto.AssetResponse_NotFound{
				NotFound: &proto.NotFound{},
			},
		}, nil
	case err != nil:
		return nil, fmt.Errorf("read file: %w", err)
	}
	return &proto.AssetResponse{
		Response: &proto.AssetResponse_File{
			File: content,
		},
	}, nil
}

func NewHandler(assetService AssetService) *handler {
	return &handler{assetService: assetService}
}
