package clipperz

import (
	"crypto/sha256"
	"encoding/hex"

	crand "crypto/rand"
)

func PrefixMatchingBits(a, b []byte) (result int) {
	c := len(a)
	if len(b) < len(a) {
		c = len(b)
	}

	i := 0
	for i < c && a[i] == b[i] {
		result += 8
		i++
	}

	if i < c {
		xor := a[i] ^ b[i]

		switch {
		case xor >= 128:
			result += 0
		case xor >= 64:
			result += 1
		case xor >= 32:
			result += 2
		case xor >= 16:
			result += 3
		case xor >= 8:
			result += 4
		case xor >= 4:
			result += 5
		case xor >= 2:
			result += 6
		case xor >= 1:
			result += 7
		}
	}

	return result
}

type Toll struct {
	RequestType string
	TargetValue []byte
	Cost        int
	Toll        []byte
}

func NewToll(requestType string, targetValue []byte, cost int) (toll *Toll) {
	return &Toll{
		RequestType: requestType,
		TargetValue: targetValue,
		Cost:        cost,
	}
}

func (toll *Toll) Pay() (err error) {
	payment := make([]byte, 32)
	_, err = crand.Reader.Read(payment)
	if err != nil {
		return err
	}

	for {
		for i := 31; i >= 0; i-- {
			payment[i]++
			if payment[i] != 0 {
				break
			}
		}

		paymentHash := sha256.New()
		paymentHash.Write(payment)
		paymentData := paymentHash.Sum(nil)

		prefixMatchingBits := PrefixMatchingBits(toll.TargetValue, paymentData)
		if prefixMatchingBits >= toll.Cost {
			break
		}
	}

	toll.Toll = payment
	return nil
}

func (toll *Toll) RequestData() (data map[string]string) {
	return map[string]string{
		"targetValue": hex.EncodeToString(toll.TargetValue),
		"toll":        hex.EncodeToString(toll.Toll),
	}
}
