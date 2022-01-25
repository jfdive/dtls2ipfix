package server

import (
    "crypto"
    "crypto/ecdsa"
    "crypto/rsa"
    "crypto/tls"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "io/ioutil"
    "net"
    "path/filepath"
    "strings"
    "fmt"

    "github.com/jfdive/dtls2ipfix/config"
    "github.com/jfdive/dtls2ipfix/logging"
    //"github.com/jfdive/dtls2ipfix/stats"

    "github.com/pion/dtls/v2"
)

/******************************************************************************
 * X509 certs business
 *****************************************************************************/

const bufSize = 8192

var (
    errBlockIsNotPrivateKey  = errors.New("block is not a private key, unable to load key")
    errUnknownKeyTime        = errors.New("unknown key time in PKCS#8 wrapping, unable to load key")
    errNoPrivateKeyFound     = errors.New("no private key found, unable to load key")
    errBlockIsNotCertificate = errors.New("block is not a certificate, unable to load certificates")
    errNoCertificateFound    = errors.New("no certificate found, unable to load certificates")
)

// LoadKeyAndCertificate reads certificates or key from file
func LoadKeyAndCertificate(keyPath string, certificatePath string) (*tls.Certificate, error) {
    privateKey, err := LoadKey(keyPath)
    if err != nil {
        return nil, err
    }

    certificate, err := LoadCertificate(certificatePath)
    if err != nil {
        return nil, err
    }

    certificate.PrivateKey = privateKey

    return certificate, nil
}

// LoadKey Load/read key from file
func LoadKey(path string) (crypto.PrivateKey, error) {
    rawData, err := ioutil.ReadFile(filepath.Clean(path))
    if err != nil {
        return nil, err
    }

    block, _ := pem.Decode(rawData)
    if block == nil || !strings.HasSuffix(block.Type, "PRIVATE KEY") {
        return nil, errBlockIsNotPrivateKey
    }

    if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
        return key, nil
    }

    if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
        switch key := key.(type) {
        case *rsa.PrivateKey, *ecdsa.PrivateKey:
            return key, nil
        default:
            return nil, errUnknownKeyTime
        }
    }

    if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
        return key, nil
    }

    return nil, errNoPrivateKeyFound
}

// LoadCertificate Load/read certificate(s) from file
func LoadCertificate(path string) (*tls.Certificate, error) {
    rawData, err := ioutil.ReadFile(filepath.Clean(path))
    if err != nil {
        return nil, err
    }

    var certificate tls.Certificate

    for {
        block, rest := pem.Decode(rawData)
        if block == nil {
            break
        }

        if block.Type != "CERTIFICATE" {
            return nil, errBlockIsNotCertificate
        }

        certificate.Certificate = append(certificate.Certificate, block.Bytes)
        rawData = rest
    }

    if len(certificate.Certificate) == 0 {
        return nil, errNoCertificateFound
    }

    return &certificate, nil
}

/******************************************************************************
 * Utils
 *****************************************************************************/
// Use fake connect to get local ip
func localIp() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:53")
    if err != nil {
        logging.Log.Errorf("failed to net.Dial to get local ip: %s", err)
        return  net.IP{}
    }
    defer conn.Close()
    return conn.LocalAddr().(*net.UDPAddr).IP
}

/******************************************************************************
 * Init and Run
 *****************************************************************************/
var dtlsConfig dtls.Config

func Init() error {

    serverCert, err := LoadKeyAndCertificate(config.Config.ServerCertKey,
                                             config.Config.ServerCert)
    if err != nil {
        return err
    }

    dtlsConfig.Certificates = append(dtlsConfig.Certificates, *serverCert)

    //caCert, err := LoadCertificate(config.Config.ServerCa)
    //if err != nil {
    //    return err
    //}

    //for _, theCaCert := range(caCert.Certificate) {
    //    dtlsConfig.Certificates = append(dtlsConfig.Certificates, theCaCert)
    //}

    return nil
}

func handleConnection(conn net.Conn) {
    buffer := make([]byte, 8192)
    for {
        l, err := conn.Read(buffer)
        if err!= nil {
            logging.Log.Errorf("connection read failed: %s", err)
            break
        }
        packetReceived(buffer, l, conn.LocalAddr(), conn.RemoteAddr())
    }
}

func Run() error {
    localAddressParsed := net.ParseIP(config.Config.LocalAddress)
    if localAddressParsed != nil {
        return fmt.Errorf("failed to parse IP: %s", config.Config.LocalAddress)
    }
    localAddr := &net.UDPAddr{
        IP:     localAddressParsed,
        Port :  config.Config.LocalPort,
    }

    listener, err := dtls.Listen("udp", localAddr, &dtlsConfig)
    if err != nil {
        return err
    }
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            logging.Log.Errorf("Failed to accept(): %s", err)
            continue
        }
        go handleConnection(conn)
    }
}
