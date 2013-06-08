package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	ustr "github.com/metaleap/go-util/str"

	ob "github.com/openbase/ob-core"
	obsrv "github.com/openbase/ob-core/server"
)

func main() {
	dirPath := flag.String("hive", "", fmt.Sprintf("%s hive directory path to use.\nIf omitted, defaults to either current directory, or the path\nstored in the $%s environment variable ('%s').\n", ob.OB_TITLE, ob.ENV_OBHIVE, os.Getenv(ob.ENV_OBHIVE)))
	addr := flag.String("addr", ":23456", "TCP address to serve HTTP requests.\nSpecify ':http' for default HTTP port or ':https' for default HTTPS port\n")
	tlsCertFile := flag.String("tls_cert", "", "File name containing a certificate for HTTPS-serving via TLS.\nIf the certificate is signed by a certificate authority, tls_cert should be\nthe concatenation of the server's certificate followed by the CA's certificate.\n")
	tlsKeyFile := flag.String("tls_key", "", "File name containing a matching private key for TLS serving.\nFor HTTPS/TLS serving, both tls_cert and tls_key are required.")
	flag.Parse()
	if len(*dirPath) == 0 {
		*dirPath = ob.Hive.GuessDirPath(*dirPath)
	}
	if err := ob.Init(*dirPath); err != nil {
		log.Fatal(err)
	}
	defer ob.Dispose()
	proto := "http"
	if len(*tlsCertFile) > 0 && len(*tlsKeyFile) > 0 {
		proto = "https"
	}
	host, _ := os.Hostname()
	if len(host) == 0 {
		host = "localhost"
	}
	log.Printf("[LIVE]\t@ %s://%s%s", proto, host, ustr.StripSuffix(*addr, ":"+proto))
	log.Fatal(obsrv.ListenAndServe(*addr, *tlsCertFile, *tlsKeyFile))
}
