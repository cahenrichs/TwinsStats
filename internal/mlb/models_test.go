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
				Stat: json.RawMessage([]byte(tt.input)),
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

func TestPitchingStats(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *PitchingStats
		wantErr bool
	}{
		{
			name:    "valid stats",
			input:   `{"wins": 5, "losses": 3, "era": 3.50, "strikeouts": 80}`,
			want:    &PitchingStats{Wins: 5, Losses: 3, ERA: 3.50, Strikeouts: 80},
			wantErr: false,
		},
		{
			name:    "malformed json",
			input:   `{"wins": "invalid"}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			split := &StatSplit{
				Stat: json.RawMessage([]byte(tt.input)),
			}
			got, err := split.GetPitchingStats()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPitchingStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPitchingStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
