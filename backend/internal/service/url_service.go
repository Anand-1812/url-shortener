// Package service for the shortening the url
package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"

	"backend/internal/repository"
)

const charset = "abdcdefghijklmonpqrstuvwxyzAPCDEFGHIJLKMONPQRSTUVWXYZ012346789"

type URLService struct {
	repo *repository.URLRepository
}

func NewURLService(repo *repository.URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *URLService) generateCode(length int) (string, error) {
	code := make([]byte, length)
	for i := range code {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		code[i] = charset[idx.Int64()]
	}

	return string(code), nil
}

func (s *URLService) Shorten(ctx context.Context, rawURL string) (repository.URLRecords, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return repository.URLRecords{}, fmt.Errorf("invalid url: must include http:// or https://")
	}

	shortCode, err := s.generateCode(6)
	if err != nil {
		return repository.URLRecords{}, err
	}

	return s.repo.Create(ctx, rawURL, shortCode)
}

func (s *URLService) Resolve(ctx context.Context, shortCode string) (repository.URLRecords, error) {
	return s.repo.FindAndIncrementsClicks(ctx, shortCode)
}
