package main

import (
	"context"
	"crypto/tls"
	"log"
	"strings"
	"time"

	xmpp "github.com/xmppo/go-xmpp"
)

// XMPPClient manages the connection to the XMPP server and routes messages.
type XMPPClient struct {
	cfg        AppConfig
	dispatcher *Dispatcher
}

// NewXMPPClient creates an XMPPClient with the given config and dispatcher.
func NewXMPPClient(cfg AppConfig, d *Dispatcher) *XMPPClient {
	return &XMPPClient{cfg: cfg, dispatcher: d}
}

// Connect establishes (and re-establishes) the XMPP connection until ctx is cancelled.
func (c *XMPPClient) Connect(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		log.Printf("xmpp: connecting to %s as %s", c.cfg.Server, c.cfg.JID)
		err := c.runSession(ctx)
		if ctx.Err() != nil {
			return ctx.Err()
		}

		log.Printf("xmpp: session ended (%v), reconnecting in 10s", err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(10 * time.Second):
		}
	}
}

// runSession opens a single XMPP session and blocks until it ends.
func (c *XMPPClient) runSession(ctx context.Context) error {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: c.cfg.TLSSkipVerify, //nolint:gosec // opt-in for self-signed homelabs
	}

	opts := xmpp.Options{
		Host:                         ensurePort(c.cfg.Server),
		User:                         c.cfg.JID,
		Password:                     c.cfg.Password,
		NoTLS:                        true,
		StartTLS:                     c.cfg.StartTLS,
		TLSConfig:                    tlsCfg,
		InsecureAllowUnencryptedAuth: c.cfg.InsecureAllowUnencryptedAuth,
		Debug:                        false,
		Session:                      false,
		Status:                       "xa",
		StatusMessage:                "Vox Imperialis online",
	}

	client, err := opts.NewClient()
	if err != nil {
		return err
	}
	defer client.Close()

	log.Printf("xmpp: connected as %s", c.cfg.JID)

	// Cancel the blocking Recv when context is done.
	connDone := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			client.Close()
		case <-connDone:
		}
	}()
	defer close(connDone)

	for {
		stanza, err := client.Recv()
		if err != nil {
			return err
		}

		switch v := stanza.(type) {
		case xmpp.Chat:
			if v.Text != "" && v.Type != "groupchat" && v.Type != "error" {
				c.handleMessage(client, v)
			}
		}
	}
}

// handleMessage authenticates the sender, parses the command, and sends a reply.
func (c *XMPPClient) handleMessage(client *xmpp.Client, msg xmpp.Chat) {
	if !IsAllowed(msg.Remote, c.cfg.AllowedUsers) {
		log.Printf("xmpp: ignoring message from unauthorized JID %s", msg.Remote)
		return
	}

	cmd := ParseCommand(msg.Text, msg.Remote)
	response := c.dispatcher.Dispatch(cmd)

	reply := xmpp.Chat{Remote: msg.Remote, Type: "chat", Text: response}
	if _, err := client.Send(reply); err != nil {
		log.Printf("xmpp: failed to send reply to %s: %v", msg.Remote, err)
	}
}

// ensurePort appends :5222 to host if no port is specified.
func ensurePort(host string) string {
	if strings.Contains(host, ":") {
		return host
	}
	return host + ":5222"
}
