package rules

import "testing"

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{name: "empty", config: Config{}, wantErr: true},
		{name: "sequence with prefix", config: Config{
			Prefix:   "p",
			Sequence: &SequenceConfig{Start: 1, Step: 1, Digits: 3},
		}, wantErr: true},
		{name: "empty replacement old", config: Config{
			Replacements: []Replacement{{Old: "", New: "x"}},
		}, wantErr: true},
		{name: "zero sequence step", config: Config{
			Sequence: &SequenceConfig{Start: 1, Step: 0, Digits: 3},
		}, wantErr: true},
		{name: "valid prefix", config: Config{Prefix: "p"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr && err == nil {
				t.Fatal("expected validation error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestEngine(t *testing.T) {
	engine, err := NewEngine(Config{
		Prefix: "IMG_",
		Replacements: []Replacement{
			{Old: "raw", New: "edited"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := engine.Rename("raw-photo.jpg", 0); got != "IMG_edited-photo.jpg" {
		t.Fatalf("unexpected prefix+replace result: %s", got)
	}

	seq, err := NewEngine(Config{
		Sequence: &SequenceConfig{Start: 1, Step: 2, Digits: 3},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := seq.Rename("photo.jpg", 2); got != "005.jpg" {
		t.Fatalf("unexpected sequence result: %s", got)
	}
}
