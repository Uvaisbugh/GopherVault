package service_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Uvaisbugh/gophervault/internal/service"
	"github.com/Uvaisbugh/gophervault/internal/storage"
)

// helper to instantiate a test service with an isolated temp file.
// 📚 Go concept: t.TempDir() returns a temporary directory that is
// automatically removed when the test and all its subtests complete!
func setupTestService(t *testing.T) (*service.VaultService, string) {
	t.Helper()
	tempFile := filepath.Join(t.TempDir(), "test_vault.json")
	store, err := storage.NewJSONStore(tempFile)
	if err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	return service.NewVaultService(store), tempFile
}

func TestVaultService_AddEntry_Validation(t *testing.T) {
	svc, _ := setupTestService(t)

	// Test empty title
	_, err := svc.AddEntry("", "some body", "devops", nil)
	if !errors.Is(err, service.ErrEmptyTitle) {
		t.Errorf("expected ErrEmptyTitle, got %v", err)
	}

	// Test empty body
	_, err = svc.AddEntry("Valid Title", "", "devops", nil)
	if !errors.Is(err, service.ErrEmptyBody) {
		t.Errorf("expected ErrEmptyBody, got %v", err)
	}

	// Test valid entry with default category
	entry, err := svc.AddEntry("Valid Title", "some command", "", []string{"tag1", " ", "tag2"})
	if err != nil {
		t.Fatalf("unexpected error adding valid entry: %v", err)
	}

	if entry.Category != "general" {
		t.Errorf("expected default category 'general', got %q", entry.Category)
	}
	if len(entry.Tags) != 2 {
		t.Errorf("expected 2 clean tags (whitespace tag stripped), got %d", len(entry.Tags))
	}
}

func TestVaultService_GetAndSearch(t *testing.T) {
	svc, _ := setupTestService(t)

	e1, err := svc.AddEntry("Docker prune", "docker system prune", "docker", []string{"devops", "cleanup"})
	if err != nil {
		t.Fatalf("failed to add entry: %v", err)
	}

	e2, err := svc.AddEntry("Git status", "git status -s", "git", []string{"vcs"})
	if err != nil {
		t.Fatalf("failed to add entry: %v", err)
	}

	// Test exact ID lookup
	found, err := svc.GetEntry(e1.ID)
	if err != nil || found.ID != e1.ID {
		t.Fatalf("failed to get entry by exact ID: %v", err)
	}

	// Test prefix ID lookup
	foundPrefix, err := svc.GetEntry(e2.ID[:4])
	if err != nil || foundPrefix.ID != e2.ID {
		t.Fatalf("failed to get entry by prefix ID: %v", err)
	}

	// Test search
	results, err := svc.SearchEntries("prune")
	if err != nil {
		t.Fatalf("search error: %v", err)
	}
	if len(results) != 1 || results[0].ID != e1.ID {
		t.Errorf("expected 1 result matching 'prune', got %d", len(results))
	}
}

func TestVaultService_ListAndFilter(t *testing.T) {
	svc, _ := setupTestService(t)

	_, _ = svc.AddEntry("Docker 1", "run", "docker", []string{"devops"})
	_, _ = svc.AddEntry("Docker 2", "ps", "docker", []string{"monitor"})
	_, _ = svc.AddEntry("Git 1", "log", "git", []string{"devops"})

	// Filter by tag
	devopsList, err := svc.ListEntries("", "devops")
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(devopsList) != 2 {
		t.Errorf("expected 2 devops entries, got %d", len(devopsList))
	}

	// Filter by category
	dockerList, err := svc.ListEntries("docker", "")
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(dockerList) != 2 {
		t.Errorf("expected 2 docker entries, got %d", len(dockerList))
	}

	// Combined filter
	combined, err := svc.ListEntries("docker", "devops")
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(combined) != 1 {
		t.Errorf("expected 1 combined entry, got %d", len(combined))
	}
}

func TestVaultService_Delete(t *testing.T) {
	svc, _ := setupTestService(t)

	entry, err := svc.AddEntry("Temporary note", "to be deleted", "temp", nil)
	if err != nil {
		t.Fatalf("failed to add entry: %v", err)
	}

	if err := svc.DeleteEntry(entry.ID); err != nil {
		t.Fatalf("failed to delete entry: %v", err)
	}

	_, err = svc.GetEntry(entry.ID)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected ErrNotFound after deletion, got %v", err)
	}
}

func TestVaultService_GetStats(t *testing.T) {
	svc, _ := setupTestService(t)

	_, _ = svc.AddEntry("Docker prune", "prune", "docker", []string{"devops", "cleanup"})
	_, _ = svc.AddEntry("Docker stop", "stop", "docker", []string{"devops"})
	_, _ = svc.AddEntry("Go test", "go test ./...", "golang", []string{"devops", "test"})

	stats, err := svc.GetStats()
	if err != nil {
		t.Fatalf("stats error: %v", err)
	}

	if stats.TotalEntries != 3 {
		t.Errorf("expected 3 entries, got %d", stats.TotalEntries)
	}
	if stats.Categories["docker"] != 2 {
		t.Errorf("expected 2 docker category entries, got %d", stats.Categories["docker"])
	}
	if len(stats.PopularTags) == 0 || stats.PopularTags[0].Tag != "devops" || stats.PopularTags[0].Count != 3 {
		t.Errorf("expected devops to be top tag with count 3, got %+v", stats.PopularTags)
	}
}

func TestVaultService_ExportAndImport(t *testing.T) {
	svc1, _ := setupTestService(t)
	_, _ = svc1.AddEntry("Exported Snippet", "echo export", "test", []string{"backup"})

	backupFile := filepath.Join(t.TempDir(), "backup.json")
	if err := svc1.Export(backupFile); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	svc2, _ := setupTestService(t)
	importedCount, err := svc2.Import(backupFile)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}

	if importedCount != 1 {
		t.Errorf("expected 1 imported entry, got %d", importedCount)
	}

	results, err := svc2.SearchEntries("Exported")
	if err != nil || len(results) != 1 {
		t.Errorf("expected to find imported entry in new vault, got %d results", len(results))
	}
}
