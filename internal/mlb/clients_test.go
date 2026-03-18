package mlb

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindPlayerByName(t *testing.T) {
	//set up fake server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/teams/142/roster" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{
			 roster: [
				{
					person: {
						id: 12345,
						fullName: "Buxton"
					}
				}
			 ]
			}`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	}))
	defer server.Close()
}
