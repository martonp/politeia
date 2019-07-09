# Politeia Proposal Data

This repo contains the data of the proposals submitted to politeia. The "Anchors" folder and anchor_audit_trail.txt contain data related to the anchoring of data to the decred blockchain through dcrtime, which facilitates Politeia’s transparent censorship functionality. The folders named with a 64 byte hash contain the data relevant to each proposal. 

Within each proposal folder there are sequentially numbered sub-folders, one for each version of the proposal. The first folder (‘/1’) corresponds to the initial version of the proposal submitted to Politeia. If a proposal author updates their proposal (typically in response to community feedback), the updated version and its associated data are put into folder (‘/2’), and so on.

Below is a table with descriptions of the files and folders found in each proposal version folder. 

### Proposal Data


| Folder/File   | Description |
| ------------- |:-------------:| -----:|
| /payload      | Folder containing an index.md file, which has the full text of the proposal and any images associated with the proposal. |
| /plugins/decred  |  Folder containing the comments.journal file, which contains comments and up/down vote data for comments. If a proposal has begun voting (or finished voting), this folder will also contain a ballot.journal file containing the ticket holder voting data. |
| 00.metadata.txt  |    Data about the proposal submission:

    timestamp
    pubkey of submitter
    title
    signature |

| 02.metadata.txt  | Data about admin approval for display:

    timestamp
    pubkey of admin |   

| 13.metadata.txt |  Data about proposal owner authorizing the start of voting. Where proposal owner has authorized and then rescinded authorization this will appear in multiple commits:

    timestamp
    owner pubkey
    transaction hash|

| 14.metadata.txt |	Data about admin authorizing the start of voting. Includes specification of the vote: duration (2016 blocks):

    quorum requirement (20%) (subject to change)
    approval threshold (60%) (subject to change)
    id and descriptions for voting options |

| 15.metadata.txt |	Data about voting period:

    starting block height and hash
    ending block height
    a list of every ticket eligible to vote on the proposal |

| recordmetadata.json  | 	Metadata about the record. |