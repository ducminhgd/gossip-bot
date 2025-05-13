package services

import (
	"testing"

	"github.com/ducminhgd/gossip-bot/internal/models"
)

func TestNewSourceService(t *testing.T) {
	// Create test sources
	sources := []models.Source{
		{
			Name:  "TestSource1",
			Type:  "test",
			URL:   "https://example.com",
			Limit: 10,
		},
		{
			Name:      "TestSource2",
			Type:      "test",
			URL:       "https://example.org",
			Limit:     5,
			SubSource: "sub",
		},
	}

	// Create service
	service := NewSourceService(sources)

	// Check that the service was created correctly
	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	// Check that the sources were stored correctly
	if len(service.sources) != len(sources) {
		t.Fatalf("Expected %d sources, got %d", len(sources), len(service.sources))
	}

	for i, source := range service.sources {
		if source.Name != sources[i].Name {
			t.Errorf("Expected source name %s, got %s", sources[i].Name, source.Name)
		}
		if source.Type != sources[i].Type {
			t.Errorf("Expected source type %s, got %s", sources[i].Type, source.Type)
		}
		if source.URL != sources[i].URL {
			t.Errorf("Expected source URL %s, got %s", sources[i].URL, source.URL)
		}
		if source.Limit != sources[i].Limit {
			t.Errorf("Expected source limit %d, got %d", sources[i].Limit, source.Limit)
		}
		if source.SubSource != sources[i].SubSource {
			t.Errorf("Expected source subsource %s, got %s", sources[i].SubSource, source.SubSource)
		}
	}
}

func TestGetSources(t *testing.T) {
	// Create test sources
	sources := []models.Source{
		{
			Name:  "TestSource1",
			Type:  "test",
			URL:   "https://example.com",
			Limit: 10,
		},
		{
			Name:      "TestSource2",
			Type:      "test",
			URL:       "https://example.org",
			Limit:     5,
			SubSource: "sub",
		},
	}

	// Create service
	service := NewSourceService(sources)

	// Get sources
	result := service.GetSources()

	// Check that the correct sources were returned
	if len(result) != len(sources) {
		t.Fatalf("Expected %d sources, got %d", len(sources), len(result))
	}

	for i, source := range result {
		if source.Name != sources[i].Name {
			t.Errorf("Expected source name %s, got %s", sources[i].Name, source.Name)
		}
	}
}

func TestGetSourceByName_Found(t *testing.T) {
	// Create test sources
	sources := []models.Source{
		{
			Name:  "TestSource1",
			Type:  "test",
			URL:   "https://example.com",
			Limit: 10,
		},
		{
			Name:      "TestSource2",
			Type:      "test",
			URL:       "https://example.org",
			Limit:     5,
			SubSource: "sub",
		},
	}

	// Create service
	service := NewSourceService(sources)

	// Get source by name
	result, err := service.GetSourceByName("TestSource2")

	// Check that no error was returned
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that the correct source was returned
	if result.Name != "TestSource2" {
		t.Errorf("Expected source name TestSource2, got %s", result.Name)
	}
	if result.Type != "test" {
		t.Errorf("Expected source type test, got %s", result.Type)
	}
	if result.URL != "https://example.org" {
		t.Errorf("Expected source URL https://example.org, got %s", result.URL)
	}
	if result.Limit != 5 {
		t.Errorf("Expected source limit 5, got %d", result.Limit)
	}
	if result.SubSource != "sub" {
		t.Errorf("Expected source subsource sub, got %s", result.SubSource)
	}
}

func TestGetSourceByName_NotFound(t *testing.T) {
	// Create test sources
	sources := []models.Source{
		{
			Name:  "TestSource1",
			Type:  "test",
			URL:   "https://example.com",
			Limit: 10,
		},
	}

	// Create service
	service := NewSourceService(sources)

	// Get source by name
	_, err := service.GetSourceByName("NonExistentSource")

	// Check that an error was returned
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that the error message is correct
	expectedError := "source not found: NonExistentSource"
	if err.Error() != expectedError {
		t.Errorf("Expected error message %q, got %q", expectedError, err.Error())
	}
}
