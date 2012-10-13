package clipperz

import (
	"encoding/hex"
	"fmt"
	"hash"
	"math/big"

	crand "crypto/rand"
)

var n *big.Int
var g *big.Int

func init() {
	n_, ok := big.NewInt(0).SetString("115b8b692e0e045692cf280b436735c77a5a9e8a9e7ed56c965f87db5b2a2ece3", 16)
	if !ok {
		panic("Could not initialise SRP n")
	}

	n = n_
	g = big.NewInt(2)
}

type SRPConnection struct {
	C, P                []byte
	H                   hash.Hash
	a, A, s, B, x, u, S *big.Int
	K, M1, M2           string
}

func NewSRPConnection(C []byte, P []byte, H hash.Hash) (srp *SRPConnection, err error) {
	secretBytes := make([]byte, 32)
	_, err = crand.Reader.Read(secretBytes)
	if err != nil {
		return nil, err
	}

	secretHex := hex.EncodeToString(secretBytes)
	a, ok := big.NewInt(0).SetString(secretHex, 16)
	if !ok {
		return nil, fmt.Errorf("Could not initialise SRP a")
	}

	A := big.NewInt(0).Exp(g, a, n)

	srp = &SRPConnection{
		C: C,
		P: P,
		H: H,
		a: a,
		A: A,
	}

	return srp, nil
}

func (srp *SRPConnection) SetResponseData(s, B *big.Int) {
	srp.s = s
	srp.B = B

	var ok bool

	srp.x, ok = big.NewInt(0).SetString(srp.stringHash(
		fmt.Sprintf("%064X", srp.s),
		hex.EncodeToString(srp.P),
	), 16)

	if !ok {
		panic("Could not initialise SRP x")
	}

	srp.u, ok = big.NewInt(0).SetString(srp.stringHash(
		fmt.Sprintf("%d", srp.B),
	), 16)

	if !ok {
		panic("Could not initialise SRP u")
	}

	srp.S = big.NewInt(0).Exp(
		big.NewInt(0).Sub(
			srp.B,
			big.NewInt(0).Exp(g, srp.x, n),
		),

		big.NewInt(0).Add(
			srp.a,
			big.NewInt(0).Mul(srp.u, srp.x),
		),

		n,
	)

	srp.K = srp.stringHash(
		fmt.Sprintf("%d", srp.S),
	)

	srp.M1 = srp.stringHash(
		fmt.Sprintf("%d", srp.A),
		fmt.Sprintf("%d", srp.B),
		srp.K,
	)

	srp.M2 = srp.stringHash(
		fmt.Sprintf("%d", srp.A),
		srp.M1,
		srp.K,
	)
}

func (srp *SRPConnection) stringHash(inputs ...string) (sum string) {
	srp.H.Reset()

	for _, input := range inputs {
		srp.H.Write([]byte(input))
	}

	sum = hex.EncodeToString(srp.H.Sum(nil))
	srp.H.Reset()

	return sum
}
