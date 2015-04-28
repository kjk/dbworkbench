package ga_event

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const beaconURL = "http://www.google-analytics.com/collect"

var (
	MissingUserAgent     = errors.New("Missing user agent")
	MissingClientId      = errors.New("Missing client ID")
	MissingPagePath      = errors.New("Missing page path")
	MissingIp            = errors.New("Missing IP")
	MissingEventCategory = errors.New("Missing Event Category")
	MissingEventAction   = errors.New("Missing Event Action")
	NegativeValue        = errors.New("Negative Event Value")
)

func generateUUID(cid *string) error {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}

	b[8] = (b[8] | 0x80) & 0xBF
	b[6] = (b[6] | 0x40) & 0x4F
	*cid = hex.EncodeToString(b)
	return nil
}

type gaContext struct {
	GAID       string
	DomainName string
	Client     *http.Client
}

func NewGAContext(gaid string, domainName string) *gaContext {
	client := &http.Client{}
	return &gaContext{gaid, domainName, client}
}

func (gac *gaContext) logToGA(ua string, values url.Values) error {
	req, _ := http.NewRequest("POST", beaconURL, strings.NewReader(values.Encode()))

	if ua != "" {
		req.Header.Add("User-Agent", ua)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err := gac.Client.Do(req)

	if err != nil {
		log.Println("GAEvent Recording error", err)
		return err
	}

	return nil
}

type gaPageView struct {
	context   *gaContext
	UserAgent string
	Cid       string
	Ip        string
	PagePath  string
	PageTitle string
	Params    map[string]string
}

func (gac *gaContext) NewPageView(ua string, cid string, ip string, pagepath string, pagetitle string, params map[string]string) *gaPageView {
	return &gaPageView{gac, ua, cid, ip, pagepath, pagetitle, params}
}

func (gapv *gaPageView) Log() error {
	if gapv.Cid == "" {
		return MissingClientId
	}
	if gapv.Ip == "" {
		return MissingIp
	}
	if gapv.PagePath == "" {
		return MissingPagePath
	}
	if gapv.UserAgent == "" {
		return MissingUserAgent
	}

	payload := url.Values{
		"v":   {"1"},                     // protocol version = 1
		"t":   {"pageview"},              // hit type
		"tid": {gapv.context.GAID},       // tracking / property ID
		"cid": {gapv.Cid},                // unique client ID (server generated UUID)
		"dp":  {gapv.PagePath},           // page path
		"uip": {gapv.Ip},                 // IP address of the user
		"dh":  {gapv.context.DomainName}, // Domain name of site
	}

	if gapv.PageTitle != "" {
		payload["dt"] = []string{gapv.PageTitle}
	}

	for key, val := range gapv.Params {
		payload[key] = []string{val}
	}

	return gapv.context.logToGA(gapv.UserAgent, payload)
}

type gaEvent struct {
	context  *gaContext
	Cid      string
	Category string
	Action   string
	Label    string
	Value    string
	Params   map[string]string
}

func (gac *gaContext) NewEvent(cid string, category string, action string, label string, value string, params map[string]string) *gaEvent {
	return &gaEvent{gac, cid, category, action, label, value, params}
}

func (gae *gaEvent) Log() error {
	if gae.Cid == "" {
		return MissingClientId
	}
	if gae.Category == "" {
		return MissingEventCategory
	}
	if gae.Action == "" {
		return MissingEventAction
	}

	val, intErr := strconv.Atoi(gae.Value)
	if intErr != nil || val < 0 {
		return NegativeValue
	}

	payload := url.Values{
		"v":   {"1"},              // protocol version = 1
		"t":   {"event"},          // hit type
		"tid": {gae.context.GAID}, // tracking / property ID
		"cid": {gae.Cid},          // unique client ID (server generated UUID)
		"ec":  {gae.Category},     // Event Category
		"ea":  {gae.Action},       // Action taken
	}

	if gae.Label != "" {
		payload["el"] = []string{gae.Label}
	}

	if gae.Value != "" {
		payload["ev"] = []string{gae.Value}
	}

	for key, val := range gae.Params {
		payload[key] = []string{val}
	}

	return gae.context.logToGA("", payload)
}
