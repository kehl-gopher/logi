package utils

import (
	"errors"
	"net/url"
	"strconv"
)

var ErrInvalidPort = errors.New("invalid port provided")

func PortResolver(port string) (int, error) {
	prt, err := strconv.Atoi(port)
	if err != nil {
		u, err := url.Parse(port)
		if err != nil {
			return 0, err
		}
		prt, err = strconv.Atoi(u.Port())
		if err != nil {
			return 0, err
		}
	}
	return prt, nil
}
