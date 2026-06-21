package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func pgtypeUUIDToDomain(id pgtype.UUID) uuid.UUID {
	return id.Bytes
}

func timestamptzToDomain(t pgtype.Timestamptz) time.Time {
	return t.Time
}

func timeToTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func stringToPgtypeText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func pgtypeTextToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func isNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
