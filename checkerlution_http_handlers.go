package checkerlution

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterHandlers(trainer *CheckerlutionTrainer) {

	r := mux.NewRouter()

	cortexHandler := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cortexUuid := vars["cortex_uuid"]
		cortex := trainer.population.Find(cortexUuid)
		fmt.Fprintf(w, "%v", cortex)
	}

	cortexSvgHandler := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		cortexUuid := vars["cortex_uuid"]
		cortex := trainer.population.Find(cortexUuid)
		cortex.RenderSVG(w)
	}

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/cortex/{cortex_uuid}", cortexHandler)
	r.HandleFunc("/cortex/{cortex_uuid}/svg", cortexSvgHandler)
	http.Handle("/", r)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Endpoints: /cortex/{cortex_uuid} and /cortex/{cortex_uuid}/svg")
}
