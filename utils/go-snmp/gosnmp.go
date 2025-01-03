// Copyright 2012 Andreas Louca. All rights reserved.
// Use of this source code is goverend by a BSD-style
// license that can be found in the LICENSE file.

package go_snmp

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

// GoSNMP represents the GoSNMP poller structure
type GoSNMP struct {
	Target    string
	Community string
	Version   SnmpVersion
	Timeout   time.Duration
	Conn      net.Conn
}

// StreamWalk will start walking a specified OID, and push through a channel the results
// as it receives them, without waiting for the whole process to finish to return the
// results. Once it has completed the walk, the channel is closed.
func (x *GoSNMP) StreamWalk(oid string, c chan SnmpPDU) error {
	if oid == "" {
		close(c)
		return fmt.Errorf("No OID given\n")
	}

	requestOid := oid

	for {
		res, err := x.GetNext(oid)
		if err != nil {
			close(c)
			return err
		}
		if res != nil {
			if len(res.Variables) > 0 {
				if strings.Index(res.Variables[0].Name, requestOid) > -1 {
					if res.Variables[0].Value == "endOfMib" {
						break
					}
					c <- res.Variables[0]
					// Set to the next
					oid = res.Variables[0].Name
				} else {
					break
				}
			} else {
				break
			}
		} else {
			break
		}

	}
	close(c)
	return nil
}

// BulkWalk sends an walks the target using SNMP BULK-GET requests. This returns
// a Variable with the response and the error condition
func (x *GoSNMP) BulkWalk(maxRepetitions uint8, oid string) (results []SnmpPDU, err error) {
	if oid == "" {
		return nil, fmt.Errorf("No OID given\n")
	}
	return x._bulkWalk(maxRepetitions, oid, oid)
}
func (x *GoSNMP) _bulkWalk(maxRepetitions uint8, searchingOid string, rootOid string) (results []SnmpPDU, err error) {
	response, err := x.GetBulk(0, maxRepetitions, searchingOid)
	if err != nil {
		return
	}
	for i, v := range response.Variables {
		if v.Value == "endOfMib" {
			return
		}
		// is this variable still in the requested oid range
		if strings.HasPrefix(v.Name, rootOid) {
			results = append(results, v)
			// is the last oid received still in the requested range
			if i == len(response.Variables)-1 {
				var subResults []SnmpPDU
				subResults, err = x._bulkWalk(maxRepetitions, v.Name, rootOid)
				if err != nil {
					return
				}
				results = append(results, subResults...)
			}
		}
	}
	return
}

// Walk will SNMP walk the target, blocking until the process is complete
func (x *GoSNMP) Walk(oid string) (results []SnmpPDU, err error) {
	if oid == "" {
		return nil, fmt.Errorf("No OID given\n")
	}
	results = make([]SnmpPDU, 0)
	requestOid := oid

	for {
		res, err := x.GetNext(oid)
		if err != nil {
			return results, err
		}
		if res != nil {
			if len(res.Variables) > 0 {
				if strings.Index(res.Variables[0].Name, requestOid) > -1 {
					results = append(results, res.Variables[0])
					// Set to the next
					oid = res.Variables[0].Name
				} else {
					break
				}
			} else {
				break
			}
		} else {
			break
		}

	}
	return
}

// sendPacket marshals & send an SNMP request. Unmarshals the response and
// returns back the parsed SNMP packet
func (x *GoSNMP) sendPacket(packet *SnmpPacket) (*SnmpPacket, error) {
	// Set timeouts on the connection
	deadline := time.Now()
	x.Conn.SetDeadline(deadline.Add(x.Timeout))

	// Create random Request-ID
	packet.RequestID = rand.Uint32()

	// Marshal it
	fBuf, err := packet.marshal()

	if err != nil {
		return nil, err
	}

	// Send the packet!
	_, err = x.Conn.Write(fBuf)
	if err != nil {
		return nil, fmt.Errorf("Error writing to socket: %s\n", err.Error())
	}
	// Try to read the response
	resp := make([]byte, 8192, 8192)
	n, err := x.Conn.Read(resp)

	if err != nil {
		return nil, fmt.Errorf("Error reading from UDP: %s\n", err.Error())
	}

	// Unmarshal the read bytes
	pdu, err := Unmarshal(resp[:n])

	if err != nil {
		return nil, fmt.Errorf("Unable to decode packet: %s\n", err.Error())
	}

	if len(pdu.Variables) < 1 {
		return nil, fmt.Errorf("No responses received.")
	}

	// check Request-ID
	if pdu.RequestID != packet.RequestID {
		return nil, fmt.Errorf("Request ID mismatch")
	}

	return pdu, nil
}

// GetNext sends an SNMP Get Next Request to the target. Returns the next
// variable response from the OID given or an error
func (x *GoSNMP) GetNext(oid string) (*SnmpPacket, error) {
	return x.request(GetNextRequest, oid)
}

// Debug function. Unmarshals raw bytes and returns the result without the network part
func (x *GoSNMP) Debug(data []byte) (*SnmpPacket, error) {
	packet, err := Unmarshal(data)

	if err != nil {
		return nil, fmt.Errorf("Unable to decode packet: %s\n", err.Error())
	}
	return packet, nil
}

// GetBulk sends an SNMP BULK-GET request to the target. Returns a Variable with
// the response or an error
func (x *GoSNMP) GetBulk(nonRepeaters, maxRepetitions uint8, oids ...string) (*SnmpPacket, error) {
	// Create and send the packet
	return x.sendPacket(&SnmpPacket{
		Version:        x.Version,
		Community:      x.Community,
		RequestType:    GetBulkRequest,
		NonRepeaters:   nonRepeaters,
		MaxRepetitions: maxRepetitions,
		Variables:      oidsToPbus(oids...),
	})
}

// Get sends an SNMP GET request to the target. Returns a Variable with the
// response or an error
func (x *GoSNMP) Get(oid string) (*SnmpPacket, error) {
	return x.request(GetRequest, oid)
}

// GetMulti sends an SNMP GET request to the target. Returns a Variable with the
// response or an error
func (x *GoSNMP) GetMulti(oids []string) (*SnmpPacket, error) {
	return x.request(GetRequest, oids...)
}

func (x *GoSNMP) request(requestType Asn1BER, oids ...string) (*SnmpPacket, error) {
	// Create and send the packet
	return x.sendPacket(&SnmpPacket{
		Version:     x.Version,
		Community:   x.Community,
		RequestType: requestType,
		Variables:   oidsToPbus(oids...),
	})
}

func oidsToPbus(oids ...string) []SnmpPDU {
	pdus := make([]SnmpPDU, len(oids))
	for i, oid := range oids {
		pdus[i] = SnmpPDU{Name: oid, Type: Null}
	}
	return pdus
}
