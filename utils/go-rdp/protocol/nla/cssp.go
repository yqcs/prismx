package nla

import (
	"encoding/asn1"
)

type NegoToken struct {
	Data []byte `asn1:"explicit,tag:0"`
}

type TSRequest struct {
	Version    int         `asn1:"explicit,tag:0"`
	NegoTokens []NegoToken `asn1:"optional,explicit,tag:1"`
	AuthInfo   []byte      `asn1:"optional,explicit,tag:2"`
	PubKeyAuth []byte      `asn1:"optional,explicit,tag:3"`
	//ErrorCode  int         `asn1:"optional,explicit,tag:4"`
}

type TSCredentials struct {
	CredType    int    `asn1:"explicit,tag:0"`
	Credentials []byte `asn1:"explicit,tag:1"`
}

type TSPasswordCreds struct {
	DomainName []byte `asn1:"explicit,tag:0"`
	UserName   []byte `asn1:"explicit,tag:1"`
	Password   []byte `asn1:"explicit,tag:2"`
}

type TSCspDataDetail struct {
	KeySpec       int    `asn1:"explicit,tag:0"`
	CardName      string `asn1:"explicit,tag:1"`
	ReaderName    string `asn1:"explicit,tag:2"`
	ContainerName string `asn1:"explicit,tag:3"`
	CspName       string `asn1:"explicit,tag:4"`
}

type TSSmartCardCreds struct {
	Pin        string            `asn1:"explicit,tag:0"`
	CspData    []TSCspDataDetail `asn1:"explicit,tag:1"`
	UserHint   string            `asn1:"explicit,tag:2"`
	DomainHint string            `asn1:"explicit,tag:3"`
}

func EncodeDERTRequest(msgs []Message, authInfo []byte, pubKeyAuth []byte) []byte {
	req := TSRequest{
		Version: 2,
	}

	if len(msgs) > 0 {
		req.NegoTokens = make([]NegoToken, 0, len(msgs))
	}

	for _, msg := range msgs {
		token := NegoToken{msg.Serialize()}
		req.NegoTokens = append(req.NegoTokens, token)
	}

	if len(authInfo) > 0 {
		req.AuthInfo = authInfo
	}

	if len(pubKeyAuth) > 0 {
		req.PubKeyAuth = pubKeyAuth
	}

	result, _ := asn1.Marshal(req)

	return result
}

func DecodeDERTRequest(s []byte) (*TSRequest, error) {
	treq := &TSRequest{}
	_, err := asn1.Unmarshal(s, treq)
	return treq, err
}
func EncodeDERTCredentials(domain, username, password []byte) []byte {
	tpas := TSPasswordCreds{domain, username, password}
	result, _ := asn1.Marshal(tpas)

	tcre := TSCredentials{1, result}
	result, _ = asn1.Marshal(tcre)

	return result
}

func DecodeDERTCredentials(s []byte) (*TSCredentials, error) {
	tcre := &TSCredentials{}
	_, err := asn1.Unmarshal(s, tcre)
	return tcre, err
}