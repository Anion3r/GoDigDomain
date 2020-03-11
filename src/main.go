package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	_ "github.com/miekg/dns"
	"io"
	"log"
	"os"
	"strings"
)

//const DEBUG = true
const A uint16 = 0x01

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
			fmt.Printf("[+] %s -> %s\n", domain, answer.(*dns.A).A.String())
		}
	}
	return true
}

func usage() {
	fmt.Fprintf(os.Stderr,
		`GoDigDomain: A Domain Name Burst Tool

Usage: gdd -dn Domain [-ds DNSServer]

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
	f, err := os.Open("E:/GoProject/DNSCheck/dic.txt")
	if err != nil {
		log.Fatalf("[x] 读取字典错误：%s", err)
	}
	defer f.Close()
	lines := bufio.NewReader(f)
	for {
		line, _, err := lines.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("[x] 读取文件错误：%s", err)
		}
		domainPrefix := string(line)
		GetAName(*dnsServer, fmt.Sprintf("%s.%s", strings.TrimPrefix(domainPrefix, "\n\r"), *domain))
	}

}
