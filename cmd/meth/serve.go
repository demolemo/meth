// TODO: consolidate this functionality into the methd part
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/demolemo/meth"
	"github.com/demolemo/meth/parse"

	"github.com/spf13/cobra"
)

// a Health Export of several verbose years is a few MB, so this is room to
// spare and still a bound: without one, any request body fills memory
const maxImportBytes = 32 << 20

var serveAddr string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the HTTP API over the database",
	Long: `Serve the HTTP API over the database.

There is no authentication, so the default address is loopback only. Anything
that exposes this to a network has to put a reverse proxy in front of it.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		srv := &http.Server{
			Addr:              serveAddr,
			Handler:           routes(),
			ReadHeaderTimeout: 5 * time.Second,
		}

		errc := make(chan error, 1)
		go func() {
			fmt.Printf("meth listening on http://%s\n", srv.Addr)
			errc <- srv.ListenAndServe()
		}()

		select {
		case err := <-errc:
			return err
		case <-ctx.Done():
			// the deferred db.Close in root runs after this returns, so the
			// handlers have to be done with it first
			shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			fmt.Println("\nshutting down")
			return srv.Shutdown(shutdown)
		}
	},
}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /import", handleImport)
	return mux
}

// handleImport takes a raw Health Export Kit file as the body and stores its
// bodyMass readings as values of ?metric=. The conversion lives here and not
// on the device that posts it, so the only thing a client has to know is a URL.
func handleImport(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("metric")
	if name == "" {
		writeError(w, http.StatusBadRequest, errors.New("metric query parameter is required"))
		return
	}

	m, err := ms.GetMetric(name)
	if err != nil {
		writeError(w, statusFor(err), err)
		return
	}

	mvs, err := parse.HealthExportWeightFrom(http.MaxBytesReader(w, r.Body, maxImportBytes), "request body", m.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	res, err := vs.CreateMetricValues(mvs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("import %s: %d new, %d updated, %d unchanged", name, res.Inserted, res.Updated, res.Unchanged)
	writeJSON(w, http.StatusOK, importResponse{
		Metric:    name,
		Inserted:  res.Inserted,
		Updated:   res.Updated,
		Unchanged: res.Unchanged,
	})
}

// the wire shape is this package's business, not the domain's, which is why
// meth.ImportResult is copied into it rather than tagged and marshalled
type importResponse struct {
	Metric    string `json:"metric"`
	Inserted  int    `json:"inserted"`
	Updated   int    `json:"updated"`
	Unchanged int    `json:"unchanged"`
}

func statusFor(err error) int {
	if errors.Is(err, meth.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("writing response: %v", err)
	}
}

// {"error": "..."} is the shape appd already returns on the box, so one client
// can read a failure from either
func writeError(w http.ResponseWriter, code int, err error) {
	log.Printf("%d: %v", code, err)
	writeJSON(w, code, map[string]string{"error": err.Error()})
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&serveAddr, "addr", "127.0.0.1:8080", "address to listen on")
}
