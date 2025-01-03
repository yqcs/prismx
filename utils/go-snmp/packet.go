package go_snmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type SnmpVersion uint8

const (
	Version1  SnmpVersion = 0x0
	Version2c SnmpVersion = 0x1
)

func (s SnmpVersion) String() string {
	if s == Version1 {
		return "1"
	} else if s == Version2c {
		return "2c"
	}
	return "U"
}

type SnmpPacket struct {
	Version        SnmpVersion
	Community      string
	RequestType    Asn1BER
	RequestID      uint32
	Error          uint8
	ErrorIndex     uint8
	NonRepeaters   uint8
	MaxRepetitions uint8
	Variables      []SnmpPDU
}

type SnmpPDU struct {
	Name  string
	Type  Asn1BER
	Value interface{}
}

func Unmarshal(packet []byte) (*SnmpPacket, error) {

	//var err error
	response := new(SnmpPacket)
	response.Variables = make([]SnmpPDU, 0, 5)

	// Start parsing the packet
	var cursor uint64 = 0

	// First bytes should be 0x30
	if Asn1BER(packet[0]) == Sequence {
		// Parse packet length
		ber, err := parseField(packet)

		if err != nil {
			return nil, err
		}

		cursor += ber.HeaderLength
		// Parse SNMP Version
		rawVersion, err := parseField(packet[cursor:])

		if err != nil {
			return nil, fmt.Errorf("Error parsing SNMP packet version: %s", err.Error())
		}

		cursor += rawVersion.DataLength + rawVersion.HeaderLength
		if version, ok := rawVersion.BERVariable.Value.(int); ok {
			response.Version = SnmpVersion(version)
		}

		// Parse community
		rawCommunity, err := parseField(packet[cursor:])

		cursor += rawCommunity.DataLength + rawCommunity.HeaderLength

		if community, ok := rawCommunity.BERVariable.Value.(string); ok {
			response.Community = community
		}

		rawPDU, err := parseField(packet[cursor:])

		response.RequestType = rawPDU.Type

		switch rawPDU.Type {
		default:
		case GetRequest, GetResponse, GetBulkRequest:
			cursor += rawPDU.HeaderLength

			// Parse Request ID
			rawRequestId, err := parseField(packet[cursor:])

			if err != nil {
				return nil, err
			}

			cursor += rawRequestId.DataLength + rawRequestId.HeaderLength
			if requestid, ok := rawRequestId.BERVariable.Value.(int); ok {
				response.RequestID = uint32(requestid)
			}

			// Parse Error
			rawError, err := parseField(packet[cursor:])

			if err != nil {
				return nil, err
			}

			cursor += rawError.DataLength + rawError.HeaderLength
			if errorNo, ok := rawError.BERVariable.Value.(int); ok {
				response.Error = uint8(errorNo)
			}

			// Parse Error Index
			rawErrorIndex, err := parseField(packet[cursor:])

			if err != nil {
				return nil, err
			}

			cursor += rawErrorIndex.DataLength + rawErrorIndex.HeaderLength

			if errorindex, ok := rawErrorIndex.BERVariable.Value.(int); ok {
				response.ErrorIndex = uint8(errorindex)
			}

			rawResp, err := parseField(packet[cursor:])

			if err != nil {
				return nil, err
			}

			cursor += rawResp.HeaderLength
			// Loop & parse Varbinds
			for cursor < uint64(len(packet)) {

				rawVarbind, err := parseField(packet[cursor:])

				if err != nil {
					return nil, err
				}

				cursor += rawVarbind.HeaderLength
				// Parse OID
				rawOid, err := parseField(packet[cursor:])

				if err != nil {
					return nil, err
				}

				cursor += rawOid.HeaderLength + rawOid.DataLength

				rawValue, err := parseField(packet[cursor:])

				if err != nil {
					return nil, err
				}
				cursor += rawValue.HeaderLength + rawValue.DataLength

				if oid, ok := rawOid.BERVariable.Value.([]int); ok {
					response.Variables = append(response.Variables, SnmpPDU{oidToString(oid), rawValue.Type, rawValue.BERVariable.Value})
				}
			}

		}
	} else {
		return nil, fmt.Errorf("Invalid packet header\n")
	}

	return response, nil
}

type RawBER struct {
	Type         Asn1BER
	HeaderLength uint64
	DataLength   uint64
	Data         []byte
	BERVariable  *Variable
}

// Parses a given field, return the ASN.1 BER Type, its header length and the data
func parseField(data []byte) (*RawBER, error) {
	var err error

	if len(data) == 0 {
		return nil, fmt.Errorf("Unable to parse BER: Data length 0")
	}

	ber := new(RawBER)

	ber.Type = Asn1BER(data[0])

	// Parse Length
	length := data[1]

	// Check if this is padded or not
	if length > 0x80 {
		length = length - 0x80
		ber.DataLength = Uvarint(data[2 : 2+length])

		ber.HeaderLength = 2 + uint64(length)

	} else {
		ber.HeaderLength = 2
		ber.DataLength = uint64(length)
	}

	// Do sanity checks
	if ber.DataLength > uint64(len(data)) {
		return nil, fmt.Errorf("Unable to parse BER: provided data length is longer than actual data (%d vs %d)", ber.DataLength, len(data))
	}

	ber.Data = data[ber.HeaderLength : ber.HeaderLength+ber.DataLength]

	ber.BERVariable, err = decodeValue(ber.Type, ber.Data)

	if err != nil {
		return nil, fmt.Errorf("Unable to decode value: %s\n", err.Error())
	}

	return ber, nil
}

func (packet *SnmpPacket) marshal() ([]byte, error) {
	// Marshal the SNMP PDU
	snmpPduBuffer := make([]byte, 0, 1024)
	snmpPduBuf := bytes.NewBuffer(snmpPduBuffer)

	requestIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(requestIDBytes, packet.RequestID)
	snmpPduBuf.Write(append([]byte{byte(packet.RequestType), 0, 2, 4}, requestIDBytes...))

	switch packet.RequestType {
	case GetBulkRequest:
		snmpPduBuf.Write([]byte{
			2, 1, packet.NonRepeaters,
			2, 1, packet.MaxRepetitions,
		})
	default:
		snmpPduBuf.Write([]byte{
			2, 1, packet.Error,
			2, 1, packet.ErrorIndex,
		})
	}

	snmpPduBuf.Write([]byte{byte(Sequence), 0})

	pduLength := 0
	for _, varlist := range packet.Variables {
		pdu, err := marshalPDU(&varlist)

		if err != nil {
			return nil, err
		}
		pduLength += len(pdu)
		snmpPduBuf.Write(pdu)
	}

	pduBytes := snmpPduBuf.Bytes()
	// Varbind list length
	pduBytes[15] = byte(pduLength)
	// SNMP PDU length (PDU header + varbind list length)
	pduBytes[1] = byte(pduLength + 14)

	// Prepare the buffer to send
	buffer := make([]byte, 0, 1024)
	buf := bytes.NewBuffer(buffer)

	// Write the message type 0x30
	buf.Write([]byte{byte(Sequence)})

	// Find the size of the whole data
	// The extra 5 bytes are snmp verion (3 bytes) + community string type and
	// community string length
	dataLength := len(pduBytes) + len(packet.Community) + 5

	// If the data is 128 bytes or larger then we need to use 2 or more bytes
	// to represent the length
	if dataLength >= 128 {
		// Work out how many bytes we require
		bytesNeeded := int(math.Ceil(math.Log2(float64(dataLength)) / 8))
		// Set the most significant bit to 1 to show we are using the long for,
		// then the 7 least significant bits to show how many bytes will be used
		// to represent the length
		lengthIdentifier := 128 + bytesNeeded
		buf.Write([]byte{uint8(lengthIdentifier)})
		lengthBytes := make([]byte, bytesNeeded)
		for i := bytesNeeded; i >= 1; i-- {
			lengthBytes[i-1] = uint8(dataLength % 256)
			dataLength = dataLength >> 8
		}
		buf.Write(lengthBytes)
	} else {
		buf.Write([]byte{uint8(dataLength)})
	}

	// Write the Version
	buf.Write([]byte{2, 1, byte(packet.Version)})

	// Write Community
	buf.Write([]byte{4, uint8(len(packet.Community))})
	buf.WriteString(packet.Community)

	// Write the PDU
	buf.Write(pduBytes)

	// Write the
	//buf.Write([]byte{packet.RequestType, uint8(17 + len(mOid)), 2, 1, 1, 2, 1, 0, 2, 1, 0, 0x30, uint8(6 + len(mOid)), 0x30, uint8(4 + len(mOid)), 6, uint8(len(mOid))})
	//buf.Write(mOid)
	//buf.Write([]byte{5, 0})

	return buf.Bytes(), nil
}

func marshalPDU(pdu *SnmpPDU) ([]byte, error) {
	oid, err := marshalOID(pdu.Name)
	if err != nil {
		return nil, err
	}

	pduBuffer := make([]byte, 0, 1024)
	pduBuf := bytes.NewBuffer(pduBuffer)

	// Mashal the PDU type into the appropriate BER
	switch pdu.Type {
	case Null:
		pduBuf.Write([]byte{byte(Sequence), byte(len(oid) + 4)})
		pduBuf.Write([]byte{byte(ObjectIdentifier), byte(len(oid))})
		pduBuf.Write(oid)
		pduBuf.Write([]byte{Null, 0x00})
	default:
		return nil, fmt.Errorf("Unable to marshal PDU: unknown BER type %d", pdu.Type)
	}

	return pduBuf.Bytes(), nil
}

func oidToString(oid []int) (ret string) {
	values := make([]interface{}, len(oid))
	for i, v := range oid {
		values[i] = v
	}
	return fmt.Sprintf(strings.Repeat(".%d", len(oid)), values...)
}

func marshalOID(oid string) ([]byte, error) {
	var err error

	// Encode the oid
	oid = strings.Trim(oid, ".")
	oidParts := strings.Split(oid, ".")
	oidBytes := make([]int, len(oidParts))

	// Convert the string OID to an array of integers
	for i := 0; i < len(oidParts); i++ {
		oidBytes[i], err = strconv.Atoi(oidParts[i])
		if err != nil {
			return nil, fmt.Errorf("Unable to parse OID: %s\n", err.Error())
		}
	}

	mOid, err := marshalObjectIdentifier(oidBytes)

	if err != nil {
		return nil, fmt.Errorf("Unable to marshal OID: %s\n", err.Error())
	}

	return mOid, err
}
