package services

import (
	"fmt"

	"github.com/ducminhgd/gossip-bot/internal/models"
)

// SourceService handles operations related to news sources
type SourceService struct {
	sources []models.Source
}

// NewSourceService creates a new SourceService
func NewSourceService(sources []models.Source) *SourceService {
	return &SourceService{
		sources: sources,
	}
}

// GetSources returns all sources
func (s *SourceService) GetSources() []models.Source {
	return s.sources
}

// GetSourceByName returns a source by name
func (s *SourceService) GetSourceByName(name string) (models.Source, error) {
	for _, source := range s.sources {
		if source.Name == name {
			return source, nil
		}
	}
	return models.Source{}, fmt.Errorf("source not found: %s", name)
}
