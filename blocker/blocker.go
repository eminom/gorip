package blocker

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

var (
	ipv6pattern = regexp.MustCompile(`^::(\d+)$`)
)

type RipIPList map[string]bool

func (r RipIPList) Dump() {
	log.Print("** dump a rip-list **")
	for ent := range r {
		log.Print(ent)
	}
	log.Print("** ri-list dumped **")
}

func (r RipIPList) IsAllowed(ip net.IP) bool {
	ips := ip.String()
	_, ok := r[ips]
	return ok
}

func (r RipIPList) IsEmpty() bool {
	return len(r) == 0
}

var (
	DefAl RipIPList
)

func init() {
	//_ = DefAl.LoadFromFile("iplist")
	decoded, err := decryptEncoded(allowedConstantString, hintKey)
	if err != nil {
		panic(err)
	}
	log.Print("load buffer for blocker")
	err = DefAl.LoadFromBuffer(string(decoded))
	if err != nil {
		log.Print("load for blocker error: %v", err)
		return
	}
	log.Print("blocker buffer load done")
}

func (rl *RipIPList) LoadFromBuffer(c string) error {
	input := bytes.NewBuffer([]byte(c))
	return rl.loadFromReader(bufio.NewReader(input))
}

func (rl *RipIPList) LoadFromFile(file string) error {
	fin, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fin.Close()
	r := bufio.NewReader(fin)
	return rl.loadFromReader(r)
}

func (rl *RipIPList) loadFromReader(r *bufio.Reader) error {
	rvs := make(RipIPList)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		var a, b, c, d int
		line = strings.TrimSpace(line)
		iCount, iErr := fmt.Sscanf(line, "%d.%d.%d.%d", &a, &b, &c, &d)
		if nil == iErr {
			if 4 == iCount {
				log.Printf("<%v> is loaded", line)
				rvs[net.IPv4(byte(a),
					byte(b), byte(c), byte(d)).String()] = true
			}
		} else {
			if ipv6pattern.MatchString(line) {
				rvs[line] = true
			}
		}
	}
	*rl = rvs
	return nil
}
