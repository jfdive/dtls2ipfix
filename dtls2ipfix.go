package main

import (
    "os"
    "flag"

    "github.com/jfdive/dtls2ipfix/logging"
    "github.com/jfdive/dtls2ipfix/config"
    "github.com/jfdive/dtls2ipfix/stats"
)

func main() {
    /* XXX */
    /* runtime.GOMAXPROCS(runtime.NumCPU()) */

    /* starting */
    logging.Init()
    logging.Log.Infof("Starting")

    /* flags */
    cfgConfigFile := flag.String("config", "", "Config file path")
    cfgDebug := flag.Bool("debug", false, "Enable debug")
    cfgQuiet := flag.Bool("quiet", false, "dont print log on stdout")
    flag.Parse()

    if *cfgConfigFile == "" {
        logging.Log.Errorf("missing configuration file")
        flag.PrintDefaults()
        os.Exit(1)
    }

    /* config */
    cfg, err := config.ReadConfig(*cfgConfigFile)
    if err != nil {
        os.Exit(1)
    }

    /* logging */
    logging.SetOutput(cfg.LogFile, !*cfgQuiet)
    if (*cfgDebug || cfg.Debug) {
        logging.SetLevel(logging.DEBUG)
    } else {
        logging.SetLevel(logging.INFO)
    }

    /* init stats */
    err = stats.Init()
    if err != nil {
        logging.Log.Error(err)
        os.Exit(1)
    }

}
