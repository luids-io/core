package encoding

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/luisguillenc/tlslayer"
	"github.com/luisguillenc/tlslayer/tlsproto"

	pb "github.com/luids-io/core/protogen/tlsutilpb"
	"github.com/luids-io/core/tlsutil"
)

// ConnectionData convert to model connectiondata
func ConnectionData(src *pb.ConnectionData) *tlsutil.ConnectionData {
	dst := &tlsutil.ConnectionData{}
	dst.ID = src.GetId()
	info := src.GetInfo()
	if info != nil {
		dst.Info = ConnectionInfo(info)
	}
	sstream := src.GetSendStream()
	if sstream != nil {
		dst.SendStream = StreamData(sstream)
	}
	rstream := src.GetRcvdStream()
	if rstream != nil {
		dst.RcvdStream = StreamData(rstream)
	}
	chello := src.GetClientHello()
	if chello != nil {
		dst.ClientHello = ClientHelloData(chello)
	}
	shello := src.GetServerHello()
	if shello != nil {
		dst.ServerHello = ServerHelloData(shello)
	}
	ccerts := src.GetClientCerts()
	if len(ccerts) > 0 {
		dst.ClientCerts = make([]tlsutil.CertSummary, 0, len(ccerts))
		for _, c := range ccerts {
			dst.ClientCerts = append(dst.ClientCerts, CertSummary(c))
		}
	}
	scerts := src.GetServerCerts()
	if len(scerts) > 0 {
		dst.ServerCerts = make([]tlsutil.CertSummary, 0, len(scerts))
		for _, c := range scerts {
			dst.ServerCerts = append(dst.ServerCerts, CertSummary(c))
		}
	}
	stags := src.GetTags()
	if len(stags) > 0 {
		dst.Tags = make([]string, 0, len(stags))
		for _, t := range stags {
			dst.Tags = append(dst.Tags, t)
		}
	}
	return dst
}

// ConnectionDataPB convert to protobuf connectiondata
func ConnectionDataPB(src *tlsutil.ConnectionData) *pb.ConnectionData {
	dst := &pb.ConnectionData{}
	dst.Id = src.ID
	if src.Info != nil {
		dst.Info = ConnectionInfoPB(src.Info)
	}
	if src.SendStream != nil {
		dst.SendStream = StreamDataPB(src.SendStream)
	}
	if src.RcvdStream != nil {
		dst.RcvdStream = StreamDataPB(src.RcvdStream)
	}
	if src.ClientHello != nil {
		dst.ClientHello = ClientHelloDataPB(src.ClientHello)
	}
	if src.ServerHello != nil {
		dst.ServerHello = ServerHelloDataPB(src.ServerHello)
	}
	if len(src.ClientCerts) > 0 {
		dst.ClientCerts = make([]*pb.ConnectionData_CertSummary, 0, len(src.ClientCerts))
		for _, c := range src.ClientCerts {
			dst.ClientCerts = append(dst.ClientCerts, CertSummaryPB(c))
		}
	}
	if len(src.ServerCerts) > 0 {
		dst.ServerCerts = make([]*pb.ConnectionData_CertSummary, 0, len(src.ServerCerts))
		for _, c := range src.ServerCerts {
			dst.ServerCerts = append(dst.ServerCerts, CertSummaryPB(c))
		}
	}
	if len(src.Tags) > 0 {
		dst.Tags = make([]string, 0, len(src.Tags))
		for _, t := range src.Tags {
			dst.Tags = append(dst.Tags, t)
		}
	}
	return dst
}

// ConnectionInfo convert to model ConnectionInfo
func ConnectionInfo(src *pb.ConnectionData_ConnectionInfo) *tlsutil.ConnectionInfo {
	dst := &tlsutil.ConnectionInfo{}
	dst.Start, _ = ptypes.Timestamp(src.GetStart())
	dst.End, _ = ptypes.Timestamp(src.GetEnd())
	dst.Duration = time.Duration(src.GetDurationNsecs())
	dst.ClientIP = src.GetClientIp()
	dst.ClientPort = int(src.GetClientPort())
	dst.ServerIP = src.GetServerIp()
	dst.ServerPort = int(src.GetServerPort())
	dst.Uncompleted = src.GetUncompleted()
	dst.DetectedError = src.GetDetectedError()
	dst.CompletedHandshake = src.GetCompletedHandshake()
	return dst
}

// ConnectionInfoPB convert to protobuf ConnectionInfo
func ConnectionInfoPB(src *tlsutil.ConnectionInfo) *pb.ConnectionData_ConnectionInfo {
	dst := &pb.ConnectionData_ConnectionInfo{}
	dst.Start, _ = ptypes.TimestampProto(src.Start)
	dst.End, _ = ptypes.TimestampProto(src.End)
	dst.DurationNsecs = int64(src.Duration)
	dst.ClientIp = src.ClientIP
	dst.ClientPort = uint32(src.ClientPort)
	dst.ServerIp = src.ServerIP
	dst.ServerPort = uint32(src.ServerPort)
	dst.Uncompleted = src.Uncompleted
	dst.DetectedError = src.DetectedError
	dst.CompletedHandshake = src.CompletedHandshake
	return dst
}

// StreamData convert to model StreamData
func StreamData(src *pb.ConnectionData_StreamData) *tlsutil.StreamData {
	dst := &tlsutil.StreamData{}
	dst.ID = src.GetId()
	info := src.GetInfo()
	if info != nil {
		dst.Info = StreamInfo(info)
	}
	plain := src.GetPlaintextAcc()
	if plain != nil {
		dst.PlaintextAcc = PlaintextSummary(plain)
	}
	cipher := src.GetCiphertextAcc()
	if cipher != nil {
		dst.CiphertextAcc = CiphertextSummary(cipher)
	}
	hskseq := src.GetHandshakes()
	if len(hskseq) > 0 {
		dst.HandshakeSeq = make([]tlsutil.HandshakeItem, 0, len(hskseq))
		for _, hsk := range hskseq {
			dst.HandshakeSeq = append(dst.HandshakeSeq,
				tlsutil.HandshakeItem{
					Len:  hsk.GetLen(),
					Type: tlsproto.HandshakeType(hsk.GetHtype()),
				})
		}
	}
	dst.HandshakeSum = int(src.GetHandshakeSum())
	return dst
}

// StreamDataPB convert to protobuf StreamData
func StreamDataPB(src *tlsutil.StreamData) *pb.ConnectionData_StreamData {
	dst := &pb.ConnectionData_StreamData{}
	dst.Id = src.ID
	if src.Info != nil {
		dst.Info = StreamInfoPB(src.Info)
	}
	if src.PlaintextAcc != nil {
		dst.PlaintextAcc = PlaintextSummaryPB(src.PlaintextAcc)
	}
	if src.CiphertextAcc != nil {
		dst.CiphertextAcc = CiphertextSummaryPB(src.CiphertextAcc)
	}
	if len(src.HandshakeSeq) > 0 {
		dst.Handshakes = make([]*pb.ConnectionData_StreamData_HandshakeItem, 0, len(src.HandshakeSeq))
		for _, hsk := range src.HandshakeSeq {
			info := &pb.ConnectionData_StreamData_HandshakeItem{
				Htype: uint32(hsk.Type),
				Len:   hsk.Len,
			}
			dst.Handshakes = append(dst.Handshakes, info)
		}
	}
	dst.HandshakeSum = int32(src.HandshakeSum)
	return dst
}

// StreamInfo convert to model StreamInfo
func StreamInfo(src *pb.ConnectionData_StreamData_StreamInfo) *tlsutil.StreamInfo {
	dst := &tlsutil.StreamInfo{}
	dst.Start, _ = ptypes.Timestamp(src.GetStart())
	dst.End, _ = ptypes.Timestamp(src.GetEnd())
	dst.Duration = time.Duration(src.GetDurationNsecs())
	dst.SawStart = src.GetSawStart()
	dst.SawEnd = src.GetSawEnd()
	dst.DetectedError = src.GetDetectedError()
	if dst.DetectedError {
		dst.ErrorType = src.GetErrorType()
		dst.ErrorTime, _ = ptypes.Timestamp(src.GetErrorTime())
	}
	dst.SrcIP4 = src.GetSrcIp()
	dst.SrcPort = int(src.GetSrcPort())
	dst.DstIP4 = src.GetDstIp()
	dst.DstPort = int(src.GetDstPort())
	dst.Bytes = src.GetBytes()
	dst.Packets = src.GetPackets()
	dst.BPS = src.GetBps()
	dst.PPS = src.GetPps()
	return dst
}

// StreamInfoPB convert to protobuf StreamInfo
func StreamInfoPB(src *tlsutil.StreamInfo) *pb.ConnectionData_StreamData_StreamInfo {
	dst := &pb.ConnectionData_StreamData_StreamInfo{}
	dst.Start, _ = ptypes.TimestampProto(src.Start)
	dst.End, _ = ptypes.TimestampProto(src.End)
	dst.DurationNsecs = int64(src.Duration)
	dst.SawStart = src.SawStart
	dst.SawEnd = src.SawEnd
	dst.DetectedError = src.DetectedError
	if src.DetectedError {
		dst.ErrorType = src.ErrorType
		dst.ErrorTime, _ = ptypes.TimestampProto(src.ErrorTime)
	}
	dst.SrcIp = src.SrcIP4
	dst.SrcPort = uint32(src.SrcPort)
	dst.DstIp = src.DstIP4
	dst.DstPort = uint32(src.DstPort)
	dst.Bytes = src.Bytes
	dst.Packets = src.Packets
	dst.Bps = src.BPS
	dst.Pps = src.PPS
	return dst
}

// PlaintextSummary convert to model PlaintextSummary
func PlaintextSummary(src *pb.ConnectionData_StreamData_PlaintextSummary) *tlsutil.PlaintextSummary {
	dst := &tlsutil.PlaintextSummary{}
	dst.HskRecords = src.HskRecords
	dst.HskBytes = src.HskBytes
	dst.AlertRecords = src.AlertRecords
	dst.AlertBytes = src.AlertBytes
	dst.CCTRecords = src.CctRecords
	dst.CCTBytes = src.CctBytes
	dst.AppDataRecords = src.AppdataRecords
	dst.AppDataBytes = src.AppdataBytes
	dst.FragmentedRecords = int(src.FragmentedRecords)
	dst.MaxMessagesInRecord = int(src.MaxMessages)
	return dst
}

// PlaintextSummaryPB convert to protobuf PlaintextSummary
func PlaintextSummaryPB(src *tlsutil.PlaintextSummary) *pb.ConnectionData_StreamData_PlaintextSummary {
	dst := &pb.ConnectionData_StreamData_PlaintextSummary{}
	dst.HskRecords = src.HskRecords
	dst.HskBytes = src.HskBytes
	dst.AlertRecords = src.AlertRecords
	dst.AlertBytes = src.AlertBytes
	dst.CctRecords = src.CCTRecords
	dst.CctBytes = src.CCTBytes
	dst.AppdataRecords = src.AppDataRecords
	dst.AppdataBytes = src.AppDataBytes
	dst.FragmentedRecords = int32(src.FragmentedRecords)
	dst.MaxMessages = int32(src.MaxMessagesInRecord)
	return dst
}

// CiphertextSummary convert to model CiphertextSummary
func CiphertextSummary(src *pb.ConnectionData_StreamData_CiphertextSummary) *tlsutil.CiphertextSummary {
	dst := &tlsutil.CiphertextSummary{}
	dst.HskRecords = src.HskRecords
	dst.HskBytes = src.HskBytes
	dst.AlertRecords = src.AlertRecords
	dst.AlertBytes = src.AlertBytes
	dst.CCTRecords = src.CctRecords
	dst.CCTBytes = src.CctBytes
	dst.AppDataRecords = src.AppdataRecords
	dst.AppDataBytes = src.AppdataBytes
	return dst
}

// CiphertextSummaryPB convert to protobuf CiphertextSummary
func CiphertextSummaryPB(src *tlsutil.CiphertextSummary) *pb.ConnectionData_StreamData_CiphertextSummary {
	dst := &pb.ConnectionData_StreamData_CiphertextSummary{}
	dst.HskRecords = src.HskRecords
	dst.HskBytes = src.HskBytes
	dst.AlertRecords = src.AlertRecords
	dst.AlertBytes = src.AlertBytes
	dst.CctRecords = src.CCTRecords
	dst.CctBytes = src.CCTBytes
	dst.AppdataRecords = src.AppDataRecords
	dst.AppdataBytes = src.AppDataBytes
	return dst
}

// ClientHelloData convert to model ClientHelloData
func ClientHelloData(src *pb.ConnectionData_ClientHelloData) *tlsutil.ClientHelloData {
	dst := &tlsutil.ClientHelloData{}
	dst.ClientVersion = tlslayer.ProtocolVersion(src.GetClientVersion())
	dst.RandomLen = int(src.GetRandomLen())
	dst.SessionIDLen = int(src.GetSessionIdLen())
	dst.SessionID = src.GetSessionId()
	dst.CipherSuitesLen = int(src.GetCipherSuitesLen())
	suites := src.GetCipherSuites()
	if len(suites) > 0 {
		dst.CipherSuites = make([]tlsproto.CipherSuite, 0, len(suites))
		for _, s := range suites {
			dst.CipherSuites = append(dst.CipherSuites, tlsproto.CipherSuite(s))
		}
	}
	comp := src.GetCompressMethods()
	if len(comp) > 0 {
		dst.CompressMethods = make([]tlsproto.CompressionMethod, 0, len(comp))
		for _, c := range comp {
			dst.CompressMethods = append(dst.CompressMethods, tlsproto.CompressionMethod(c))
		}
	}
	dst.ExtensionLen = int(src.GetExtensionLen())
	exts := src.GetExtensions()
	if len(exts) > 0 {
		dst.Extensions = make([]tlsutil.ExtensionItem, 0, len(exts))
		for _, e := range exts {
			dst.Extensions = append(dst.Extensions,
				tlsutil.ExtensionItem{
					Len:  uint16(e.GetLen()),
					Type: tlsproto.ExtensionType(e.GetEtype()),
				})
		}
	}
	info := src.GetExtensionInfo()
	if info != nil {
		dst.ExtensionInfo = DecodedInfo(info)
	}
	dst.UseGREASE = src.GetUseGrease()
	dst.JA3 = src.GetJa3()
	dst.JA3digest = src.GetJa3Digest()

	return dst
}

// ClientHelloDataPB convert to protobuf ClientHelloData
func ClientHelloDataPB(src *tlsutil.ClientHelloData) *pb.ConnectionData_ClientHelloData {
	dst := &pb.ConnectionData_ClientHelloData{}
	dst.ClientVersion = uint32(src.ClientVersion)
	dst.RandomLen = uint32(src.RandomLen)
	dst.SessionIdLen = uint32(src.SessionIDLen)
	dst.SessionId = src.SessionID
	dst.CipherSuitesLen = uint32(src.CipherSuitesLen)
	if len(src.CipherSuites) > 0 {
		dst.CipherSuites = make([]uint32, 0, len(src.CipherSuites))
		for _, s := range src.CipherSuites {
			dst.CipherSuites = append(dst.CipherSuites, uint32(s))
		}
	}
	if len(src.CompressMethods) > 0 {
		dst.CompressMethods = make([]uint32, 0, len(src.CompressMethods))
		for _, s := range src.CompressMethods {
			dst.CompressMethods = append(dst.CompressMethods, uint32(s))
		}
	}
	dst.ExtensionLen = int32(src.ExtensionLen)
	if len(src.Extensions) > 0 {
		dst.Extensions = make([]*pb.ConnectionData_ExtensionItem, 0, len(src.Extensions))
		for _, e := range src.Extensions {
			dst.Extensions = append(dst.Extensions,
				&pb.ConnectionData_ExtensionItem{
					Len:   uint32(e.Len),
					Etype: uint32(e.Type),
				})
		}
	}
	if src.ExtensionInfo != nil {
		dst.ExtensionInfo = DecodedInfoPB(src.ExtensionInfo)
	}
	dst.UseGrease = src.UseGREASE
	dst.Ja3 = src.JA3
	dst.Ja3Digest = src.JA3digest

	return dst
}

// ServerHelloData convert to model ServerHelloData
func ServerHelloData(src *pb.ConnectionData_ServerHelloData) *tlsutil.ServerHelloData {
	dst := &tlsutil.ServerHelloData{}
	dst.ServerVersion = tlslayer.ProtocolVersion(src.GetServerVersion())
	dst.RandomLen = int(src.GetRandomLen())
	dst.SessionIDLen = int(src.GetSessionIdLen())
	dst.SessionID = src.GetSessionId()
	dst.CipherSuiteSel = tlsproto.CipherSuite(src.GetCipherSuiteSel())
	dst.CompressMethodSel = tlsproto.CompressionMethod(src.GetCompressMethodSel())
	dst.ExtensionLen = int(src.GetExtensionLen())
	exts := src.GetExtensions()
	if len(exts) > 0 {
		dst.Extensions = make([]tlsutil.ExtensionItem, 0, len(exts))
		for _, e := range exts {
			dst.Extensions = append(dst.Extensions,
				tlsutil.ExtensionItem{
					Len:  uint16(e.GetLen()),
					Type: tlsproto.ExtensionType(e.GetEtype()),
				})
		}
	}
	info := src.GetExtensionInfo()
	if info != nil {
		dst.ExtensionInfo = DecodedInfo(info)
	}

	return dst
}

// ServerHelloDataPB convert to protobuf ServerHelloData
func ServerHelloDataPB(src *tlsutil.ServerHelloData) *pb.ConnectionData_ServerHelloData {
	dst := &pb.ConnectionData_ServerHelloData{}
	dst.ServerVersion = uint32(src.ServerVersion)
	dst.RandomLen = uint32(src.RandomLen)
	dst.SessionIdLen = uint32(src.SessionIDLen)
	dst.SessionId = src.SessionID
	dst.CipherSuiteSel = uint32(src.CipherSuiteSel)
	dst.CompressMethodSel = uint32(src.CompressMethodSel)
	dst.ExtensionLen = int32(src.ExtensionLen)
	if len(src.Extensions) > 0 {
		dst.Extensions = make([]*pb.ConnectionData_ExtensionItem, 0, len(src.Extensions))
		for _, e := range src.Extensions {
			dst.Extensions = append(dst.Extensions,
				&pb.ConnectionData_ExtensionItem{
					Len:   uint32(e.Len),
					Etype: uint32(e.Type),
				})
		}
	}
	if src.ExtensionInfo != nil {
		dst.ExtensionInfo = DecodedInfoPB(src.ExtensionInfo)
	}
	return dst
}

// DecodedInfo convert to model DecodedInfo
func DecodedInfo(src *pb.ConnectionData_DecodedInfo) *tlsutil.DecodedInfo {
	dst := &tlsutil.DecodedInfo{}
	dst.SNI = src.GetSni()
	schemes := src.GetSignatureSchemes()
	if len(schemes) > 0 {
		dst.SignatureSchemes = make([]tlsproto.SignatureScheme, 0, len(schemes))
		for _, s := range src.SignatureSchemes {
			dst.SignatureSchemes = append(dst.SignatureSchemes, tlsproto.SignatureScheme(s))
		}
	}
	versions := src.GetSupportedVersions()
	if len(versions) > 0 {
		dst.SupportedVersions = make([]tlsproto.SupportedVersion, 0, len(versions))
		for _, s := range versions {
			dst.SupportedVersions = append(dst.SupportedVersions, tlsproto.SupportedVersion(s))
		}
	}
	groups := src.GetSupportedGroups()
	if len(groups) > 0 {
		dst.SupportedGroups = make([]tlsproto.SupportedGroup, 0, len(groups))
		for _, s := range groups {
			dst.SupportedGroups = append(dst.SupportedGroups, tlsproto.SupportedGroup(s))
		}
	}
	ecpoints := src.GetEcPointFormats()
	if len(ecpoints) > 0 {
		dst.ECPointFormats = make([]tlsproto.ECPointFormat, 0, len(ecpoints))
		for _, s := range ecpoints {
			dst.ECPointFormats = append(dst.ECPointFormats, tlsproto.ECPointFormat(s))
		}
	}
	dst.OSCP = src.GetOscp()
	alpns := src.GetAlpns()
	if len(alpns) > 0 {
		dst.ALPNs = make([]string, 0, len(alpns))
		for _, s := range alpns {
			dst.ALPNs = append(dst.ALPNs, s)
		}
	}
	//TODO: KeyShareEntries
	//TODO: PSKKeyExchangeModes
	return dst
}

// DecodedInfoPB convert to protobuf DecodedInfo
func DecodedInfoPB(src *tlsutil.DecodedInfo) *pb.ConnectionData_DecodedInfo {
	dst := &pb.ConnectionData_DecodedInfo{}
	dst.Sni = src.SNI
	if len(src.SignatureSchemes) > 0 {
		dst.SignatureSchemes = make([]uint32, 0, len(src.SignatureSchemes))
		for _, s := range src.SignatureSchemes {
			dst.SignatureSchemes = append(dst.SignatureSchemes, uint32(s))
		}
	}
	if len(src.SupportedVersions) > 0 {
		dst.SupportedVersions = make([]uint32, 0, len(src.SupportedVersions))
		for _, s := range src.SupportedVersions {
			dst.SupportedVersions = append(dst.SupportedVersions, uint32(s))
		}
	}
	if len(src.SupportedGroups) > 0 {
		dst.SupportedGroups = make([]uint32, 0, len(src.SupportedGroups))
		for _, s := range src.SupportedGroups {
			dst.SupportedGroups = append(dst.SupportedGroups, uint32(s))
		}
	}
	if len(src.ECPointFormats) > 0 {
		dst.EcPointFormats = make([]uint32, 0, len(src.ECPointFormats))
		for _, s := range src.ECPointFormats {
			dst.EcPointFormats = append(dst.EcPointFormats, uint32(s))
		}
	}
	dst.Oscp = src.OSCP
	if len(src.ALPNs) > 0 {
		dst.Alpns = make([]string, 0, len(src.ALPNs))
		for _, s := range src.ALPNs {
			dst.Alpns = append(dst.Alpns, s)
		}
	}
	//TODO: KeyShareEntries
	//TODO: PSKKeyExchangeModes
	return dst
}

// CertSummary convert to model CertSummary
func CertSummary(src *pb.ConnectionData_CertSummary) tlsutil.CertSummary {
	dst := tlsutil.CertSummary{}
	dst.Digest = src.GetDigest()
	dst.Issuer = src.GetIssuer()
	dst.Subject = src.GetSubject()
	dst.IsCA = src.GetIsCa()
	return dst
}

// CertSummaryPB convert to protobuf CertSummary
func CertSummaryPB(src tlsutil.CertSummary) *pb.ConnectionData_CertSummary {
	dst := &pb.ConnectionData_CertSummary{}
	dst.Digest = src.Digest
	dst.Issuer = src.Issuer
	dst.Subject = src.Subject
	dst.IsCa = src.IsCA
	return dst
}
