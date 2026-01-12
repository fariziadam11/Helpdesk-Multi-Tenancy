package ticket

import "fmt"

// TransformInvGateTicket transforms InvGate's incidents.details.by.view response format
// to the frontend-compatible flat structure
func TransformInvGateTicket(invGateTicket map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	// ID - direct mapping
	if id, ok := invGateTicket["id"]; ok {
		transformed["id"] = id
		transformed["inv_gate_id"] = fmt.Sprintf("%v", id)
		transformed["wrk_ticket_id"] = fmt.Sprintf("WRK-#%v", id)
		transformed["pretty_id"] = fmt.Sprintf("ARM-#%v", id)
	}

	// Extract from nested request object
	if request, ok := invGateTicket["request"].(map[string]interface{}); ok {
		// Title from request.subject
		if subject, ok := request["subject"]; ok {
			transformed["title"] = subject
		}

		// Category from request.category
		if category, ok := request["category"].(map[string]interface{}); ok {
			if categoryID, ok := category["id"]; ok {
				transformed["category_id"] = categoryID
			}
		}

		// Type from request.type
		if typeObj, ok := request["type"].(map[string]interface{}); ok {
			if typeID, ok := typeObj["id"]; ok {
				transformed["type_id"] = typeID
			}
		}
	}

	// Status
	if status, ok := invGateTicket["status"].(map[string]interface{}); ok {
		if statusLabel, ok := status["label"]; ok {
			transformed["status"] = statusLabel
		}
		if statusID, ok := status["id"]; ok {
			transformed["status_id"] = statusID
		}
	}

	// Priority
	if priority, ok := invGateTicket["priority"].(map[string]interface{}); ok {
		if priorityID, ok := priority["id"]; ok {
			transformed["priority_id"] = priorityID
		}
	}

	// Direct field mappings
	if customer, ok := invGateTicket["customer"]; ok {
		transformed["user_id"] = customer
	}

	if creator, ok := invGateTicket["creator"]; ok {
		transformed["creator_id"] = creator
	}

	if agent, ok := invGateTicket["agent"]; ok {
		transformed["assigned_id"] = agent
	}

	if helpdesk, ok := invGateTicket["helpdesk"]; ok {
		transformed["assigned_group_id"] = helpdesk
	}

	if location, ok := invGateTicket["location"]; ok {
		transformed["location_id"] = location
	}

	if category, ok := invGateTicket["category"]; ok {
		transformed["category_id"] = category
	}

	// Description fields
	if description, ok := invGateTicket["description"]; ok {
		transformed["description"] = description
	}

	if unformattedDesc, ok := invGateTicket["unformatted_description"]; ok {
		transformed["unformatted_description"] = unformattedDesc
	}

	// Timestamps - extract value from nested objects
	if creationDate, ok := invGateTicket["creation_date"].(map[string]interface{}); ok {
		if value, ok := creationDate["value"]; ok {
			transformed["created_at"] = value
			transformed["date_ocurred"] = value
		}
	}

	if lastUpdate, ok := invGateTicket["last_update"].(map[string]interface{}); ok {
		if value, ok := lastUpdate["value"]; ok {
			transformed["last_update"] = value
		}
	}

	if closingDate, ok := invGateTicket["closing_date"].(map[string]interface{}); ok {
		if value, ok := closingDate["value"]; ok {
			transformed["closed_at"] = value
			transformed["solved_at"] = value
		}
	} else {
		transformed["closed_at"] = nil
		transformed["solved_at"] = nil
	}

	if rejectionDate, ok := invGateTicket["rejection_date"]; ok {
		transformed["rejection_date"] = rejectionDate
	}

	// Source
	if source, ok := invGateTicket["source"].(map[string]interface{}); ok {
		if sourceID, ok := source["id"]; ok {
			transformed["source_id"] = sourceID
		}
	}

	// Rating
	if rating, ok := invGateTicket["rating"]; ok {
		transformed["rating"] = rating
	}

	// Collaborators (array)
	if collaborators, ok := invGateTicket["collaborators"]; ok {
		transformed["collaborators"] = collaborators
	}

	// Customer companies and groups (arrays)
	if customerCompanies, ok := invGateTicket["customer_companies"]; ok {
		transformed["customer_companies"] = customerCompanies
	}

	if customerGroups, ok := invGateTicket["customer_groups"]; ok {
		transformed["customer_groups"] = customerGroups
	}

	// SLA fields
	if firstResponseSLA, ok := invGateTicket["first_response_sla_status"]; ok {
		transformed["sla_incident_first_reply"] = firstResponseSLA
	}

	if resolutionSLA, ok := invGateTicket["resolution_sla_status"]; ok {
		transformed["sla_incident_resolution"] = resolutionSLA
	}

	// Subject (duplicate of title for compatibility)
	if subject, ok := invGateTicket["subject"]; ok {
		transformed["subject"] = subject
	}

	// Type (duplicate for compatibility)
	if typeObj, ok := invGateTicket["type"].(map[string]interface{}); ok {
		if typeID, ok := typeObj["id"]; ok {
			transformed["type_id"] = typeID
		}
	}

	// Waiting for
	if waitingFor, ok := invGateTicket["waiting_for"].(map[string]interface{}); ok {
		if waitingType, ok := waitingFor["type"].(map[string]interface{}); ok {
			if typeID, ok := waitingType["id"]; ok {
				transformed["waiting_for_type_id"] = typeID
			}
		}
		if value, ok := waitingFor["value"]; ok {
			transformed["waiting_for_value"] = value
		}
	}

	return transformed
}

// TransformInvGateTicketList transforms a list of InvGate tickets
func TransformInvGateTicketList(invGateTickets []interface{}) []map[string]interface{} {
	transformed := make([]map[string]interface{}, 0, len(invGateTickets))

	for _, ticket := range invGateTickets {
		if ticketMap, ok := ticket.(map[string]interface{}); ok {
			transformed = append(transformed, TransformInvGateTicket(ticketMap))
		}
	}

	return transformed
}
