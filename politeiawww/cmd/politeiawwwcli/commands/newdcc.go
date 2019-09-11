// Copyright (c) 2017-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package commands

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/decred/politeia/politeiad/api/v1/mime"
	cms "github.com/decred/politeia/politeiawww/api/cms/v1"
	www "github.com/decred/politeia/politeiawww/api/www/v1"
	"github.com/decred/politeia/util"
)

// domainTypes gives human readable output for the various domain types available
var domainTypes = map[cms.DomainTypeT]string{
	cms.DomainTypeDeveloper:     "(1) Developer",
	cms.DomainTypeMarketing:     "(2) Marketing",
	cms.DomainTypeCommunity:     "(3) Community Management",
	cms.DomainTypeResearch:      "(4) Research",
	cms.DomainTypeDesign:        "(5) Design",
	cms.DomainTypeDocumentation: "(6) Documentation",
}

// contractorTypes gives human readable output for the various contractor types available
var contractorTypes = map[cms.ContractorTypeT]string{
	cms.ContractorTypeDirect:        "(1) Direct",
	cms.ContractorTypeSubContractor: "(3) Sub-contractor",
}

// NewDCCCmd submits a new dcc.
type NewDCCCmd struct {
	Args struct {
		Type        uint     `positional-arg-name:"type"`            // 1 for Issuance, 2 for Revocation
		Attachments []string `positional-arg-name:"attachmentfiles"` // DCC attachment files
	} `positional-args:"true" optional:"true"`
	Type           string ``
	NomineeUserID  string `long:"nomineeuserid" optional:"true" description:"The UserID of the Nominated User"`
	Statement      string `long:"statement" optional:"true" description:"Statement in support of the DCC"`
	Domain         string `long:"domain" optional:"true" description:"The domain of the nominated user"`
	ContractorType string `long:"contractortype" optional:"true" description:"The contractor type of the nominated user"`
}

// Execute executes the new dcc command.
func (cmd *NewDCCCmd) Execute(args []string) error {

	// Check for a valid DCC type
	if int(cmd.Args.Type) <= 0 || int(cmd.Args.Type) > 2 {
		return errInvalidDCCType
	}

	// Check for user identity
	if cfg.Identity == nil {
		return errUserIdentityNotFound
	}

	// Get server public key
	vr, err := client.Version()
	if err != nil {
		return err
	}
	var domainType int
	var contractorType int
	if cmd.Statement == "" || cmd.NomineeUserID == "" || cmd.Domain == "" ||
		cmd.ContractorType == "" {
		reader := bufio.NewReader(os.Stdin)
		if cmd.Statement == "" {
			fmt.Print("Enter your statement to support the DCC: ")
			cmd.Statement, _ = reader.ReadString('\n')
		}
		if cmd.NomineeUserID == "" {
			fmt.Print("Enter the nominee user id: ")
			cmd.NomineeUserID, _ = reader.ReadString('\n')
		}
		if cmd.Domain == "" {
			for {
				fmt.Printf("Domain Type: " +
					domainTypes[cms.DomainTypeDeveloper] + ", " +
					domainTypes[cms.DomainTypeMarketing] + ", " +
					domainTypes[cms.DomainTypeCommunity] + ", " +
					domainTypes[cms.DomainTypeResearch] + ", " +
					domainTypes[cms.DomainTypeDesign] + ", " +
					domainTypes[cms.DomainTypeDocumentation] + ": ")
				cmd.Domain, _ = reader.ReadString('\n')
				domainType, err = strconv.Atoi(strings.TrimSpace(cmd.Domain))
				if err != nil {
					fmt.Println("Invalid entry, please try again.")
					continue
				}
				if domainType < 1 || domainType > 6 {
					fmt.Println("Invalid domain type entered, please try again.")
					continue
				}
				str := fmt.Sprintf(
					"Your current Domain setting is: \"%v\" Keep this?",
					domainType)
				update, err := promptListBool(reader, str, "yes")
				if err != nil {
					return err
				}
				if update {
					break
				}
			}
		} else {
			domainType, err = strconv.Atoi(strings.TrimSpace(cmd.Domain))
			if err != nil {
				return fmt.Errorf("invalid domain type: %v", err)
			}
			if domainType < 1 || domainType > 6 {
				return fmt.Errorf("invalid domain type")
			}
		}
		if cmd.ContractorType == "" {
			for {
				fmt.Printf("Contractor Type: " +
					contractorTypes[cms.ContractorTypeDirect] + ", " +
					contractorTypes[cms.ContractorTypeSubContractor] + ": ")
				cmd.ContractorType, _ = reader.ReadString('\n')
				contractorType, err = strconv.Atoi(strings.TrimSpace(cmd.ContractorType))
				if err != nil {
					fmt.Println("Invalid entry, please try again.")
					continue
				}
				if contractorType != 1 && contractorType != 3 {
					fmt.Println("Invalid contractor type entered, please try again.")
					continue
				}
				str := fmt.Sprintf(
					"Your current Contractor Type setting is: \"%v\" Keep this?",
					domainType)
				update, err := promptListBool(reader, str, "yes")
				if err != nil {
					return err
				}
				if update {
					break
				}
			}
		} else {
			contractorType, err = strconv.Atoi(strings.TrimSpace(cmd.ContractorType))
			if err != nil {
				return fmt.Errorf("invalid contractor type: %v", err)
			}
			if contractorType != 1 && contractorType != 3 {
				return fmt.Errorf("invalid contractor type")
			}
		}
		fmt.Print("\nPlease carefully review your information and ensure it's " +
			"correct. If not, press Ctrl + C to exit. Or, press Enter to continue.")
		reader.ReadString('\n')
	}

	dccInput := &cms.DCCInput{}
	dccInput.SponsorStatement = strings.TrimSpace(cmd.Statement)
	dccInput.NomineeUserID = strings.TrimSpace(cmd.NomineeUserID)
	dccInput.Type = cms.DCCTypeT(int(cmd.Args.Type))
	dccInput.Domain = cms.DomainTypeT(domainType)
	dccInput.ContractorType = cms.ContractorTypeT(contractorType)

	// Print request details
	err = printJSON(dccInput)
	if err != nil {
		return err
	}
	b, err := json.Marshal(dccInput)
	if err != nil {
		return fmt.Errorf("Marshal: %v", err)
	}

	f := www.File{
		Name:    "dcc.json",
		MIME:    mime.DetectMimeType(b),
		Digest:  hex.EncodeToString(util.Digest(b)),
		Payload: base64.StdEncoding.EncodeToString(b),
	}

	files := make([]www.File, 0, 1)
	files = append(files, f)

	// Compute merkle root and sign it
	sig, err := signedMerkleRoot(files, cfg.Identity)
	if err != nil {
		return fmt.Errorf("SignMerkleRoot: %v", err)
	}

	// Setup new dcc request
	nd := cms.NewDCC{
		File:      f,
		PublicKey: hex.EncodeToString(cfg.Identity.Public.Key[:]),
		Signature: sig,
	}

	// Print request details
	err = printJSON(nd)
	if err != nil {
		return err
	}

	// Send request
	ndr, err := client.NewDCC(nd)
	if err != nil {
		return err
	}

	ndFiles := make([]www.File, 0, 1)
	ndFiles = append(ndFiles, nd.File)

	// Verify the censorship record
	pr := www.ProposalRecord{
		Files:            ndFiles,
		PublicKey:        nd.PublicKey,
		Signature:        nd.Signature,
		CensorshipRecord: ndr.CensorshipRecord,
	}
	err = verifyProposal(pr, vr.PubKey)
	if err != nil {
		return fmt.Errorf("unable to verify proposal %v: %v",
			pr.CensorshipRecord.Token, err)
	}

	// Print response details
	return printJSON(ndr)
}
