package stats

import (
    "time"
    "sync/atomic"
    "encoding/json"
    "os"

    "github.com/jfdive/dtls2ipfix/config"
    "github.com/jfdive/dtls2ipfix/logging"
)

type Counter64 uint64

func (c *Counter64) Inc() uint64 {
    return atomic.AddUint64((*uint64)(c), 1)
}

func (c *Counter64) Get() uint64 {
    return atomic.LoadUint64((*uint64)(c))
}

type Dtls2IpfixStats_internal struct {
    InputPacket             Counter64
    ClientCount             Counter64
    OutputIpfixPackets      Counter64
    PacketMalformedError    Counter64
    DecodeError             Counter64
    ProcessingError         Counter64
    SendError               Counter64
}

var stats *Dtls2IpfixStats_internal

type Dtls2IpfixStats struct {
    InputPacket             uint64
    ClientCount             uint64
    OutputIpfixPackets      uint64
    PacketMalformedError    uint64
    DecodeError             uint64
    ProcessingError         uint64
    SendError               uint64
}

func Init() error {
    stats = new(Dtls2IpfixStats_internal)
    if config.Config.StatsFile != "" {
        go UpdateStatsFile()
    }
    return nil
}

func Get() *Dtls2IpfixStats_internal {
    if stats == nil {
        Init()
    }
    return stats
}

func dataGet(r *Dtls2IpfixStats) {
    r.InputPacket = stats.InputPacket.Get()
    r.ClientCount = stats.ClientCount.Get()
    r.OutputIpfixPackets = stats.OutputIpfixPackets.Get()
    r.PacketMalformedError = stats.PacketMalformedError.Get()
    r.DecodeError = stats.DecodeError.Get()
    r.ProcessingError = stats.ProcessingError.Get()
    r.SendError = stats.SendError.Get()
}

func UpdateStatsFile() {
    updateInterval := "10s"
    if config.Config.StatsUpdateInterval != "" {
        updateInterval = config.Config.StatsUpdateInterval
    }
    d,err := time.ParseDuration(updateInterval)
    if err != nil {
        logging.Log.Errorf("Failed to parse stats update interval: %s", err)
        p,_ := time.ParseDuration("10s")
        d = p
    }

    logging.Log.Debugf("update stats file starting")

    var outStats Dtls2IpfixStats

    for {
        // sleep for x time
        time.Sleep(d)
        // create|append file
        fd, err := os.OpenFile(config.Config.StatsFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
        if err != nil {
            logging.Log.Errorf("failed to open/create stat file %s: %s", config.Config.StatsFile, err)
            continue
        }
        // update stats
        dataGet(&outStats)
        // generate json
        jsonContent, err := json.Marshal(outStats)
        if err != nil {
            logging.Log.Errorf("failed to json marshal stats: %s", err)
            continue
        }
        // write stat file
        fd.Write(jsonContent)
        // close
        fd.Close()
    }
}
