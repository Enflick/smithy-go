package rand

import (
	"github.com/aws/smithy-go/internal/uuid"
	"io"
)

const dash byte = '-'

// UUIDIdempotencyToken provides a utility to get idempotency tokens in the
// UUID format.
type UUIDIdempotencyToken struct {
	uuid *UUID
}

// NewUUIDIdempotencyToken returns a idempotency token provider returning
// tokens in the UUID random format using the reader provided.
func NewUUIDIdempotencyToken(r io.Reader) *UUIDIdempotencyToken {
	return &UUIDIdempotencyToken{uuid: NewUUID(r)}
}

// GetIdempotencyToken returns a random UUID value for Idempotency token.
func (u UUIDIdempotencyToken) GetIdempotencyToken() (string, error) {
	return u.uuid.GetUUID()
}

// UUID provides computing random UUID version 4 values from a random source
// reader.
type UUID struct {
	randSrc io.Reader
}

// NewUUID returns an initialized UUID value that can be used to retrieve
// random UUID version 4 values.
func NewUUID(r io.Reader) *UUID {
	return &UUID{randSrc: r}
}

// GetUUID returns a random UUID version 4 string representation sourced from the random reader the
// UUID was created with. Returns an error if unable to compute the UUID.
func (r *UUID) GetUUID() (string, error) {
	var b [16]byte
	if _, err := io.ReadFull(r.randSrc, b[:]); err != nil {
		return "", err
	}
	r.makeUUIDv4(b[:])
	return uuid.Format(b), nil
}

// GetBytes returns a byte slice containing a random UUID version 4 sourced from the random reader the
// UUID was created with. Returns an error if unable to compute the UUID.
func (r *UUID) GetBytes() (u []byte, err error) {
	u = make([]byte, 16)
	if _, err = io.ReadFull(r.randSrc, u); err != nil {
		return u, err
	}
	r.makeUUIDv4(u)
	return u, nil
}

func (r *UUID) makeUUIDv4(u []byte) {
	// 13th character is "4"
	u[6] = (u[6] & 0x0f) | 0x40 // Version 4
	// 17th character is "8", "9", "a", or "b"
	u[8] = (u[8] & 0x3f) | 0x80 // Variant most significant bits are 10x where x can be either 1 or 0
}
