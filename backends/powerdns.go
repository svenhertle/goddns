package backends

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/stdlib"
	"github.com/spf13/viper"
)

// PowerDNSBackend is a GoDDNS backend for PowerDNS with PostgresSQL
type PowerDNSBackend struct {
	domain string
	ttl    int

	dbConnectionString string
	dbDriver           string

	pdnsDomainID int
}

// dbVendorDriverMap maps the supported database vendors to the used drivers
var dbVendorDriverMap = map[string]string{
	"postgres": "pgx",
}

// Configure reads the database config for the PowerDNS + PostgresSQL backend
func (backend *PowerDNSBackend) Configure(cfg *viper.Viper, domain string, ttl int) error {
	// store dns settings
	backend.domain = domain
	backend.ttl = ttl

	// set some defaults
	cfg.SetDefault("vendor", "postgres")
	cfg.SetDefault("host", "localhost")
	cfg.SetDefault("port", 5432)

	// validation
	if cfg.GetString("database") == "" {
		return errors.New("Database name missing")
	}
	if cfg.GetString("user") == "" {
		return errors.New("User missing")
	}
	if cfg.GetString("password") == "" {
		return errors.New("Password missing")
	}

	// check database vendor
	driver, found := dbVendorDriverMap[cfg.GetString("vendor")]
	if !found {
		return errors.New("Database vendor not supported")
	}
	backend.dbDriver = driver

	// force inclusion (TODO: get rid of this)
	_ = stdlib.ErrNotPgx

	// prepare db connection string
	backend.dbConnectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.GetString("host"), cfg.GetInt("port"), cfg.GetString("user"), cfg.GetString("password"),
		cfg.GetString("database"), cfg.GetString("sslmode"))
	// get domain id, we will need it
	pdnsDomainID, err := backend.getDomainID()
	if err != nil {
		return err
	}
	backend.pdnsDomainID = pdnsDomainID

	return nil
}

// Update writes a new record to the PostgresSQL database of PowerDNS
func (backend *PowerDNSBackend) Update(name string, ip string, addressType AddressType) error {
	// connect to database
	db, err := sql.Open(backend.dbDriver, backend.dbConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	// we need to insert the FQDN
	fqdn := fmt.Sprintf("%s.%s", name, backend.domain)

	// ideally we want to prevent race conditions between the select and the insert. unfortunately, a simple solution would
	// reqquire a change of the unique constraints on the PowerDNS database. we should not do this...
	// TODO: row based locking? golang mutex? table lock? ignore it?

	// check if entry exists already
	var recordID int
	var recordNotFound = false
	row := db.QueryRow("SELECT id FROM records WHERE domain_id=$1 AND name=$2 AND type=$3;", backend.pdnsDomainID, fqdn, addressType)
	err = row.Scan(&recordID)

	if err == sql.ErrNoRows {
		recordNotFound = true
	} else if err != nil {
		return err
	}

	// change record
	if recordNotFound {
		_, err = db.Exec("INSERT INTO records (domain_id, name, content, type, ttl) VALUES ($1, $2, $3, $4, $5);",
			backend.pdnsDomainID, fqdn, ip, addressType, backend.ttl)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec("UPDATE records SET content=$1, ttl=$2 WHERE id=$3;",
			ip, backend.ttl, recordID)
		if err != nil {
			return err
		}
	}

	return nil
}

// getDomainID reads the PowerDNS domain id from the database
func (backend *PowerDNSBackend) getDomainID() (int, error) {
	// connect
	db, err := sql.Open(backend.dbDriver, backend.dbConnectionString)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// query
	var id int
	row := db.QueryRow("SELECT id FROM domains WHERE name=$1;", backend.domain)
	err = row.Scan(&id)

	if err == sql.ErrNoRows {
		return 0, errors.New("Domain not found on PowerDNS")
	} else if err != nil {
		return 0, err
	}

	return id, nil
}
