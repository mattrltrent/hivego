package hivego

import (
	"bytes"
	"encoding/binary"
	"time"
)

func opIdB(opName string) byte {
	id := getHiveOpId(opName)
	return byte(id)
}

func refBlockNumB(refBlockNumber uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, refBlockNumber)
	return buf
}

func refBlockPrefixB(refBlockPrefix uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, refBlockPrefix)
	return buf
}

func expTimeB(expTime string) ([]byte, error) {
	exp, err := time.Parse("2006-01-02T15:04:05", expTime)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(exp.Unix()))
	return buf, nil
}

func countOpsB(ops []hiveOperation) []byte {
	b := make([]byte, 5)
	l := binary.PutUvarint(b, uint64(len(ops)))
	return b[0:l]
}

func extensionsB() byte {
	return byte(0x00)
}

func appendVString(s string, b *bytes.Buffer) *bytes.Buffer {
	vBuf := make([]byte, 5)
	vLen := binary.PutUvarint(vBuf, uint64(len(s)))
	b.Write(vBuf[0:vLen])

	b.WriteString(s)
	return b
}

func appendVStringArray(a []string, b *bytes.Buffer) *bytes.Buffer {
	b.Write([]byte{byte(len(a))})
	for _, s := range a {
		appendVString(s, b)
	}
	return b
}

func appendVAsset(asset string, b *bytes.Buffer) error {
    parts := strings.Split(asset, " ")
    if len(parts) != 2 {
       return errors.New("invalid asset format: " + asset)
    }

    amountStr, symbol := parts[0], parts[1]

    // all tokens have precision 3 except for VESTS    precision := 3    if symbol == "VESTS" {
       precision = 6    }

    // convert to their old names for compatibility    switch symbol {
    case "HIVE":
       symbol = "STEEM"    case "HBD":
       symbol = "SBD"    }

    // convert to float and multiply by 10^precision    amount, err := strconv.ParseFloat(amountStr, 64)

    if err != nil {
       return err    }

    amount = amount * math.Pow10(precision)

    // write the amount as int64    err = binary.Write(b, binary.LittleEndian, int64(amount))

    if err != nil {
       return err    }

    // write the precision    b.WriteByte(byte(precision))

    // write the symbol NUL padded to 8 bits    for i := 0; i < 7; i++ {
       if i < len(symbol) {
          b.WriteByte(symbol[i])
       } else {
          b.WriteByte(byte(0))
       }
    }

    return nil
}

func serializeTx(tx hiveTransaction) ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(refBlockNumB(tx.RefBlockNum))
	buf.Write(refBlockPrefixB(tx.RefBlockPrefix))
	expTime, err := expTimeB(tx.Expiration)
	if err != nil {
		return nil, err
	}
	buf.Write(expTime)

	opsB, err := serializeOps(tx.Operations)
	if err != nil {
		return nil, err
	}
	buf.Write(opsB)
	buf.Write([]byte{extensionsB()})
	return buf.Bytes(), nil
}

func serializeOps(ops []hiveOperation) ([]byte, error) {
	var opsBuf bytes.Buffer
	opsBuf.Write(countOpsB(ops))
	for _, op := range ops {
		b, err := op.serializeOp()
		if err != nil {
			return nil, err
		}
		opsBuf.Write(b)
	}
	return opsBuf.Bytes(), nil
}

func (o voteOperation) serializeOp() ([]byte, error) {
	var voteBuf bytes.Buffer
	voteBuf.Write([]byte{opIdB(o.opText)})
	appendVString(o.Voter, &voteBuf)
	appendVString(o.Author, &voteBuf)
	appendVString(o.Permlink, &voteBuf)

	weightBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(weightBuf, uint16(o.Weight))
	voteBuf.Write(weightBuf)

	return voteBuf.Bytes(), nil
}

func (o customJsonOperation) serializeOp() ([]byte, error) {
	var jBuf bytes.Buffer
	jBuf.Write([]byte{opIdB(o.opText)})
	appendVStringArray(o.RequiredAuths, &jBuf)
	appendVStringArray(o.RequiredPostingAuths, &jBuf)
	appendVString(o.Id, &jBuf)
	appendVString(o.Json, &jBuf)

	return jBuf.Bytes(), nil
}

func (o claimRewardOperation) serializeOp() ([]byte, error) {
    var claimBuf bytes.Buffer    claimBuf.Write([]byte{39})
    appendVString(o.Account, &claimBuf)
    err := appendVAsset(o.RewardHIVE, &claimBuf)

    if err != nil {
       return nil, err    }

    err = appendVAsset(o.RewardHBD, &claimBuf)

    if err != nil {
       return nil, err    }

    err = appendVAsset(o.RewardVests, &claimBuf)

    if err != nil {
       return nil, err    }

    return claimBuf.Bytes(), nil
}
