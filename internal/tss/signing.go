package tss

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
	tsslib "github.com/bnb-chain/tss-lib/v2/tss"
	log "github.com/sirupsen/logrus"
)

func (ts *TSSServiceImpl) HandleSigningMessages(ctx context.Context, inCh chan SigningMessage, outCh chan tsslib.Message, endCh chan *common.SignatureData) {
	parties := 3
	threshold := 2
	partyIDs := createPartyIDs(parties)
	peerCtx := tsslib.NewPeerContext(partyIDs)
	params := tsslib.NewParameters(tsslib.S256(), peerCtx, partyIDs[0], parties, threshold)

	msgToSign := big.NewInt(123456)
	savedData, err := loadTSSData()
	if err != nil {
		log.Errorf("Failed to load TSS data: %v", err)
		return
	}

	party := signing.NewLocalParty(msgToSign, params, *savedData, outCh, endCh)
	if err := party.Start(); err != nil {
		log.Errorf("TSS signing process failed to start: %v", err)
		return
	}

	for {
		select {
		case msg := <-inCh:
			var tssMsg tsslib.Message
			if err := json.Unmarshal([]byte(msg.Content), &tssMsg); err != nil {
				log.Errorf("Failed to unmarshal TSS signing message: %v", err)
				continue
			}
			// TODO: handle the message
		case <-ctx.Done():
			return
		}
	}
}
