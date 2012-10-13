package clipperz

import (
	"crypto/sha256"
)

type Client struct {
	Backend           *Backend
	ConnectionID      string
	LoginInfo         *LoginInfo
	OfflineCopyNeeded bool
	Lock              string
}

func NewClient(username string, passphrase string) (client *Client, err error) {
	h := sha256.New()

	h.Write([]byte(username))
	h.Write([]byte(passphrase))
	C := h.Sum(nil)
	h.Reset()

	h.Write([]byte(passphrase))
	h.Write([]byte(username))
	P := h.Sum(nil)
	h.Reset()

	srp, err := NewSRPConnection(C, P, h)
	if err != nil {
		return nil, err
	}

	backend := NewBackend("http://clipperz.com/beta/json/", false, srp)
	connectionID, loginInfo, offlineCopyNeeded, lock, err := backend.Login()
	if err != nil {
		return nil, err
	}

	client = &Client{
		Backend:           backend,
		ConnectionID:      connectionID,
		LoginInfo:         loginInfo,
		OfflineCopyNeeded: offlineCopyNeeded,
		Lock:              lock,
	}

	return client, nil
}
