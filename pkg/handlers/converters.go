package handlers

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/markbates/pop/nulls"
	"github.com/satori/go.uuid"
)

// These functions clean up converting value types into pointers which
// otherwise takes multiple lines.
func stringPointer(s string) *string { return &s }

// These functions facilitate converting from the go types the db uses
// into the strfmt types that go-swagger uses for payloads.
func fmtUUID(u uuid.UUID) *strfmt.UUID {
	fmtUUID := strfmt.UUID(u.String())
	return &fmtUUID
}

func fmtNullUUID(u nulls.UUID) *strfmt.UUID {
	if u.Valid {
		return fmtUUID(u.UUID)
	}
	return nil
}

func fmtNullBool(b nulls.Bool) *bool {
	if b.Valid {
		return &b.Bool
	}
	return nil
}

func fmtDateTime(dateTime time.Time) *strfmt.DateTime {
	fmtDateTime := strfmt.DateTime(dateTime)
	return &fmtDateTime
}

func fmtDate(date time.Time) *strfmt.Date {
	fmtDate := strfmt.Date(date)
	return &fmtDate
}

func fmtInt64(i int) *int64 {
	fmtInt := int64(i)
	return &fmtInt
}
