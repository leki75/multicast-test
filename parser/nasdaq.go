package parser

import (
	"encoding/binary"
	"fmt"

	"github.com/leki75/multicast-test/log"
	"github.com/leki75/multicast-test/net/multicast"
	"go.uber.org/zap"
)

func Nasdaq(n int) multicast.HandlerFunc {
	var logger func(string, ...zap.Field)
	expectedID := uint64(0)
	label := fmt.Sprintf("nasdaq#%d", n)

	return func(packetLength int, buff []byte) {
		sequenceID := binary.BigEndian.Uint64(buff[10:])
		blockCount := binary.BigEndian.Uint16(buff[18:])

		switch expectedID {
		case sequenceID, sequenceID + 1:
			logger = log.Logger.Info

		default:
			logger = log.Logger.Warn
		}

		logger(
			label,
			zap.Uint64("exp", expectedID),
			zap.Uint64("seq", sequenceID),
			zap.Uint16("cnt", blockCount),
			zap.Int("len", packetLength),
		)

		expectedID = sequenceID + uint64(blockCount)
	}
}
