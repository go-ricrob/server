package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-ricrob/exec/task"
	"github.com/go-ricrob/game/board"
	"github.com/go-ricrob/server/internal/exec"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, httpError int, err error) {
	w.WriteHeader(httpError)
	errorResponse := &errorResponse{Error: err.Error()}
	if b, jsonErr := json.Marshal(errorResponse); jsonErr == nil { // should always be the case
		w.Write(b)
	}
}

func boardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid http method %s", r.Method))
		return
	}

	tiles := new(task.Tiles)

	if err := tiles.ParseURL(r.URL); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	board := board.New([board.NumTile]string{
		board.TopLeft:     tiles.TopLeft,
		board.TopRight:    tiles.TopRight,
		board.BottomLeft:  tiles.BottomLeft,
		board.BottomRight: tiles.BottomRight,
	})

	response, err := json.Marshal(board.Fields)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}
	w.Write(response)
}

type solveHandler struct {
	execer *exec.Execer
}

func (h *solveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	args, err := task.ParseURL(r.URL)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	resultCh := h.execer.Run(args.CmdArgs())

	for result := range resultCh {
		if result.Err != nil {
			writeError(w, http.StatusInternalServerError, result.Err)
		} else {
			w.Write(result.Response)
		}
		flusher.Flush()
	}
}
