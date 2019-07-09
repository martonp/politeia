// Copyright (c) 2017-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"net/http"

	pd "github.com/decred/politeia/politeiad/api/v1"
	"github.com/decred/politeia/util"
)

const cmsReadme = `
	### CMS Readme
`

const politeiaReadme = `
	### Politeia Readme
`

func (p *politeiawww) updateReadme(readmeContents string) error {

	challenge, err := util.Random(pd.ChallengeSize)
	if err != nil {
		return err
	}

	urm := pd.UpdateReadme{
		Challenge: hex.EncodeToString(challenge),
		Content:   string(readmeContents),
	}

	// Send politeiad request
	response, err := p.makeRequest(http.MethodPost,
		pd.UpdateReadmeRoute, urm)

	if err != nil {
		return err
	}

	return nil
}

func (p *politeiawww) UpdateWWWReadme() error {
	return p.updateReadme(politeiaReadme)
}

func (p *politeiawww) UpdateCMSReadme() error {
	return p.updateReadme(cmsReadme)
}
