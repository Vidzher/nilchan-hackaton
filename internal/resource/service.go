package resource

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nilchan-hackaton/internal/parser"
)

const firecrawlAttempts = 2

type pageParser interface {
	ParsePage(ctx context.Context, pageURL string) (parser.Page, error)
}

type resourceRepository interface {
	FindByURL(ctx context.Context, userID int64, resourceURL string) (*Resource, error)
	FindByID(ctx context.Context, userID, resourceID int64) (*Summary, error)
	List(ctx context.Context, userID int64, status, tag string) ([]Summary, error)
	CheckCapacity(ctx context.Context, userID int64, purchaseOverflow bool) error
	CreateProcessing(ctx context.Context, draft Resource, purchaseOverflow bool) (*Resource, error)
	RetryFailed(ctx context.Context, resourceID, userID int64, allowOverflowPurchase bool) (overflowSlotPurchased bool, err error)
}

type submitter interface {
	Submit(resource *Resource)
}

type Service struct {
	repo             resourceRepository
	parser           pageParser
	submitter        submitter
	firecrawlTimeout time.Duration
}

func NewService(repo resourceRepository, pageParser pageParser, submitter submitter, firecrawlTimeout time.Duration) *Service {
	return &Service{
		repo:             repo,
		parser:           pageParser,
		submitter:        submitter,
		firecrawlTimeout: firecrawlTimeout,
	}
}

func (s *Service) Get(ctx context.Context, userID, resourceID int64) (*Summary, error) {
	found, err := s.repo.FindByID(ctx, userID, resourceID)
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrNotFound
	}
	return found, nil
}

func (s *Service) List(ctx context.Context, userID int64, status, tag string) ([]Summary, error) {
	if status != "" && !Status(status).Valid() {
		return nil, ErrInvalidStatus
	}
	return s.repo.List(ctx, userID, status, tag)
}

func (s *Service) Create(ctx context.Context, userID int64, request CreateResourceRequest) (*Resource, error) {

	normalizedURL, err := NormalizeURL(request.URL)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByURL(ctx, userID, normalizedURL)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		switch existing.Status {
		case StatusProcessing:
			return existing, nil
		case StatusFailed:
			overflowSlotPurchased, err := s.repo.RetryFailed(ctx, existing.ID, userID, request.PurchaseOverflowSlot)
			if err != nil {
				return nil, err
			}
			existing.Status = StatusProcessing
			existing.PurchasedOverflowSlot = overflowSlotPurchased
			s.submitter.Submit(existing)
			return existing, nil
		default:
			return nil, ErrDuplicate
		}
	}

	if err := s.repo.CheckCapacity(ctx, userID, request.PurchaseOverflowSlot); err != nil {
		return nil, err
	}

	page, err := s.parsePage(ctx, normalizedURL)
	if err != nil {
		return nil, err
	}
	page, err = validatePage(page)
	if err != nil {
		return nil, err
	}

	resource := Resource{
		UserID:  userID,
		URL:     normalizedURL,
		Title:   page.Title,
		Tags:    page.Tags,
		Content: page.Content,
	}

	created, err := s.repo.CreateProcessing(ctx, resource, request.PurchaseOverflowSlot)
	if errors.Is(err, ErrDuplicate) {
		existing, findErr := s.repo.FindByURL(ctx, userID, normalizedURL)
		if findErr != nil {
			return nil, findErr
		}
		if existing != nil && existing.Status == StatusProcessing {
			return existing, nil
		}
	}
	if err != nil {
		return nil, err
	}

	s.submitter.Submit(created)
	return created, nil
}

func (s *Service) parsePage(ctx context.Context, resourceURL string) (parser.Page, error) {
	parseCtx, cancel := context.WithTimeout(ctx, s.firecrawlTimeout)
	defer cancel()

	var lastErr error
	for range firecrawlAttempts {
		page, err := s.parser.ParsePage(parseCtx, resourceURL)
		if err == nil {
			return page, nil
		}
		lastErr = err
		if parseCtx.Err() != nil {
			return parser.Page{}, fmt.Errorf("%w: %v", ErrFirecrawlTimeout, parseCtx.Err())
		}
	}
	return parser.Page{}, fmt.Errorf("%w: %v", ErrFirecrawlFailed, lastErr)
}
