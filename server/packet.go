package server

import (
    "net"

//    "github.com/jfdive/dtls2ipfix/config"
    "github.com/jfdive/dtls2ipfix/logging"
)

func packetReceived(buffer []byte, length int, localAddr net.Addr, remoteAddr net.Addr) {
	logging.Log.Debugf("Received data from %s of length %d", remoteAddr.String(), length)
}
