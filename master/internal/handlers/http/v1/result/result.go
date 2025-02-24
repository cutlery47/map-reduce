package result

import (
	"net/http"

	"github.com/cutlery47/map-reduce/master/internal/service"
)

type resultRoutes struct {
	srv *service.Master
}

func (rr *resultRoutes) mapRes(w http.ResponseWriter, r *http.Request) {

}

func (rr *resultRoutes) redRes(w http.ResponseWriter, r *http.Request) {

}
