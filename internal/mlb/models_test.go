package mlb

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestHittingStats(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *HittingStats
		wantErr bool
	}{
		{
			name:    "valid stats",
			input:   `{"hits": 10, "atBats": 30, "gamesPlayed": 10, "homeRuns": 2}`,
			want:    &HittingStats{Hits: 10, AtBats: 30, GamesPlayed: 10, HomeRuns: 2},
			wantErr: false,
		},
		{
			name:    "malformed json",
			input:   `{"hits": "invalid"}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			split := &StatSplit{
				Stat: json.RawMessage(tt.input),
			}
			got, err := split.GetHittingStats()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHittingStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHittingStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
