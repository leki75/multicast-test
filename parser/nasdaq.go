package parser

import (
	"encoding/binary"
	"fmt"

	"github.com/leki75/multicast-test/log"
	"github.com/leki75/multicast-test/net/multicast"
	"go.uber.org/zap"
)

func Nasdaq(n int) multicast.HandlerFunc {
	expectedID := uint64(0)
	label := fmt.Sprintf("nasdaq#%d", n)

	// https://www.utpplan.com/DOC/UtpBinaryOutputSpec.pdf
	return func(packetLength int, buf []byte) {
		sequenceID := binary.BigEndian.Uint64(buf[10:])
		blockCount := binary.BigEndian.Uint16(buf[18:])

		var category []byte
		block := buf[20:]
		length := uint16(0)

		if blockCount == 0xFFFF {
			goto log
		}

		for i := blockCount; i > 0; i-- {
			length = binary.BigEndian.Uint16(block[0:])
			category = append(category, block[3:5]...)
			category = append(category, ',')
			block = block[length+2:]
		}

	log:
		logger := log.Logger.Info
		if expectedID < sequenceID {
			logger = log.Logger.Warn
		}

		logger(
			label,
			zap.Uint64("exp", expectedID),
			zap.Uint64("seq", sequenceID),
			zap.Uint16("cnt", blockCount),
			zap.Int("len", packetLength),
			zap.ByteString("typ", category),
		)

		expectedID = sequenceID + uint64(blockCount)
	}
}
