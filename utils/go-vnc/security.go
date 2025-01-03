// Implementation of RFC 6143 §7.2 Security Types.

package vnc

import (
	"crypto/des"
	"errors"
)

const (
	secTypeInvalid = uint8(0)
	secTypeNone    = uint8(1)
	secTypeVNCAuth = uint8(2)
)

// ClientAuth implements a method of authenticating with a remote server.
type ClientAuth interface {
	// SecurityType returns the byte identifier sent by the server to
	// identify this authentication scheme.
	SecurityType() uint8

	// Handshake is called when the authentication handshake should be
	// performed, as part of the general RFB handshake. (see 7.2.1)
	Handshake(*ClientConn) error
}

// ClientAuthNone is the "none" authentication. See 7.2.1.
type ClientAuthNone struct{}

func (*ClientAuthNone) SecurityType() uint8 {
	return secTypeNone
}

func (*ClientAuthNone) Handshake(conn *ClientConn) error {

	return nil
}

// ClientAuthVNC is the standard password authentication. See 7.2.2.
type ClientAuthVNC struct {
	Password string
}

type vncAuthChallenge [16]byte

func (*ClientAuthVNC) SecurityType() uint8 {
	return secTypeVNCAuth
}

func (auth *ClientAuthVNC) Handshake(conn *ClientConn) error {

	if auth.Password == "" {
		return errors.New("security Handshake failed; no password provided for VNCAuth")
	}

	// Read challenge block
	var challenge vncAuthChallenge
	if err := conn.receive(&challenge); err != nil {
		return err
	}

	auth.encode(&challenge)

	// Send the encrypted challenge back to server
	if err := conn.send(challenge); err != nil {
		return err
	}

	return nil
}

func (auth *ClientAuthVNC) encode(ch *vncAuthChallenge) error {
	// Copy password string to 8 byte 0-padded slice
	key := make([]byte, 8)
	copy(key, auth.Password)

	// Each byte of the password needs to be reversed. This is a
	// non RFC-documented behaviour of VNC clients and servers
	for i := range key {
		key[i] = (key[i]&0x55)<<1 | (key[i]&0xAA)>>1 // Swap adjacent bits
		key[i] = (key[i]&0x33)<<2 | (key[i]&0xCC)>>2 // Swap adjacent pairs
		key[i] = (key[i]&0x0F)<<4 | (key[i]&0xF0)>>4 // Swap the 2 halves
	}

	// Encrypt challenge with key.
	cipher, err := des.NewCipher(key)
	if err != nil {
		return err
	}
	for i := 0; i < len(ch); i += cipher.BlockSize() {
		cipher.Encrypt(ch[i:i+cipher.BlockSize()], ch[i:i+cipher.BlockSize()])
	}

	return nil
}
