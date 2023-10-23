package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// checks logging flag if program is called as ./main.go -debug
func setLogger() {

	// logging settings
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// activated check
	log.Debug().Msg("Debugging activated.")
}

func main() {

	// logger
	setLogger()

	// configs
	HostAddr := "localhost"
	HostPort := "8081"
	PathCaCrt := "certs/certificates/ca.crt"
	PathServerPem := "certs/certificates/server.pem"
	PathServerKey := "certs/certificates/server.key"
	UrlPath := "/my-btc-usdt-order"

	// parse certificate configs
	var caPath string
	flag.StringVar(&caPath, "path", PathCaCrt, "CA certificates")
	cert, _ := tls.LoadX509KeyPair(PathServerPem, PathServerKey)
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		log.Error().Err(err).Msg("ioutil.ReadFile")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// configure TLS suite
	tlsConfig := tls.Config{
		RootCAs:          caCertPool,
		CurvePreferences: []tls.CurveID{tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
		MaxVersion:       tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
		},
		NextProtos:         []string{"http/1.1"},
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		// PreferServerCipherSuites: true,
	}

	// create http server with TLS config
	server := http.Server{
		Addr:      HostAddr + ":" + HostPort,
		TLSConfig: &tlsConfig,
	}

	// set server handler
	http.HandleFunc(UrlPath, response)

	// server start listening for https connections
	log.Info().Msg("HTTPS Server " + HostAddr + ":" + HostPort + UrlPath)
	err = server.ListenAndServeTLS(PathServerPem, PathServerKey)
	if err != nil {
		log.Error().Err(err).Msg("server.ListenAndServeTLS")
	}
}

// handler
func response(w http.ResponseWriter, r *http.Request) {

	// measure server response time
	start := time.Now()

	// parse request
	r.ParseForm()

	// set response headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)

	// create response body
	message := map[string]interface{}{
		"pair":          "BTCUSDT",
		"data":          "2022.04.27",
		"time":          "12:00:00",
		"volume":        "321654",
		"price":         "38002.2",
		"all time high": "660000.5",
		"24 high":       "396564.3",
		"personal data": map[string]string{
			"age":    "20",
			"income": "1,300,561 Euro",
		},
	}

	// serialize JSON to bytes
	bytePresentation, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("json.Marshal")
		return
	}

	// calculate the required amount of dummy data
	requiredDummyDataSize := 2048 - len(bytePresentation)

	// create dummy data
	dummyData := make([]byte, requiredDummyDataSize)
	for i := 0; i < requiredDummyDataSize; i++ {
		dummyData[i] = 'a' // or any other placeholder character
	}

	message["dummy_data"] = string(dummyData)

	// re-serialize the message with dummy data
	bytePresentationWithDummyData, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Error serializing message with dummy data")
		return
	}

	// write response
	w.Write(bytePresentationWithDummyData)

	// measure elapsed response time
	elapsed := time.Since(start)
	log.Debug().Str("time", elapsed.String()).Msg("local server response time took.")

}
