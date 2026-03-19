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
				"roster": [
					{
						"person": {
							"id": 12345,
							"fullName": "Byron Buxton"
						},
						"jerseyNumber": "25",
						"position": {
							"code": "8",
							"name": "Outfielder"
						}
					}
				]
			}`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL + "/"
	client.cache = t.TempDir()

	t.Run("Player found", func(t *testing.T) {
		player, err := client.FindPlayerByName(142, "Buxton")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if player.Person.FullName != "Byron Buxton" {
			t.Errorf("Expecting to get Buxton, got %s", player.Person.FullName)
		}
	})

	t.Run("Player not found", func(t *testing.T) {
		_, err := client.FindPlayerByName(142, "Nonexistent Player")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
