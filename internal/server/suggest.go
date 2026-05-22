package server

import (
	"encoding/json"
	"net/http"
  "log"
	"github.com/Sylvester-Kapoko/risitPap/internal/store"
)

func HandleSuggest(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if len(q) < 2 {
			_,_ = w.Write([]byte("[]"))
			return
		}
		rows, err := st.SuggestItems(q)
		if err != nil {
			http.Error(w, "internal error", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rows); err != nil {
         log.Printf("suggest encode: %v", err)
    }
  }
}
