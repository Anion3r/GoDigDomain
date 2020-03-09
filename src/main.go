package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	_ "github.com/miekg/dns"
	"log"
	"os"
)

const DEBUG = true
const A uint16 = 0x01

var Result map[string]string

func GetAName(DNSServer, domain string) bool {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	in, err := dns.Exchange(msg, DNSServer+":53")
	if err != nil {
		if DEBUG {
			log.Printf("[x] %s", err)
		}
		return false
	}

	for _, answer := range in.Answer {
		if answer.Header().Rrtype == A {
			Result[answer.(*dns.A).A.String()] = domain
		}
	}
	return true
}

func usage() {
	fmt.Fprintf(os.Stderr,
		`GoDigDomain: A Domain Name Burst Tool

Usage: gdd -dn Domain [-dm DNSServer]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Usage = usage

	dnsServer := flag.String("ds", "114.114.114.114", "DNS服务器")
	domain := flag.String("dn", "", "域名")
	flag.Parse()

	if *domain == "" {
		flag.Usage()
		return
	}
	Result = make(map[string]string)
	GetAName(*dnsServer, *domain)
	for ip := range Result {
		fmt.Printf("[+] %s -> %s\n", Result[ip], ip)
	}
}
