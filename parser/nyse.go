package parser

import (
	"encoding/binary"
	"fmt"

	"github.com/leki75/multicast-test/log"
	"github.com/leki75/multicast-test/net/multicast"
	"go.uber.org/zap"
)

func Nyse(n int) multicast.HandlerFunc {
	expectedID := uint32(0)
	label := fmt.Sprintf("nyse#%d", n)

	// https://www.ctaplan.com/publicdocs/ctaplan/CTS_Pillar_Output_Specification.pdf
	return func(packetLength int, buf []byte) {
		sequenceID := binary.BigEndian.Uint32(buf[5:])
		blockCount := buf[9]

		var category []byte
		block := buf[20:]
		length := uint16(0)

		for i := blockCount; i > 0; i-- {
			length = binary.BigEndian.Uint16(block[0:])
			category = append(category, block[2:4]...)
			category = append(category, ',')
			block = block[length:]
		}

		logger := log.Logger.Info
		if expectedID < sequenceID {
			logger = log.Logger.Warn
		}

		logger(
			label,
			zap.Uint32("exp", expectedID),
			zap.Uint32("seq", sequenceID),
			zap.Uint8("cnt", blockCount),
			zap.Int("len", packetLength),
			zap.ByteString("typ", category),
		)

		expectedID = sequenceID + uint32(blockCount)
	}
}
