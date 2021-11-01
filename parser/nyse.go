package parser

import (
	"encoding/binary"
	"fmt"

	"github.com/leki75/multicast-test/log"
	"github.com/leki75/multicast-test/net/multicast"
	"go.uber.org/zap"
)

func Nyse(n int) multicast.HandlerFunc {
	var logger func(string, ...zap.Field)
	expectedID := uint32(0)
	label := fmt.Sprintf("nyse#%d", n)

	return func(packetLength int, buf []byte) {
		sequenceID := binary.BigEndian.Uint32(buf[5:])
		blockCount := buf[9]

		switch expectedID {
		case sequenceID, sequenceID + 1:
			logger = log.Logger.Info

		default:
			logger = log.Logger.Warn
		}

		logger(
			label,
			zap.Uint32("exp", expectedID),
			zap.Uint32("seq", sequenceID),
			zap.Uint8("cnt", blockCount),
			zap.Int("len", packetLength),
		)

		expectedID = sequenceID + uint32(blockCount)
	}
}
