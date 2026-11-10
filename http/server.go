package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/demolemo/meth"
	"github.com/demolemo/meth/sqlite"
)

type Server struct {
	server *http.Server

	ms sqlite.MetricService
}

func NewServer(db *sql.DB) *Server {
	s := &Server{
		&http.Server{
			Addr: ":8082",
		},
		*sqlite.NewMetricService(db),
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /metrics", s.handleMetricList)
	mux.HandleFunc("POST /metrics", s.handleMetricCreate)
	mux.HandleFunc("GET /metrics/{name}", s.handleMetricGet)
	mux.HandleFunc("DELETE /metrics/{name}", s.handleMetricDelete)
	s.server.Handler = mux
}

func (s *Server) Serve() error {
	return s.server.ListenAndServe()
}

func (s *Server) handleMetricList(w http.ResponseWriter, r *http.Request) {
	limit := -1 // all
	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "limit must be an integer", http.StatusBadRequest)
			return
		}
		limit = n
	}

	metrics, err := s.ms.GetMetrics(limit)
	if err != nil {
		http.Error(w, "Getting metrics failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, "Metrics encoding failed: "+err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Server) handleMetricCreate(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name string `json:"name"`
	}
	// cap json size at 1MiB
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // cap body at 1 MiB
	dec := json.NewDecoder(r.Body)
	// a typo like {"nmae":...} becomes a 400, not a metric with an empty name
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
		return
	}

	// create metric pointer and pass it to creation method
	met := &meth.Metric{Name: in.Name}
	err := s.ms.CreateMetric(met)
	if err != nil {
		http.Error(w, "metric creation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(fmt.Appendf(nil, "Metric created, id: %d, name: %s\n", met.ID, met.Name))
}

func (s *Server) handleMetricGet(w http.ResponseWriter, r *http.Request) {
	// decode metric id
	name := r.PathValue("name")

	// fetch metric from the db
	metric, err := s.ms.GetMetric(name)
	if err != nil {
		// TODO: throw server error
		// Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	// switch r.Header.Get("Accept") {
	// case "application/json":
	// 	w.Header().Set("Content-type", "application/json")
	// 	if err := json.NewEncoder(w).Encode(metric); err != nil {
	// 		// TODO: error
	// 		// LogError(r, err)
	// 		return
	// 	}

	// default:
	// 	return
	// 	// tmpl := html.DialViewTemplate{
	// 	// 	Dial:      dial,
	// 	// 	InviteURL: fmt.Sprintf("%s/invite/%s", s.URL(), dial.InviteCode),
	// 	// }
	// 	// tmpl.Render(r.Context(), w)
	// }

	if err := json.NewEncoder(w).Encode(metric); err != nil {
		// TODO: throw server error
		return
	}
}

func (s *Server) handleMetricDelete(w http.ResponseWriter, r *http.Request) {
	// get metric name from URL
	name := r.PathValue("name")

	// delete metric from the db
	err := s.ms.DeleteMetric(name)
	if err != nil {
		http.Error(w, "metric deletion failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
