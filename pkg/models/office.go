package models

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// AccountingInfo represents a single move queue item within a queue.
type AccountingInfo struct {
	TAC           string                         `json:"tac" db:"tac"`
	DeptIndicator internalmessages.DeptIndicator `json:"dept_indicator" db:"dept_indicator"`
}

// FetchAccountingInfo gets accounting information for a specific move
func FetchAccountingInfo(db *pop.Connection, moveID uuid.UUID) (AccountingInfo, error) {
	accountingInfo := AccountingInfo{}
	// TODO: replace hardcoded values with actual query values once data is available
	query := `
		SELECT 'F8J1' as tac,'AIRFORCE' AS dept_indicator
		FROM moves
		WHERE moves.id = $1
	`
	err := db.RawQuery(query, moveID).First(&accountingInfo)
	return accountingInfo, err
}