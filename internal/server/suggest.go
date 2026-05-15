package server

import (
    "encoding/json"
    "net/http"
    "github.com/Sylvester-Kapoko/risitPap/internal/store"
)

func HandleSuggest(st *store.Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        q := r.URL.Query().Get("q")
        if len(q) < 2 {
            w.Write([]byte("[]"))
            return
        }
        rows, err := st.SuggestItems(q)
        if err != nil {
            http.Error(w, "internal error", 500)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(rows)
    }
}