package shared

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"
	"math/rand/v2"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

func NewUUID() pgtype.UUID {
	g := pool.Get().(*generator)
	defer pool.Put(g)

	return g.newUUID()
}

var pool = sync.Pool{
	New: func() any {
		return newGenerator()
	},
}

type generator struct {
	rng rand.ChaCha8
}

func seed() [32]byte {
	var r [32]byte
	if _, err := io.ReadFull(crand.Reader, r[:]); err != nil {
		panic(err)
	}

	return r
}

func newGenerator() *generator {
	return &generator{
		rng: *rand.NewChaCha8(seed()),
	}
}

func (g *generator) newUUID() pgtype.UUID {
	var u pgtype.UUID
	u.Valid = true
	binary.NativeEndian.PutUint64(u.Bytes[:8], g.rng.Uint64())
	binary.NativeEndian.PutUint64(u.Bytes[8:], g.rng.Uint64())
	u.Bytes[6] = (u.Bytes[6] & 0x0f) | 0x40 // Version 4
	u.Bytes[8] = (u.Bytes[8] & 0x3f) | 0x80 // Variant 10

	return u
}

func EncodeUUID(uuid pgtype.UUID) string {
	src := uuid.Bytes

	var buf [36]byte

	hex.Encode(buf[0:8], src[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], src[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], src[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], src[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], src[10:])

	return string(buf[:])
}
