package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	_ "github.com/miekg/dns"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

const DEBUG = false
const A uint16 = 0x01
const MaxThread int = 100

var tmpChain = make(chan struct{}, MaxThread)
var waitGroup sync.WaitGroup

func GetAName(DNSServer, domain string) bool {
	defer waitGroup.Done()
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	tmpChain <- struct{}{}
	in, err := dns.Exchange(msg, DNSServer+":53")
	<-tmpChain
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

Usage: gdd -dn Domain [-ds DNSServer(s)] [-dt DictFile]

Options:
`)
	flag.PrintDefaults()
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func createDNSServerList(DNSServers string) []string {
	dnsServers := strings.Split(DNSServers, ",")
	ret := make([]string, len(dnsServers))
	for i := 0; i < len(dnsServers); i++ {
		if net.ParseIP(dnsServers[i]) == nil {
			continue
		}
		ret[i] = dnsServers[i]
	}
	return ret
}

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Usage = usage

	dnsServer := flag.String("ds", "114.114.114.114", "DNS服务器(多个服务器用逗号[,]隔开)")
	domain := flag.String("dn", "", "域名")
	dict := flag.String("dt", "./dict.txt", "域名字典")
	flag.Parse()

	if *domain == "" {
		flag.Usage()
		return
	} else if *dict != "" && !isFileExist(*dict) {
		fmt.Printf("[x] 指定的字典文件 %s 不存在\n", *dict)
		return
	}

	f, err := os.Open(*dict)
	if err != nil {
		if DEBUG {
			log.Fatalf("[x] 打开字典文件错误：%s", err)
		} else {
			fmt.Println("[x] 打开字典文件错误。")
		}
		return
	}
	defer f.Close()

	lines := bufio.NewReader(f)
	dnsServers := createDNSServerList(*dnsServer)

	for {
		line, _, err := lines.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			if DEBUG {
				log.Fatalf("[x] 读取字典文件错误：%s", err)
			} else {
				fmt.Println("[x] 读取字典文件错误。")
			}
			return
		}
		domainPrefix := string(line)
		for i := 0; i < len(dnsServers); i++ {
			waitGroup.Add(1)
			go GetAName(dnsServers[i], fmt.Sprintf("%s.%s", strings.TrimPrefix(domainPrefix, "\n\r"), *domain))
		}
	}
	waitGroup.Wait()
}
