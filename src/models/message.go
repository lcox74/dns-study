package models

import "encoding/binary"

type (
	OpCode       uint8
	ResponseCode uint8
	QType        uint16
	QClass       uint16
)

type DnsMessage struct {
	Hdr         DnsHeader     `json:"hdr"`    // Required for all Messages (has length of below record arrays)
	Questions   []DnsQuestion `json:"quests"` // The question/s for the Name Server
	Answers     []DnsRR       `json:"ans"`    // Resource Records that answer the questions
	Authorities []DnsRR       `json:"auths"`  // Resource Records pointing towards an authority
	Additionals []DnsRR       `json:"adds"`   // Resource Records holdign additional information
}

type DnsHeader struct {
	ID      uint16       `json:"trans-id"`   // Identifier generated to match outstanding queries
	QR      uint8        `json:"trans-type"` // 1 Bit determining; query (0) or response (1)
	OPCODE  OpCode       `json:"opcode"`     // 4 Bit determining type of query
	AA      uint8        `json:"aa"`         // 1 Bit Authoritive Answer
	TC      uint8        `json:"tc"`         // 1 Bit Truncation; specifies if the message is too long
	RD      uint8        `json:"rd"`         // 1 Bit Recursion Desired
	RA      uint8        `json:"ra"`         // 1 Bit Recursion Avaliable
	Z       uint8        `json:"z"`          // 1 Bit Reserved for future use (must be 0)
	RCODE   ResponseCode `json:"rcode"`      // Response Code
	QDCOUNT uint16       `json:"qdcount"`    // Number of entries in Question Section
	ADCOUNT uint16       `json:"adcount"`    // Number of RR in Answer Section
	NSCOUNT uint16       `json:"nscount"`    // Number of NS RR in Authoritive Records Section
	ARCOUNT uint16       `json:"arcount"`    // Number of RR in Additional Records Section
}

const HdrLen uint8 = 12 // Length of DnsHeader in Bytes

type QNameLabel struct {
	Length byte   `json:"len"`   // Length of Data in bytes
	Data   []byte `json:"label"` // Label Data
}

type DnsQuestion struct {
	QNAME  []QNameLabel `json:"name"`  // Domain name represented as a sequence of labels. Ends with an empty length label
	QTYPE  QType        `json:"type"`  // Specifies the Type of query eg. NS for an authoritative name server
	QCLASS QClass       `json:"class"` // Specifies the Class of query eg. IN for internet
}

type DnsRR struct {
	NAME     []byte `json:"name"`  // Uses Message Compression (See RFC1035 4.1.4)
	TYPE     QType  `json:"type"`  // Specifies the Type of query eg. NS for an authoritative name server
	CLASS    QClass `json:"class"` // Specifies the Class of query eg. IN for internet
	TTL      uint32 `json:"ttl"`   // Time interval in seconds RR can be cached before needing to be discarded. 0 = Don't cache
	RDLENGTH uint16 `json:"rdlen"` // Length in bytes of RDATA
	RDATA    []byte `json:"rdata"` // Data of RR, format determined by TYPE and CLASS
}

const (
	QUERY  OpCode = 0 // Standard Query
	IQUERY OpCode = 1 // Inverse Query
	STATUS OpCode = 2 // Server Status Request
)

const (
	OK             ResponseCode = 0 // No error condition
	FormatError    ResponseCode = 1 // NS was unable to interpret query
	ServerFailure  ResponseCode = 2 // NS was unable to process query due to problem with NS
	NameError      ResponseCode = 3 // Domain name referenced does not exit
	NotImplemented ResponseCode = 4 // NS does not support hte requested kind of query
	Refused        ResponseCode = 5 // NS refuses to perform the specified operation for policy reasons
)

const (
	A     QType = 1   // Host Address
	NS    QType = 2   // Authoritative Name Server
	MD    QType = 3   // Mail Destination 					 (Obsolete - use MX)
	MF    QType = 4   // Mail Forwarder 					 (Obsolete - use MX)
	CNAME QType = 5   // Canonical Name for an Alias
	SOA   QType = 6   // Marks the start of a zone of Authority
	MB    QType = 7   // Mailbox domain name 					  (EXPERIMENTAL)
	MG    QType = 8   // Mail group member 						  (EXPERIMENTAL)
	MR    QType = 9   // Mail rename domain name 				  (EXPERIMENTAL)
	NULL  QType = 10  // Null RR 								  (EXPERIMENTAL)
	WKS   QType = 11  // Well known service description
	PTR   QType = 12  // Domain name pointer
	HINFO QType = 13  // Host information
	MINFO QType = 14  // Mailbox or Mail list information
	MX    QType = 15  // Mail exchange
	TXT   QType = 16  // Text strings
	AXFR  QType = 252 // Request for a transfer of an entire zone
	MAILB QType = 253 // Request for mailbox-related records (MB, MG or MR)
	MAILA QType = 254 // Request for mail agent RR			 (Obsolete - use MX)
	ALL   QType = 255 // Request for all records
)

const (
	IN QClass = 1 // The Internet
	CS QClass = 2 // The CSNET class (Obsolete - used only for examples in some obsolete RFCs)
	CH QClass = 3 // The CHAOS class
	HS QClass = 4 // Hesiod [Dyer 87]
)

func MarshalDNS(data []byte) (DnsMessage, error) {
	var (
		message DnsMessage
		offset  uint32 = 0
		err     error
	)

	// Marshal Header
	message.Hdr, err = marshalDNSHdr(data[offset:HdrLen])
	if err != nil {
		return message, err
	}
	offset += uint32(HdrLen)

	// Marshal Questions
	if message.Hdr.QDCOUNT > 0 {
		var off uint32
		message.Questions, off, err = marshalDnsQuestions(data[offset:], message.Hdr.QDCOUNT)
		if err != nil {
			return message, err
		}
		offset += off
	}

	return message, nil
}

func marshalDNSHdr(data []byte) (DnsHeader, error) {
	var hdr DnsHeader

	// Transaction ID
	hdr.ID = binary.BigEndian.Uint16(data[0:2])

	// Flags
	hdr.QR = (data[2] & 0x80) >> 7
	hdr.OPCODE = OpCode((data[3] & 0x78) >> 3)
	hdr.AA = (data[2] & 0x04) >> 2
	hdr.TC = (data[2] & 0x02) >> 1
	hdr.RD = (data[2] & 0x01) >> 0
	hdr.RA = (data[3] & 0x80) >> 7
	hdr.Z = (data[3] & 0x70) >> 4
	hdr.RCODE = ResponseCode(data[3]&0x0F) >> 0

	// RR Counts
	hdr.QDCOUNT = binary.BigEndian.Uint16(data[4:6])
	hdr.ADCOUNT = binary.BigEndian.Uint16(data[6:8])
	hdr.NSCOUNT = binary.BigEndian.Uint16(data[8:10])
	hdr.ARCOUNT = binary.BigEndian.Uint16(data[10:12])

	return hdr, nil
}

func marshalDnsQuestions(data []byte, num uint16) ([]DnsQuestion, uint32, error) {
	var (
		questions []DnsQuestion
		offset    uint32 = 0
	)

	for iter := uint32(0); iter < uint32(num); iter++ {
		question, off, err := marshalDnsQuestion(data[offset:])
		if err != nil {
			return questions, offset, err
		}
		offset += off
		questions = append(questions, question)
	}

	return questions, offset, nil
}

func marshalDnsQuestion(data []byte) (DnsQuestion, uint32, error) {
	var (
		question DnsQuestion
		offset   uint32 = 0
		off      uint32 = 0
	)

	// Get Name
	question.QNAME, off = marshalNameLabels(data[offset:])
	offset += off

	// Get Type
	question.QTYPE = QType(binary.BigEndian.Uint16(data[offset : offset+2]))
	offset += 2

	// Get Class
	question.QCLASS = QClass(binary.BigEndian.Uint16(data[offset : offset+2]))
	offset += 2

	return question, offset, nil
}

func marshalNameLabels(data []byte) ([]QNameLabel, uint32) {
	var (
		labels []QNameLabel
		offset uint32 = 0
	)

	for data[offset] != 0x00 {
		label, off := marshalNameLabel(data[offset:])
		labels = append(labels, label)
		offset += off
	}

	return labels, offset
}

func marshalNameLabel(data []byte) (QNameLabel, uint32) {
	var label QNameLabel

	label.Length = data[0]
	label.Data = data[1 : label.Length+1]

	return label, uint32(label.Length + 1)
}
