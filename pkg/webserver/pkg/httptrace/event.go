// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptrace

// Event is a bitmask that selects which httptrace events are logged.
// Events can be OR-combined, e.g. EventDNS | EventTLS.
//
//go:generate stringer -type Event -trimprefix=Event
type Event uint32

const (
	// EventGetConn logs GetConn.
	EventGetConn Event = 1 << iota
	// EventGotConn logs GotConn.
	EventGotConn
	// EventDNSStart logs DNSStart.
	EventDNSStart
	// EventDNSDone logs DNSDone.
	EventDNSDone
	// EventConnectStart logs ConnectStart.
	EventConnectStart
	// EventConnectDone logs ConnectDone.
	EventConnectDone
	// EventTLSHandshakeStart logs TLSHandshakeStart.
	EventTLSHandshakeStart
	// EventTLSHandshakeDone logs TLSHandshakeDone.
	EventTLSHandshakeDone
	// EventWroteRequest logs WroteRequest.
	EventWroteRequest
	// EventGotFirstResponseByte logs GotFirstResponseByte.
	EventGotFirstResponseByte
	// EventPutIdleConn logs PutIdleConn.
	EventPutIdleConn
)

// Event groups for convenience.
const (
	// EventConn covers GetConn/GotConn/ConnectStart/ConnectDone/PutIdleConn.
	EventConn = EventGetConn | EventGotConn | EventConnectStart | EventConnectDone | EventPutIdleConn
	// EventDNS covers DNSStart/DNSDone.
	EventDNS = EventDNSStart | EventDNSDone
	// EventTLS covers TLSHandshakeStart/TLSHandshakeDone.
	EventTLS = EventTLSHandshakeStart | EventTLSHandshakeDone
	// EventRoundTrip covers WroteRequest and GotFirstResponseByte (TTFB).
	EventRoundTrip = EventWroteRequest | EventGotFirstResponseByte
	// EventDefault is logged when WithEvents is unset.
	// It covers connection acquisition and the DNS result.
	EventDefault = EventGetConn | EventGotConn | EventDNSDone
	// EventAll enables every supported event.
	EventAll = EventConn | EventDNS | EventTLS | EventRoundTrip
)

// has reports whether e includes all bits in mask.
func (e Event) has(mask Event) bool { return e&mask != 0 }
