package invgate

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// GetTicketsByView retrieves tickets using InvGate's incidents.details.by.view endpoint
// This endpoint returns detailed ticket information with metadata for a specific view
// creatorID filters tickets to only show those created by the specified InvGate user ID
func (s *service) GetTicketsByView(ctx context.Context, viewID int, pageKey string, creatorID int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("view_id", strconv.Itoa(viewID))

	if pageKey != "" {
		params.Set("page_key", pageKey)
	}

	// Filter by creator ID if provided
	if creatorID > 0 {
		params.Set("creator_id", strconv.Itoa(creatorID))
	}

	return s.doRequest(ctx, http.MethodGet, "incidents.details.by.view", nil, params)
}
