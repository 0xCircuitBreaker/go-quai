package types

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/dominant-strategies/go-quai/common"
	"github.com/stretchr/testify/require"
)

func TestManifestEncodeDecode(t *testing.T) {
	// Create a new manifest
	hash1 := common.BytesToHash([]byte{0x01})
	manifest := BlockManifest{hash1}
	manifest = append(manifest, hash1)

	// Encode the manifest to ProtoManifest format
	protoManifest, err := manifest.ProtoEncode()
	if err != nil {
		t.Errorf("Failed to encode manifest: %v", err)
	}

	// Decode the ProtoManifest into a new Manifest
	decodedManifest := BlockManifest{}
	err = decodedManifest.ProtoDecode(protoManifest)
	if err != nil {
		t.Errorf("Failed to decode manifest: %v", err)
	}

	require.Equal(t, manifest, decodedManifest)
}

func headerTestData() (*Header, common.Hash) {
	header := &Header{
		parentHash:               []common.Hash{common.HexToHash("0x123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0"), common.HexToHash("0x123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0")},
		uncleHash:                common.HexToHash("0x23456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef1"),
		evmRoot:                  common.HexToHash("0x456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef3"),
		quaiStateSize:            big.NewInt(1000),
		utxoRoot:                 common.HexToHash("0x56789abcdef0123456789abcdef0123456789abcdef0123456789abcdef4"),
		txHash:                   common.HexToHash("0x6789abcdef0123456789abcdef0123456789abcdef0123456789abcdef5"),
		outboundEtxHash:          common.HexToHash("0x789abcdef0123456789abcdef0123456789abcdef0123456789abcdef6"),
		etxRollupHash:            common.HexToHash("0x9abcdef0123456789abcdef0123456789abcdef0123456789abcdef8"),
		manifestHash:             []common.Hash{common.HexToHash("0xabcdef0123456789abcdef0123456789abcdef0123456789abcdef9"), common.HexToHash("0xabcdef0123456789abcdef0123456789abcdef0123456789abcdef9"), common.HexToHash("0xabcdef0123456789abcdef0123456789abcdef0123456789abcdef9")},
		receiptHash:              common.HexToHash("0xbcdef0123456789abcdef0123456789abcdef0123456789abcdefa"),
		parentEntropy:            []*big.Int{big.NewInt(123456789), big.NewInt(123456789), big.NewInt(123456789)},
		parentDeltaEntropy:       []*big.Int{big.NewInt(123456789), big.NewInt(123456789), big.NewInt(123456789)},
		parentUncledDeltaEntropy: []*big.Int{big.NewInt(123456789), big.NewInt(123456789), big.NewInt(123456789)},
		efficiencyScore:          12345,
		thresholdCount:           12345,
		expansionNumber:          123,
		etxEligibleSlices:        common.HexToHash("0xcdef0123456789abcdef0123456789abcdef0123456789abcdefb"),
		primeTerminusHash:        common.HexToHash("0xdef0123456789abcdef0123456789abcdef0123456789abcdefc"),
		interlinkRootHash:        common.HexToHash("0xef0123456789abcdef0123456789abcdef0123456789abcdefd"),
		uncledEntropy:            big.NewInt(123456789),
		number:                   []*big.Int{big.NewInt(123456789), big.NewInt(123456789)},
		gasLimit:                 123456789,
		gasUsed:                  987654321,
		baseFee:                  big.NewInt(123456789),
		extra:                    []byte("SGVsbG8gd29ybGQ="),
		stateLimit:               1234567,
		stateUsed:                1234567,
		exchangeRate:             big.NewInt(123456789),
		avgTxFees:                big.NewInt(1000000),
		totalFees:                big.NewInt(1432323123),
		kQuaiDiscount:            big.NewInt(123456),
		conversionFlowAmount:     big.NewInt(123456),
		minerDifficulty:          big.NewInt(1234899),
		primeStateRoot:           common.HexToHash("0xcdef0123456789abcdef0123456789abcdef0123456789abcdefb"),
		regionStateRoot:          common.HexToHash("0xcdef0123456789abcdef0123456789abcdef0123456789abcdefb"),
	}

	return header, header.Hash()
}

func TestHeaderHash(t *testing.T) {
	_, hash := headerTestData()
	correctHash := common.HexToHash("0x900c189590b1d12744cf9bdf760d0ed92f8862b87b0479e701a656f926dfa34f")
	require.Equal(t, correctHash, hash, "Hash not equal to expected hash")
}

var testInt64 = int64(987654321)
var testUInt8 = uint8(123)
var testUInt16 = uint16(54321)
var testUInt64 = uint64(123456789)
var testByte = []byte("test byte")

func fuzzHeaderHash(f *testing.F, getField func(*Header) common.Hash, setField func(*Header, common.Hash)) {
	header, _ := headerTestData()
	f.Add(testByte)
	f.Add(getField(header).Bytes())
	f.Fuzz(func(t *testing.T, b []byte) {
		localHeader, hash := headerTestData()
		bHash := common.BytesToHash(b)
		if getField(localHeader) != bHash {
			setField(localHeader, bHash)
			require.NotEqual(t, localHeader.Hash(), hash, "Hash equal for root \noriginal: %v, modified: %v", getField(header), bHash)
		}
	})
}

func fuzzHeaderUint64Hash(f *testing.F, getField func(*Header) uint64, setField func(*Header, uint64)) {
	header, _ := headerTestData()
	f.Add(testUInt64)
	f.Add(getField(header))
	f.Fuzz(func(t *testing.T, i uint64) {
		localHeader, hash := headerTestData()
		if getField(localHeader) != i {
			setField(localHeader, i)
			require.NotEqual(t, localHeader.Hash(), hash, "Hash equal for field \noriginal: %v, modified: %v", getField(header), i)
		}
	})
}

func fuzzHeaderHashLoopField(f *testing.F, getField func(*Header) []common.Hash, setField func(*Header, int, common.Hash)) {
	header, _ := headerTestData()
	f.Add(testByte)
	f.Add(getField(header)[0].Bytes())
	f.Fuzz(func(t *testing.T, b []byte) {
		localHeader, hash := headerTestData()
		bHash := common.BytesToHash(b)
		hashes := getField(localHeader)
		for i, h := range hashes {
			if bHash != h {
				setField(localHeader, i, bHash)
				require.NotEqual(t, localHeader.Hash(), hash, "Hash equal for hash field \noriginal: %v, modified: %v", h, bHash)
				//reset hash for next iteration
				setField(localHeader, i, h)
			}
		}
	})
}

func fuzzHeaderBigIntHash(f *testing.F, getField func(*Header) *big.Int, setField func(*Header, *big.Int)) {
	header, _ := headerTestData()
	f.Add(testInt64)
	f.Add(getField(header).Int64())
	f.Fuzz(func(t *testing.T, i int64) {
		localHeader, hash := headerTestData()
		bi := big.NewInt(i)
		if getField(localHeader).Cmp(bi) != 0 {
			setField(localHeader, bi)
			require.NotEqual(t, localHeader.Hash(), hash, "Hash equal for field \noriginal: %v, modified: %v", getField(header), bi)
		}
	})
}

func fuzzHeaderBigIntLoopHash(f *testing.F, getField func(*Header) []*big.Int, setField func(*Header, int, *big.Int)) {
	header, _ := headerTestData()
	f.Add(testInt64)
	f.Add(getField(header)[0].Int64())
	f.Fuzz(func(t *testing.T, i int64) {
		bi := big.NewInt(i)
		localHeader, hash := headerTestData()
		bigInts := getField(localHeader)
		for i, bigInt := range bigInts {
			if bigInt.Cmp(bi) != 0 {
				setField(localHeader, i, bi)
				require.NotEqual(t, localHeader.Hash(), hash, "Hash equal for bigInt field \noriginal: %v, modified: %v", bigInt, bi)
				//reset bigInt for next iteration
				setField(localHeader, i, bigInt)
			}
		}
	})
}

func fuzzHeaderUint16FieldHash(f *testing.F, getField func(*Header) uint16, setField func(*Header, uint16)) {
	header, _ := headerTestData()
	f.Add(testUInt16)
	f.Add(getField(header))
	f.Fuzz(func(t *testing.T, i uint16) {
		localHeader, hash := headerTestData()
		if getField(localHeader) != i {
			setField(localHeader, i)
			require.NotEqual(t, localHeader.Hash(), hash, "Hash equal for field \noriginal: %v, modified: %v", getField(header), i)
		}
	})
}

func FuzzHeaderParentHash(f *testing.F) {
	fuzzHeaderHashLoopField(f,
		func(h *Header) []common.Hash { return h.parentHash },
		func(h *Header, i int, hash common.Hash) { h.parentHash[i] = hash })
}

func FuzzHeaderUncleHash(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.uncleHash }, func(h *Header, hash common.Hash) { h.uncleHash = hash })
}

func FuzzHeaderEvmRootHash(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.evmRoot }, func(h *Header, hash common.Hash) { h.evmRoot = hash })
}

func FuzzHeaderUtxoRootHash(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.utxoRoot }, func(h *Header, hash common.Hash) { h.utxoRoot = hash })
}

func FuzzHeaderTxHash(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.txHash }, func(h *Header, hash common.Hash) { h.txHash = hash })
}

func FuzzHeaderOutboundEtxHash(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.outboundEtxHash }, func(h *Header, hash common.Hash) { h.outboundEtxHash = hash })
}

func FuzzHeaderEtxRollupHash(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.etxRollupHash }, func(h *Header, hash common.Hash) { h.etxRollupHash = hash })
}

func FuzzHeaderManifestHash(f *testing.F) {
	fuzzHeaderHashLoopField(f,
		func(h *Header) []common.Hash { return h.manifestHash },
		func(h *Header, i int, hash common.Hash) { h.manifestHash[i] = hash })
}

func FuzzHeaderReceiptHash(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.receiptHash }, func(h *Header, hash common.Hash) { h.receiptHash = hash })
}

func FuzzHeaderParentEntropyHash(f *testing.F) {
	fuzzHeaderBigIntLoopHash(f,
		func(h *Header) []*big.Int { return h.parentEntropy },
		func(h *Header, i int, bi *big.Int) { h.parentEntropy[i] = bi })
}

func FuzzHeaderParentDeltaEntropyHash(f *testing.F) {
	fuzzHeaderBigIntLoopHash(f,
		func(h *Header) []*big.Int { return h.parentDeltaEntropy },
		func(h *Header, i int, bi *big.Int) { h.parentDeltaEntropy[i] = bi })
}

func FuzzHeaderParentUncledDeltaEntropyHash(f *testing.F) {
	fuzzHeaderBigIntLoopHash(f,
		func(h *Header) []*big.Int { return h.parentUncledDeltaEntropy },
		func(h *Header, i int, bi *big.Int) { h.parentUncledDeltaEntropy[i] = bi })
}

func FuzzHeaderEfficiencyScoreHash(f *testing.F) {
	fuzzHeaderUint16FieldHash(f,
		func(h *Header) uint16 { return h.efficiencyScore },
		func(h *Header, i uint16) { h.efficiencyScore = i })
}

func FuzzHeaderThresholdCountHash(f *testing.F) {
	fuzzHeaderUint16FieldHash(f,
		func(h *Header) uint16 { return h.thresholdCount },
		func(h *Header, i uint16) { h.thresholdCount = i })
}

func FuzzHeaderExpansionNumberHash(f *testing.F) {
	header, _ := headerTestData()
	f.Add(testUInt8)
	f.Add(header.expansionNumber)
	f.Fuzz(func(t *testing.T, i uint8) {
		localHeader, hash := headerTestData()
		if localHeader.expansionNumber != i {
			localHeader.expansionNumber = i
			require.NotEqual(t, localHeader.Hash(), hash, "Hash equal for expansionNumber \noriginal: %v, modified: %v", header.expansionNumber, i)
		}
	})
}

func FuzzHeaderEtxEligibleSlicesHash(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.etxEligibleSlices }, func(h *Header, hash common.Hash) { h.etxEligibleSlices = hash })
}

func FuzzHeaderPrimeTerminusHash(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.primeTerminusHash }, func(h *Header, hash common.Hash) { h.primeTerminusHash = hash })
}

func FuzzHeaderInterlinkRootHashHash(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.interlinkRootHash }, func(h *Header, hash common.Hash) { h.interlinkRootHash = hash })
}

func FuzzHeaderUncledSHash(f *testing.F) {
	fuzzHeaderBigIntHash(f,
		func(h *Header) *big.Int { return h.uncledEntropy },
		func(h *Header, bi *big.Int) { h.uncledEntropy = bi })
}

func FuzzHeaderNumberHash(f *testing.F) {
	fuzzHeaderBigIntLoopHash(f,
		func(h *Header) []*big.Int { return h.number },
		func(h *Header, i int, bi *big.Int) { h.number[i] = bi })
}

func FuzzHeaderGasLimitHash(f *testing.F) {
	fuzzHeaderUint64Hash(f,
		func(h *Header) uint64 { return h.gasLimit },
		func(h *Header, i uint64) { h.gasLimit = i })
}

func FuzzHeaderGasUsedHash(f *testing.F) {
	fuzzHeaderUint64Hash(f,
		func(h *Header) uint64 { return h.gasUsed },
		func(h *Header, i uint64) { h.gasUsed = i })
}

func FuzzHeaderBaseFeeHash(f *testing.F) {
	fuzzHeaderBigIntHash(f,
		func(h *Header) *big.Int { return h.baseFee },
		func(h *Header, bi *big.Int) { h.baseFee = bi })
}

func FuzzHeaderStateLimitHash(f *testing.F) {
	fuzzHeaderUint64Hash(f,
		func(h *Header) uint64 { return h.stateLimit },
		func(h *Header, bi uint64) { h.stateLimit = bi })
}
func FuzzHeaderStateUsedHash(f *testing.F) {
	fuzzHeaderUint64Hash(f,
		func(h *Header) uint64 { return h.stateUsed },
		func(h *Header, bi uint64) { h.stateUsed = bi })
}
func FuzzHeaderQuaiStateSize(f *testing.F) {
	fuzzHeaderBigIntHash(f,
		func(h *Header) *big.Int { return h.quaiStateSize },
		func(h *Header, bi *big.Int) { h.quaiStateSize = bi })
}
func FuzzHeaderExtraHash(f *testing.F) {
	header, _ := headerTestData()
	f.Add(testByte)
	f.Add(header.extra)
	f.Fuzz(func(t *testing.T, b []byte) {
		localHeader, hash := headerTestData()
		if !bytes.Equal(localHeader.extra, b) {
			localHeader.extra = b
			require.NotEqual(t, localHeader.Hash(), hash, "Hash equal for extra \noriginal: %v, modified: %v", header.extra, b)
		}
	})
}

func FuzzHeaderExchangeRate(f *testing.F) {
	fuzzHeaderBigIntHash(f,
		func(h *Header) *big.Int { return h.exchangeRate },
		func(h *Header, bi *big.Int) { h.exchangeRate = bi })
}

func FuzzHeaderAvgTxFees(f *testing.F) {
	fuzzHeaderBigIntHash(f,
		func(h *Header) *big.Int { return h.avgTxFees },
		func(h *Header, bi *big.Int) { h.avgTxFees = bi })
}

func FuzzHeaderTotalFees(f *testing.F) {
	fuzzHeaderBigIntHash(f,
		func(h *Header) *big.Int { return h.totalFees },
		func(h *Header, bi *big.Int) { h.totalFees = bi })
}

func FuzzHeaderKQuaiDiscount(f *testing.F) {
	fuzzHeaderBigIntHash(f,
		func(h *Header) *big.Int { return h.kQuaiDiscount },
		func(h *Header, bi *big.Int) { h.kQuaiDiscount = bi })
}

func FuzzHeaderConversionFlowAmount(f *testing.F) {
	fuzzHeaderBigIntHash(f,
		func(h *Header) *big.Int { return h.conversionFlowAmount },
		func(h *Header, bi *big.Int) { h.conversionFlowAmount = bi })
}

func FuzzHeaderMinerDifficulty(f *testing.F) {
	fuzzHeaderBigIntHash(f,
		func(h *Header) *big.Int { return h.minerDifficulty },
		func(h *Header, bi *big.Int) { h.minerDifficulty = bi })
}

func FuzzHeaderPrimeStateRoot(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.primeStateRoot }, func(h *Header, hash common.Hash) { h.primeStateRoot = hash })
}

func FuzzHeaderRegionStateRoot(f *testing.F) {
	fuzzHeaderHash(f, func(h *Header) common.Hash { return h.regionStateRoot }, func(h *Header, hash common.Hash) { h.regionStateRoot = hash })
}
