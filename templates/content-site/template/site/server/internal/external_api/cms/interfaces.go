package cms

import "site/internal/external_api/cms/proto"

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name cmsClient --inpackage --testonly
type cmsClient interface {
	proto.CMSClient
}
