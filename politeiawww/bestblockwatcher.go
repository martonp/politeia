// Copyright (c) 2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import exptypes "github.com/decred/dcrdata/explorer/types/v2"

func (p *politeiawww) setupBestBlockWatcher() {
	p.wsDcrdata.subToPing()
	p.wsDcrdata.subToNewBlock()
	log.Infof("FUCK YOUU!~")

	go func() {
		for {
			log.Infof("FUCK YOUU12222!~")

			msg, ok := <-p.wsDcrdata.client.Receive()
			if !ok {
				break
			}
			if msg == nil {
				log.Errorf("ReceiveMsg failed")
				continue
			}

			switch m := msg.Message.(type) {
			case string:
				log.Infof("Message (%s): %s", msg.EventId, m)
			case int:
				log.Infof("Message (%s): %v", msg.EventId, m)
			case *exptypes.WebsocketBlock:
				log.Debugf("Message (%s): WebsocketBlock(height=%v)", m)
			default:
				log.Infof("Message of type %v unhandled. %v", msg.EventId, msg.Message)
			}
		}
	}()
}
