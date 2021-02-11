package common

import (
	"regexp"
	"strings"
)

var (
	ipRegex = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
)

const (
	IP           = "ip_address"
	CountryCode  = "country_code"
	Country      = "country"
	City         = "city"
	Latitude     = "latitude"
	Longitude    = "longitude"
	MysteryValue = "mystery_value"
)

func IsIpv4Regex(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	return ipRegex.MatchString(ipAddress)
}
