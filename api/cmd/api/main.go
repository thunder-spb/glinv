package main

import (
	"flag"
	"log"

	"restapi/internal/cmdb"

	"github.com/vharitonsky/iniflags"
)

func main() {
	//var configPath string
	//flag.StringVar(&configPath, "config-path", "vault/secrets/pgsql-cmdb-creds", "Path to config file")

	addr := flag.String("addr", ":10011", "HTTP network address")
	dsn := flag.String("dsn", "postgres://postgres:password@192.168.10.229:5432/cmdb?sslmode=disable", "PostgreSQL data source name")
	logLevel := flag.String("logLevel", "debug", "Log level")
	iniflags.Parse()

	cfg := cmdb.Config{Addr: *addr, LogLevel: *logLevel, DSN: *dsn}

	config := cmdb.NewConfig(cfg)

	// // Parsing JSON config file
	// configFile, err := os.Open(configPath)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// defer func() {
	// 	cerr := configFile.Close()

	// 	if err == nil {
	// 		err = cerr
	// 	}
	// }()

	// jsonParser := json.NewDecoder(configFile)

	// if err := jsonParser.Decode(&config); err != nil {
	// 	log.Fatal(err)
	// }

	// Start REST API
	if err := cmdb.Start(config); err != nil {
		log.Fatal(err)
	}
}
