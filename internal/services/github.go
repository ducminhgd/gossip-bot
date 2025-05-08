package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ducminhgd/gossip-bot/internal/models"
	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

// GithubService handles operations related to GitHub
type GithubService struct {
	client *github.Client
	owner  string
	repo   string
}

// NewGithubService creates a new GithubService
func NewGithubService(token, owner, repo string) *GithubService {
	// Create OAuth2 token source
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	// Create GitHub client
	client := github.NewClient(tc)

	return &GithubService{
		client: client,
		owner:  owner,
		repo:   repo,
	}
}

// CreateIssue creates a new GitHub issue with the given title and body
func (s *GithubService) CreateIssue(title, body string) (*github.Issue, error) {
	// Create issue request
	issue := &github.IssueRequest{
		Title: github.String(title),
		Body:  github.String(body),
	}

	// Create issue
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createdIssue, _, err := s.client.Issues.Create(ctx, s.owner, s.repo, issue)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return createdIssue, nil
}

// GenerateIssueContent generates the content for a GitHub issue
func (s *GithubService) GenerateIssueContent(newsMap map[string][]models.News) (string, error) {
	// Get current date
	now := time.Now().UTC()
	date := now.Format("2006-01-02")

	// Build issue content
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", date))

	// Add news for each source
	for source, newsList := range newsMap {
		if len(newsList) == 0 {
			continue
		}

		// Add source header
		sb.WriteString(fmt.Sprintf("## %s\n\n", source))

		// Add news items - only titles, no descriptions
		for i, news := range newsList {
			sb.WriteString(fmt.Sprintf("%d. [%s](%s)\n", i+1, news.Title, news.URL))
		}

		sb.WriteString("\n")
	}

	return sb.String(), nil
}
