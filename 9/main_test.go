package main

import "testing"

type test struct {
	name     string
	input    string
	expected string
	wantErr  bool
}

func TestUnzip(t *testing.T) {
	data := []test{
		{"1", "a4bc2d5e", "aaaabccddddde", false},
		{"2", "abcd", "abcd", false},
		{"3", "a1b2c3", "abbccc", false},
		{"4", "z9", "zzzzzzzzz", false},
		{"5", "\\3", "3", false},
		{"6", "a\\4b", "a4b", false},
		{"7", "a0b1", "b", false},
		{"8", "a2\\3", "aa3", false},
		{"9", "", "", true},
		{"10", "3abc", "", true},
		{"11", "a\\", "", true},
		{"12", "a12", "", true},
		{"13", "б3в2", "бббвв", false},
		{"14", "界2世3", "界界世世世", false},
		{"15", "🙂2🙃", "🙂🙂🙃", false},
		{"16", "\\32", "33", false},
	}
	for _, ts := range data {
		t.Run(ts.name, func(t *testing.T) {
			got, err := unzipString(ts.input)

			if ts.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil, result = %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != ts.expected {
				t.Fatalf("got %v, want %v,", got, ts.expected)
			}
		})
	}
}
