package models_test

import (
	"testing"

	"github.com/Uvaisbugh/gophervault/internal/models"
)

// 📚 Go concept: In Go, unit tests live in files ending in _test.go.
// A test function starts with "Test" and takes a single parameter: (t *testing.T).
// We use the testing standard library — no external test runners required!

func TestNewEntry(t *testing.T) {
	title := "Git Commit Amend"
	body := "git commit --amend --no-edit"
	category := "git"
	tags := []string{"git", "vcs"}

	entry := models.NewEntry(title, body, category, tags)

	if entry.ID == "" {
		t.Fatal("expected generated ID to not be empty")
	}
	if entry.Title != title {
		t.Errorf("expected title %q, got %q", title, entry.Title)
	}
	if entry.Body != body {
		t.Errorf("expected body %q, got %q", body, entry.Body)
	}
	if entry.Category != category {
		t.Errorf("expected category %q, got %q", category, entry.Category)
	}
	if len(entry.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(entry.Tags))
	}
	if entry.CreatedAt.IsZero() || entry.UpdatedAt.IsZero() {
		t.Error("expected CreatedAt and UpdatedAt to be set")
	}
}

func TestEntry_HasTag(t *testing.T) {
	entry := models.NewEntry("Docker", "docker run", "devops", []string{"Docker", "CONTAINER"})

	// 📚 Go concept: Table-driven tests!
	// This is the idiomatic, standard way to write Go tests.
	tests := []struct {
		name     string
		tag      string
		expected bool
	}{
		{"exact match", "Docker", true},
		{"lowercase match", "docker", true},
		{"uppercase match", "DOCKER", true},
		{"container tag", "container", true},
		{"missing tag", "kubernetes", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entry.HasTag(tt.tag)
			if got != tt.expected {
				t.Errorf("HasTag(%q) = %v; want %v", tt.tag, got, tt.expected)
			}
		})
	}
}

func TestEntry_MatchesSearch(t *testing.T) {
	entry := models.NewEntry("Kubectl get pods", "kubectl get pods -A -o wide", "k8s", []string{"kubernetes", "debug"})

	tests := []struct {
		name     string
		query    string
		expected bool
	}{
		{"match title", "pods", true},
		{"case insensitive title", "PODS", true},
		{"match body", "-o wide", true},
		{"match tag", "debug", true},
		{"no match", "terraform", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entry.MatchesSearch(tt.query)
			if got != tt.expected {
				t.Errorf("MatchesSearch(%q) = %v; want %v", tt.query, got, tt.expected)
			}
		})
	}
}

func TestEntry_String(t *testing.T) {
	entry := models.NewEntry("My Title", "echo hello", "shell", []string{"bash"})
	str := entry.String()

	if str == "" {
		t.Error("expected non-empty string representation")
	}
}
