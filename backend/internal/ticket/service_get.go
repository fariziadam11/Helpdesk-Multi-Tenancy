package ticket

import (
	"context"
	"sort"

	"werk-ticketing/internal/errors"
)

func (s *service) GetTickets(ctx context.Context, tenantID, creatorID string, page, limit int) (map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Get user from database to retrieve InvGateUserID
	user, err := s.userRepo.GetByEmail(ctx, tenantID, creatorID)
	if err != nil {
		s.logger.WithError(err).
			WithField("email", creatorID).
			Error("failed to get user from database")
		return nil, errors.NewAppError(
			errors.ErrCodeInternal,
			"failed to fetch user information",
			err,
		)
	}

	if user == nil {
		s.logger.WithField("email", creatorID).
			Warn("user not found")
		return nil, errors.NewAppError(
			errors.ErrCodeNotFound,
			"user not found",
			nil,
		)
	}

	// Calculate page_key for InvGate pagination
	// For now, we'll use empty string for first page
	// InvGate will return next_page_key for subsequent pages
	pageKey := ""

	// Call InvGate API with view_id=7 (hardcoded as per requirement)
	// Note: We don't pass creatorID to InvGate because the API doesn't support filtering by creator
	// We'll filter the results ourselves after getting the response
	resp, err := s.client.GetTicketsByView(ctx, 7, pageKey, 0)
	if err != nil {
		s.logger.WithError(err).
			WithField("invGateUserID", user.InvGateUserID).
			Error("failed to get tickets from InvGate")
		return nil, errors.NewAppError(
			errors.ErrCodeExternalService,
			"failed to fetch tickets from external service",
			err,
		)
	}

	// Extract data array from response
	var allTickets []map[string]interface{}
	if data, ok := resp["data"].([]interface{}); ok {
		// Transform each ticket to frontend-compatible format
		allTickets = TransformInvGateTicketList(data)
	}

	// Sort tickets by ID descending (newest first)
	// This ensures the most recent tickets appear at the top
	sort.Slice(allTickets, func(i, j int) bool {
		idI, okI := allTickets[i]["id"]
		idJ, okJ := allTickets[j]["id"]

		if !okI || !okJ {
			return false
		}

		// Convert to float64 for comparison (JSON unmarshaling uses float64 for numbers)
		var numI, numJ float64
		switch v := idI.(type) {
		case float64:
			numI = v
		case int:
			numI = float64(v)
		case int64:
			numI = float64(v)
		}

		switch v := idJ.(type) {
		case float64:
			numJ = v
		case int:
			numJ = float64(v)
		case int64:
			numJ = float64(v)
		}

		return numI > numJ // Descending order (newest first)
	})

	// Filter tickets by creator (InvGate user ID)
	// Only include tickets where creator field matches the logged-in user's InvGateUserID
	var filteredTickets []map[string]interface{}
	for _, ticket := range allTickets {
		// Check if creator field matches user's InvGateUserID
		if creatorID, ok := ticket["creator_id"]; ok {
			// Handle both int and float64 from JSON unmarshaling
			var creatorIDInt int
			switch v := creatorID.(type) {
			case int:
				creatorIDInt = v
			case float64:
				creatorIDInt = int(v)
			case int64:
				creatorIDInt = int(v)
			default:
				continue
			}

			if creatorIDInt == user.InvGateUserID {
				filteredTickets = append(filteredTickets, ticket)
			}
		}
	}

	// Build pagination metadata based on filtered results
	totalCount := len(filteredTickets)
	totalPages := 1
	if totalCount > 0 {
		totalPages = (totalCount + limit - 1) / limit
	}

	// Apply client-side pagination to match requested page/limit
	startIdx := (page - 1) * limit
	endIdx := startIdx + limit

	var paginatedTickets []map[string]interface{}
	if startIdx >= len(filteredTickets) {
		paginatedTickets = []map[string]interface{}{}
	} else {
		if endIdx > len(filteredTickets) {
			endIdx = len(filteredTickets)
		}
		paginatedTickets = filteredTickets[startIdx:endIdx]
	}

	return map[string]interface{}{
		"data": paginatedTickets,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       totalCount,
			"total_pages": totalPages,
			"has_next":    page < totalPages,
			"has_prev":    page > 1,
		},
	}, nil
}

func (s *service) GetTicketDetail(ctx context.Context, ticketID string) (map[string]interface{}, error) {
	resp, err := s.client.GetTicketDetail(ctx, ticketID)
	if err != nil {
		s.logger.WithError(err).
			WithField("ticketID", ticketID).
			Error("failed to get ticket detail from InvGate")
		return nil, errors.NewAppError(
			errors.ErrCodeExternalService,
			"failed to fetch ticket detail from external service",
			err,
		)
	}

	if statusID, ok := resp["status_id"]; ok {
		resp["status"] = getStatusName(statusID)
	}

	return resp, nil
}
