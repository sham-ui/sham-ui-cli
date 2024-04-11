package logout

import "net/http"

//go:generate go run github.com/vektra/mockery/v2@v2.20.0 --name Service --inpackage --testonly
type Service interface {
	Delete(rw http.ResponseWriter, r *http.Request) error
}
