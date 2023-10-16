package server

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	Log *Log
}

func NewServer() *Server {
	return &Server{
		Log: NewLog(),
	}
}

type ReadRequest struct {
	Offset uint64 `json:"offset"`
}

type ReadResponse struct {
	Record Record `json:"record"`
}

type WriteRequest struct {
	Record Record `json:"record"`
}

type WriteResponse struct {
	Offset uint64 `json:"offset"`
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	// If the request URL path does not exactly match "/", then send
	// a 404: Not Found response to the client.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Home"))
}

func (s *Server) logRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// Let the client know which request method is supported.
		w.Header().Set("Allow", "GET")

		// Use the http.Error() function to send a 405 status code and
		// "Method Not Allowed" string as the response body.
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReadRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ReadResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *Server) logWrite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Let the client know which request methods are supported.
		w.Header().Set("Allow", "POST")

		// Use the http.Error() function to send a 405 status code and
		// "Method Not Allowed" string as the response body.
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WriteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	offset, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := WriteResponse{Offset: offset}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func CreateServer(addr string) *http.Server {
	srv := NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.home)
	mux.HandleFunc("/log/read", srv.logRead)
	mux.HandleFunc("/log/write", srv.logWrite)

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
