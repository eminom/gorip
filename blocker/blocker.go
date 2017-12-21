package blocker

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type RipIPList map[string]bool

func (r RipIPList) Dump() {
	log.Print("dump a rip-list **")
	for ent := range r {
		log.Print(ent)
	}
}

func (r RipIPList) IsAllowed(ip net.IP) bool {
	ips := ip.String()
	_, ok := r[ips]
	return ok
}

var (
	DefAl RipIPList
)

func init() {
	//_ = DefAl.LoadFromFile("iplist")
	decoded, err := decryptEncoded(allowedConstantString, hintKey)
	if err != nil {
		log.Fatal(err)
	}
	_ = DefAl.LoadFromBuffer(string(decoded))
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
		if 4 == iCount && iErr == nil {
			rvs[net.IPv4(byte(a),
				byte(b), byte(c), byte(d)).String()] = true
		}
	}
	*rl = rvs
	return nil
}
