package notifications

import (
	"fmt"
	// "net/url"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// MoveCanceled has notification content for approved moves
type MoveCanceled struct {
	db      *pop.Connection
	logger  *zap.Logger
	moveID  uuid.UUID
	session *auth.Session // TODO - remove this when we move permissions up to handlers and out of models
}

// NewMoveCanceled returns a new move approval notification
func NewMoveCanceled(db *pop.Connection,
	logger *zap.Logger,
	session *auth.Session,
	moveID uuid.UUID) *MoveCanceled {

	return &MoveCanceled{
		db:      db,
		logger:  logger,
		moveID:  moveID,
		session: session,
	}
}

func (m MoveCanceled) emails() ([]emailContent, error) {
	var emails []emailContent

	move, err := models.FetchMove(m.db, m.session, m.moveID)
	if err != nil {
		return emails, err
	}

	orders, err := models.FetchOrder(m.db, m.session, move.OrdersID)
	if err != nil {
		return emails, err
	}

	serviceMember, err := models.FetchServiceMember(m.db, m.session, orders.ServiceMemberID)
	if err != nil {
		return emails, err
	}

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	if serviceMember.DutyStation.Name == "" || orders.NewDutyStation.Name == "" {
		return emails, fmt.Errorf("missing current or new duty station for service member")
	}

	if serviceMember.DutyStation.TransportationOffice.PhoneLines == nil {
		return emails, fmt.Errorf("missing contact information for origin PPPO")
	}

	// Set up various text segments. Copy comes from here:
	// https://docs.google.com/document/d/1bgE0Q_-_c93uruMP8dcNSHugXo8Pidz6YFojWBKn1Gg/edit#heading=h.h3ys1ur2qhpn
	// TODO: we will want some sort of templating system

	introText := `Your move has been canceled.`
	nextSteps := fmt.Sprintf("Your move from %s to %s with the move locator ID %s was cancelled.",
		serviceMember.DutyStation.Name, orders.NewDutyStation.Name, move.Locator)
	closingText := fmt.Sprintf("Contact your local PPPO %s at %s if you have any questions.",
		serviceMember.DutyStation.Name, serviceMember.DutyStation.TransportationOffice.PhoneLines[0].Number)

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "MOVE.MIL: Your move has been canceled.",
		htmlBody:       fmt.Sprintf("%s<br/>%s<br/>%s", introText, nextSteps, closingText),
		textBody:       fmt.Sprintf("%s\n%s\n%s", introText, nextSteps, closingText),
	}

	m.logger.Info("Sent move cancellation email to service member",
		zap.String("sevice member email address", *serviceMember.PersonalEmail))

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}