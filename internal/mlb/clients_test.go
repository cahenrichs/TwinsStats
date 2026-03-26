package mlb

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchAndCache(t *testing.T) {
	tests := []struct {
		name              string
		teamId            int
		setupCache        func(cacheDir string, filename string)
		mockAPIResponse   string
		mockAPIStatusCode int
		expectAPICall     bool
		expectError       bool
		expectedName      string
	}{
		{
			name:              "Cache miss",
			teamId:            142,
			setupCache:        func(cacheDir string, filename string) {},
			mockAPIResponse:   `{"roster": [{"person": {"fullName": "Byron Buxton"}}]}`,
			mockAPIStatusCode: http.StatusOK,
			expectAPICall:     true,
			expectError:       false,
			expectedName:      "Byron Buxton",
		},
		{
			name: "Cache Hit - No API Call",
			setupCache: func(cacheDir string, fileName string) {
				// Pre-populate the cache with Royce Lewis
				data := []byte(`{"roster": [{"person": {"fullName": "Royce Lewis"}}]}`)
				_ = os.MkdirAll(cacheDir, 0755)
				_ = os.WriteFile(filepath.Join(cacheDir, fileName), data, 0644)
			},
			mockAPIResponse:   `{}`,
			mockAPIStatusCode: http.StatusInternalServerError, // API error shouldn't matter if cache is hit
			expectAPICall:     false,
			expectedName:      "Royce Lewis",
			expectError:       false,
		},
		{
			name: "Corrupted Cache - Falls back to API",
			setupCache: func(dir, file string) {
				_ = os.MkdirAll(dir, 0755)
				_ = os.WriteFile(filepath.Join(dir, file), []byte(`{invalid-json}`), 0644)
			},
			mockAPIResponse:   `{"roster": [{"person": {"fullName": "Carlos Correa"}}]}`,
			mockAPIStatusCode: http.StatusOK,
			expectAPICall:     true,
			expectedName:      "Carlos Correa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheDir := t.TempDir()
			fileName := "test_data.json"
			tt.setupCache(cacheDir, fileName)
			requestCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCount++
				w.WriteHeader(tt.mockAPIStatusCode)
				fmt.Fprintln(w, tt.mockAPIResponse)
			}))
			defer server.Close()
			client := &Client{
				baseURL: server.URL + "/",
				http:    server.Client(),
				cache:   cacheDir,
			}
			result, err := fetchAndCache[RosterResponse](client, server.URL, fileName)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedName, result.Roster[0].Person.FullName)

				if tt.expectAPICall {
					assert.Equal(t, 1, requestCount, "Expected an API call")
				} else {
					assert.Equal(t, 0, requestCount, "Expected NO API call (Cache should have been hit)")
				}
			}
		})
	}
}
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
