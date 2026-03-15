package main

import "strings"

// IsAllowed reports whether jid is in the configured allowlist.
// The resource part of the JID is stripped before comparison.
func IsAllowed(jid string, allowedUsers []string) bool {
	bare := bareJID(jid)
	for _, u := range allowedUsers {
		if strings.EqualFold(bare, strings.TrimSpace(u)) {
			return true
		}
	}
	return false
}

// bareJID strips the resource from a full JID.
// "user@domain/resource" → "user@domain"
func bareJID(jid string) string {
	if i := strings.Index(jid, "/"); i >= 0 {
		return jid[:i]
	}
	return jid
}
