package model

type Asset struct {
	Path    string
	Content []byte
}

type AssetNotFoundError struct {
	Path string
}

func (e AssetNotFoundError) Error() string {
	return "asset not found: " + e.Path
}

func (e AssetNotFoundError) Is(target error) bool {
	//nolint:errorlint
	_, ok := target.(AssetNotFoundError)
	return ok
}

func NewAssetNotFoundError(path string) error {
	return AssetNotFoundError{
		Path: path,
	}
}
