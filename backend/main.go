package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ResultRequest struct {
	Exam  string `json:"exam"`
	Year  string `json:"year"`
	Board string `json:"board"`
	Roll  string `json:"roll"`
	Reg   string `json:"reg"`
}

type Subject struct {
	Code    string `json:"code"`
	Subject string `json:"subject"`
	Grade   string `json:"grade"`
}

type ResultResponse struct {
	Demo   bool      `json:"demo"`
	Name   string    `json:"name"`
	Roll   string    `json:"roll"`
	Board  string    `json:"board"`
	Exam   string    `json:"exam"`
	Year   string    `json:"year"`
	GPA    string    `json:"gpa"`
	Result string    `json:"result"`
	Grades []Subject `json:"grades"`
}

type Visitor struct {
	Count int
	Reset time.Time
}

var (
	visitors = make(map[string]*Visitor)
	mu       sync.Mutex
	digits   = regexp.MustCompile(`^[0-9]{4,20}$`)
)

func sendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func rateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			sendError(w, http.StatusBadRequest, "Invalid connection")
			return
		}
		now := time.Now()

		mu.Lock()

		v, exists := visitors[ip]

		if !exists || now.After(v.Reset) {
			visitors[ip] = &Visitor{
				Count: 1,
				Reset: now.Add(time.Minute),
			}
		} else {
			v.Count++

			if v.Count > 30 {
				mu.Unlock()
				sendError(w, http.StatusTooManyRequests,
					"Too many requests. Try again later.")
				return
			}
		}

		mu.Unlock()

		next(w, r)
	}
}

func validate(req ResultRequest) error {
	if req.Exam != "hsc" && req.Exam != "ssc" {
		return errors.New("invalid exam")
	}

	if req.Year != "2024" &&
		req.Year != "2025" &&
		req.Year != "2026" {
		return errors.New("invalid year")
	}

	boards := map[string]bool{
		"dhaka":      true,
		"rajshahi":   true,
		"comilla":    true,
		"chittagong": true,
		"jessore":    true,
		"barisal":    true,
		"sylhet":     true,
		"dinajpur":   true,
		"mymensingh": true,
		"madrasah":   true,
		"technical":  true,
	}

	if !boards[req.Board] {
		return errors.New("invalid board")
	}

	if !digits.MatchString(req.Roll) ||
		!digits.MatchString(req.Reg) {
		return errors.New("invalid roll or registration")
	}

	return nil
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	origin := r.Header.Get("Origin")

	if origin == "http://localhost:5500" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var req ResultRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	req.Exam = strings.ToLower(strings.TrimSpace(req.Exam))
	req.Board = strings.ToLower(strings.TrimSpace(req.Board))

	if err := validate(req); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := ResultResponse{
		Demo:   true,
		Name:   "DEMO STUDENT",
		Roll:   req.Roll,
		Board:  req.Board,
		Exam:   req.Exam,
		Year:   req.Year,
		GPA:    "4.50",
		Result: "SAMPLE ONLY",
		Grades: []Subject{
			{"101", "Bangla", "A"},
			{"107", "English", "A-"},
			{"275", "ICT", "A+"},
			{"121", "Logic", "A"},
		},
	}

	json.NewEncoder(w).Encode(result)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/result", rateLimit(resultHandler))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	log.Println("API running on port", port)

	log.Fatal(server.ListenAndServe())
}
