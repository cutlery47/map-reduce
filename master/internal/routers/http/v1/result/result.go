package result

import (
	"net/http"

	"github.com/cutlery47/map-reduce/master/internal/domain"
)

type resultRoutes struct {
	srv *domain.Master
}

func (rr *resultRoutes) mapRes(w http.ResponseWriter, r *http.Request) {

}

func (rr *resultRoutes) redRes(w http.ResponseWriter, r *http.Request) {

}
