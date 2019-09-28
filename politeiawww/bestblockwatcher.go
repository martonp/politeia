// Copyright (c) 2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import exptypes "github.com/decred/dcrdata/explorer/types"

func (p *politeiawww) setupBestBlockWatcher() {
	p.wsDcrdata.subToPing()
	p.wsDcrdata.subToNewBlock()

	go func() {
		for {

			msg, ok := <-p.wsDcrdata.client.Receive()
			if !ok {
				log.Infof("BAD MESSAGE")
				break
			}
			if msg == nil {
				log.Errorf("ReceiveMsg failed")
				continue
			}

			switch m := msg.Message.(type) {
			case *exptypes.WebsocketBlock:
				log.Infof("Message (%s): WebsocketBlock(hash=%s)", msg.EventId, m.Block.Height)
				p.bestBlock = uint64(m.Block.Height)
			default:
				log.Debugf("Message of type %v unhandled.", m)
				continue
			}
		}
	}()
}
