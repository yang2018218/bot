package netutil

import (
	"net"

	"github.com/ip2location/ip2location-go/v9"
)

var (
	db *ip2location.DB
)

func OpenDB(dbPath string) (*ip2location.DB, error) {
	var err error
	db, err = ip2location.OpenDB(dbPath)
	return db, err
}

func IPv4ToLocation(ip string) (results ip2location.IP2Locationrecord, err error) {
	results, err = db.Get_all(ip)
	return
}

func IsIpv4Private(ip string) bool {
	ipAddress := net.ParseIP(ip)
	return ipAddress.IsPrivate()
}
