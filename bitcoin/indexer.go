package bitcoin

import (
	"fmt"

	"errors"

	"github.com/btcsuite/btcd/txscript"
	"github.com/evmos/ethermint/types"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

var (
	ErrParsePkScript       = errors.New("parse pkscript err")
	ErrDecodeListenAddress = errors.New("decode listen address err")
)

// BitcoinIndexer bitcoin indexer, parse and forward data
type BitcoinIndexer struct {
	client        *rpcclient.Client // call bitcoin rpc client
	chainParams   *chaincfg.Params  // bitcoin network params, e.g. mainnet, testnet, etc.
	listenAddress btcutil.Address   // need listened bitcoin address
}

// NewBitcoinIndexer new bitcoin indexer
func NewBitcoinIndexer(client *rpcclient.Client, chainParams *chaincfg.Params, listenAddress string) (*BitcoinIndexer, error) {
	// check listenAddress
	address, err := btcutil.DecodeAddress(listenAddress, chainParams)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrDecodeListenAddress, err)
	}
	return &BitcoinIndexer{
		client:        client,
		chainParams:   chainParams,
		listenAddress: address,
	}, nil
}

// ParseBlock parse block data by block height
// NOTE: Currently, only transfer transactions are supported.
func (b *BitcoinIndexer) ParseBlock(height int64) ([]types.BitcoinTxParseResult, error) {
	blockResult, err := b.getBlockByHeight(height)
	if err != nil {
		return nil, err
	}

	blockParsedResult := make([]types.BitcoinTxParseResult, 0)
	for _, v := range blockResult.Transactions {
		sendDatas, err := b.parseTx(v.TxHash())
		if err != nil {
			return nil, err
		}
		blockParsedResult = append(blockParsedResult, sendDatas...)
	}

	return blockParsedResult, nil
}

// getBlockByHeight returns a raw block from the server given its height
func (b *BitcoinIndexer) getBlockByHeight(height int64) (*wire.MsgBlock, error) {
	blockhash, err := b.client.GetBlockHash(height)
	if err != nil {
		return nil, err
	}
	return b.client.GetBlock(blockhash)
}

// parseTx parse transaction data
func (b *BitcoinIndexer) parseTx(txHash chainhash.Hash) (parsedResult []types.BitcoinTxParseResult, err error) {
	txResult, err := b.client.GetRawTransaction(&txHash)
	if err != nil {
		return nil, fmt.Errorf("get raw transaction err:%w", err)
	}

	for _, v := range txResult.MsgTx().TxOut {
		pkAddress, err := b.parseAddress(v.PkScript)
		if err != nil {
			if errors.Is(err, ErrParsePkScript) {
				continue
			}
			return nil, err
		}

		// if pk address eq dest listened address, after parse from address by vin prev tx
		if pkAddress == b.listenAddress.EncodeAddress() {
			fromAddress, err := b.parseFromAddress(txResult)
			if err != nil {
				return nil, fmt.Errorf("vin parse err:%w", err)
			}
			parsedResult = append(parsedResult, types.BitcoinTxParseResult{Value: v.Value, From: fromAddress, To: pkAddress})
		}
	}

	return
}

// parseFromAddress from vin parse to address
// return all possible values parsed to
func (b *BitcoinIndexer) parseFromAddress(txResult *btcutil.Tx) (toAddress []string, err error) {
	for _, vin := range txResult.MsgTx().TxIn {
		// get prev tx hash
		prevTxId := vin.PreviousOutPoint.Hash
		vinResult, err := b.client.GetRawTransaction(&prevTxId)
		if err != nil {
			return nil, fmt.Errorf("vin get raw transaction err:%w", err)
		}
		if len(vinResult.MsgTx().TxOut) <= 0 {
			return nil, fmt.Errorf("vin txOut is null")
		}
		vinPKScript := vinResult.MsgTx().TxOut[vin.PreviousOutPoint.Index].PkScript
		//  script to address
		vinPkAddress, err := b.parseAddress(vinPKScript)
		if err != nil {
			if errors.Is(err, ErrParsePkScript) {
				continue
			}
			return nil, err
		}

		toAddress = append(toAddress, vinPkAddress)
	}
	return
}

// parseAddress from pkscript parse address
func (b *BitcoinIndexer) parseAddress(pkScript []byte) (string, error) {
	pk, err := txscript.ParsePkScript(pkScript)
	if err != nil {
		return "", fmt.Errorf("%w:%w", ErrParsePkScript, err)
	}

	//  encodes the script into an address for the given chain.
	pkAddress, err := pk.Address(b.chainParams)
	if err != nil {
		return "", fmt.Errorf("PKScript to address err:%w", err)
	}
	return pkAddress.EncodeAddress(), nil
}