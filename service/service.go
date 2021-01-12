package service

import (
	"errors"
	"fmt"
	"strings"
)

const defDomain = "trellis"

// Service service basic info
type Service struct {
	Domain  string `json:"domain" yaml:"domain"`
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

func (p *Service) init() {
	if p.Domain == "" {
		p.Domain = defDomain
	}
}

// FullName Service full name
func (p *Service) FullName() string {
	p.init()
	return fmt.Sprintf("/%s/%s", ReplaceURL(p.Domain), ReplaceURL(p.Name))
}

// FullPath Service full path
func (p *Service) FullPath(ps ...string) string {
	p.init()

	ss := []string{ReplaceURL(p.Domain), ReplaceURL(p.Name), ReplaceURL(p.Version)}

	for _, s := range ps {
		ss = append(ss, ReplaceURL(s))
	}

	return "/" + strings.Join(ss, "/")
}

// ParseService parse a string to base service
func ParseService(s string) (*Service, error) {
	ss := strings.Split(s, "/")

	lenSS := len(ss)

	var bs *Service
	if lenSS == 3 {
		bs = &Service{
			Name:    ss[1],
			Version: ss[2],
		}
	} else if lenSS > 3 {
		bs = &Service{
			Domain:  ss[1],
			Name:    ss[2],
			Version: ss[3],
		}
	} else {
		return nil, errors.New("failed parse base service")
	}

	bs.init()

	return bs, nil
}

// ReplaceURL replace url
func ReplaceURL(str string) string {
	str = strings.ToLower(str)
	str = strings.Replace(str, ":", "_", -1)
	str = strings.Replace(str, "/", "_", -1)
	return str
}
