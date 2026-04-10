package sentinel_transport_http

import "net/http"

func GetImagery(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(city))
}
