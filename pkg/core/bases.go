package core

import (
	"crypto/md5"
	"io"
	"regexp"
	"sort"
)

type PrometheusAlert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations struct {
		Summary     string `json:"summary"`
		Description string `json:"description"`
	} `json:"annotations"`
	StartsAt     string `json:"startsAt"`
	EndsAt       string `json:"endsAt"`
	GeneratorURL string `json:"generatorURL"`
}

type PrometheusAlertManagerMessage struct {
	Receiver string             `json:"receiver"`
	Status   string             `json:"status"`
	Alerts   []*PrometheusAlert `json:"alerts"`
}

type Patterns struct {
	Labels map[string]*regexp.Regexp
}

type Sender interface {
	Send([]*PrometheusAlert)
}

type Receiver interface {
	Subscribe(sender Sender, patterns *Patterns)
	Receive(message []*PrometheusAlert)
}


func MakeHashOfSortingList(strings []string) string {
	sort.Strings(strings)
	hash := md5.New()
	for _, message := range strings {
		io.WriteString(hash, message)
	}
	return string(hash.Sum(nil))
}