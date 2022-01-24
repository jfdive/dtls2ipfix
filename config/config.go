package config

import (
    "encoding/json"
    "os"

    "github.com/jfdive/dtls2ipfix/logging"
)

type ConfigForwarderEntry struct {
    Address string
    Port    int
}

type ConfigCtxt struct {
    Debug               bool
    LogFile             string
    LocalAddress        string
    LocalPort           int
    NumWorker           int
    ServerCertKey       string
    ServerCert          string
    ServerCa            string
    StatsFile           string
    StatsUpdateInterval string
    Forwarders          []ConfigForwarderEntry
}

var Config ConfigCtxt

func ReadConfig(configFile string) (*ConfigCtxt,error) {
    fh, err := os.Open(configFile)
    if err != nil {
        logging.Log.Errorf("Failed to read config file %s: %s", configFile, err.Error())
        return nil,err
    }
    defer fh.Close()
    jsonParser := json.NewDecoder(fh)
    err = jsonParser.Decode(&Config)
    if err != nil {
        logging.Log.Errorf("Failed to parse config file %s: %s", configFile, err.Error())
        return nil,err
    }
    logging.Log.Infof("Read configuration file %s", configFile)
    return &Config,nil
}

