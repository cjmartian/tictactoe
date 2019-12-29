package main

import (
    "log"
    "net/http"
    "encoding/json"
    "strconv"

    "github.com/gorilla/mux"
)

type server struct{}

var board [3][3]string

func getBoard(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    encoded, err := json.Marshal(board)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{"error": "error marshalling board data"}`))
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(encoded)
}

func resetBoard(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    board = [3][3]string{}
    encoded, err := json.Marshal(board)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{"error": "error marshalling board data"}`))
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(encoded)
}

func handleMove(w http.ResponseWriter, r *http.Request) {
    queries := mux.Vars(r)
    w.Header().Set("Content-Type", "application/json")
    var err error
    player := queries["player"]

    y := -1
    if val, ok := queries["y"]; ok {
        y, err = strconv.Atoi(val)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(`{"error": "error converting string data to int"}`))
            return
        }
    }

    x := -1
    if val, ok := queries["x"]; ok {
        x, err = strconv.Atoi(val)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte(`{"error": "error converting string data to int"}`))
            return
        }
    }

    if len(board[y][x]) == 0 {
       board[y][x] = player
       encoded, err := json.Marshal(board)
       if err != nil {
           w.WriteHeader(http.StatusInternalServerError)
           w.Write([]byte(`{"error": "error marshalling board data"}`))
           return
       }
       w.WriteHeader(http.StatusOK)
       w.Write(encoded)
       return
    }
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(`{"error": "piece already exists at that location"}`))
}

func main() {
    r := mux.NewRouter()

    api := r.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/board", getBoard).Methods(http.MethodGet)
    api.HandleFunc("/board/{x}/{y}/{player}", handleMove).Methods(http.MethodPost)
    api.HandleFunc("/board/reset", resetBoard).Methods(http.MethodGet)

    log.Fatal(http.ListenAndServe(":8080", r))
}
