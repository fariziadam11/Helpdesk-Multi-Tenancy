package ticket

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"werk-ticketing/internal/errors"
	"werk-ticketing/internal/invgate"
)

func (s *service) CreateTicket(ctx context.Context, tenantID string, req TicketRequest, creatorEmail string) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByEmail(ctx, tenantID, creatorEmail)
	if err != nil {
		s.logger.WithError(err).WithField("creatorEmail", creatorEmail).Error("failed to get user by email")
		return nil, errors.NewAppError(
			errors.ErrCodeInternal,
			"failed to retrieve user information",
			err,
		)
	}
	if user == nil {
		return nil, errors.NewAppError(
			errors.ErrCodeNotFound,
			"user not found",
			nil,
		)
	}

	invgateUserID := user.InvGateUserID
	if invgateUserID == 0 {
		s.logger.WithField("creatorEmail", creatorEmail).Error("user has no invgate_user_id")
		return nil, errors.NewAppError(
			errors.ErrCodeInternal,
			"user is not synchronized with InvGate. Please contact administrator.",
			nil,
		)
	}

	payload := invgate.CreateTicketPayload{
		SourceID:    req.SourceID,
		CreatorID:   invgateUserID,
		CustomerID:  invgateUserID,
		CategoryID:  req.CategoryID,
		TypeID:      req.TypeID,
		PriorityID:  req.PriorityID,
		Title:       req.Title,
		Description: req.Description,
		DateOcurred: req.DateOcurred,
	}

	var invgateResp map[string]interface{}
	if len(req.AttachmentFiles) > 0 {
		invgateResp, err = s.client.CreateTicketWithAttachments(ctx, payload, req.AttachmentFiles)
	} else {
		invgateResp, err = s.client.CreateTicket(ctx, payload)
	}

	if err != nil {
		s.logger.WithError(err).
			WithFields(logrus.Fields{
				"creator_id":      payload.CreatorID,
				"customer_id":     payload.CustomerID,
				"category_id":     payload.CategoryID,
				"type_id":         payload.TypeID,
				"priority_id":     payload.PriorityID,
				"source_id":       payload.SourceID,
				"has_attachments": len(req.AttachmentFiles) > 0,
			}).
			Error("failed to create ticket in InvGate")

		errorMsg := "failed to create ticket in external service"
		if err.Error() != "" {
			errorMsg = fmt.Sprintf("InvGate API error: %s", err.Error())
		}

		return nil, errors.NewAppError(
			errors.ErrCodeExternalService,
			errorMsg,
			err,
		)
	}

	// Extract InvGate ID for logging
	invGateID, err := extractInvGateID(invgateResp)
	if err != nil {
		s.logger.WithError(err).WithField("response", invgateResp).Warn("failed to extract InvGate ID from response")
		invGateID = ""
	}

	s.logger.WithFields(logrus.Fields{
		"invGateID":    invGateID,
		"title":        req.Title,
		"creatorEmail": creatorEmail,
	}).Info("ticket created successfully in InvGate")

	return invgateResp, nil
}
