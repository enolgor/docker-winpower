package service

import (
	"flag"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"
)

var UpsURL *url.URL
var PostURL *url.URL
var Script string
var Rate time.Duration
var Timeout *time.Duration
var LogLevel slog.Level

func init() {
	var upsstr, poststr, timeoutstr, ratestr, loglevelstr string
	flag.StringVar(&upsstr, "u", "", "ups json url, e.g. http://localhost:8888/0/json")
	flag.StringVar(&poststr, "p", "", "post to remote url (optional)")
	flag.StringVar(&Script, "s", "", "run script after timeout (optional)")
	flag.StringVar(&timeoutstr, "t", "", "timeout to run script (optional, mandatory if -s), golang duration format")
	flag.StringVar(&ratestr, "r", "1s", "rate to poll ups status, golang duration format, default 1s")
	flag.StringVar(&loglevelstr, "l", "INFO", "log level, default WARN. Values: DEBUG, INFO, WARN")
	flag.Parse()
	upsstr = strings.TrimSpace(upsstr)
	poststr = strings.TrimSpace(poststr)
	Script = strings.TrimSpace(Script)
	timeoutstr = strings.TrimSpace(timeoutstr)
	ratestr = strings.TrimSpace(ratestr)
	loglevelstr = strings.TrimSpace(loglevelstr)
	if upsstr == "" {
		log.Println("-u is mandatory")
		flag.Usage()
		os.Exit(1)
	}
	if ratestr == "" {
		log.Println("-r is mandatory")
		flag.Usage()
		os.Exit(1)
	}
	if loglevelstr == "" {
		log.Println("-l is mandatory")
		flag.Usage()
		os.Exit(1)
	}
	if (Script != "" && timeoutstr == "") || (Script == "" && timeoutstr != "") {
		log.Println("-s and -t are mandatory if one is defined")
		flag.Usage()
		os.Exit(1)
	}
	var err error
	if timeoutstr != "" {
		var timeout time.Duration
		if timeout, err = time.ParseDuration(timeoutstr); err != nil {
			log.Printf("can't parse timeout %s, use golang format\n", timeoutstr)
			flag.Usage()
			os.Exit(1)
		} else {
			Timeout = &timeout
		}
	}
	if Rate, err = time.ParseDuration(ratestr); err != nil {
		log.Printf("can't parse rate %s, use golang format\n", ratestr)
		flag.Usage()
		os.Exit(1)
	}
	if UpsURL, err = url.Parse(upsstr); err != nil {
		log.Printf("can't parse ups url %s\n", upsstr)
		flag.Usage()
		os.Exit(1)
	}
	if poststr != "" {
		if PostURL, err = url.Parse(poststr); err != nil {
			log.Printf("can't parse post url %s\n", poststr)
			flag.Usage()
			os.Exit(1)
		}
	}
	switch loglevelstr {
	case "WARN":
		LogLevel = slog.LevelWarn
	case "INFO":
		LogLevel = slog.LevelInfo
	case "DEBUG":
		LogLevel = slog.LevelDebug
	default:
		log.Println("valid -l values are: WARN, INFO, DEBUG")
		flag.Usage()
		os.Exit(1)
	}
}
