package subkill3r

import (
	"errors"

	"github.com/miekg/dns"
)

func LookupA(fqdn, serverAddr string) ([]string, error) {
	var m dns.Msg
	var ips []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	in, err := dns.Exchange(&m, serverAddr)
	if err != nil {
		return ips, err
	}
	if len(in.Answer) < 1 {
		return ips, errors.New("no answer")
	}
	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String())
			return ips, nil
		}
	}
	return ips, nil
}

func LookupCNAME(fqdn, serverAddr string) ([]string, error) {
	var m dns.Msg
	var fqdns []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	in, err := dns.Exchange(&m, serverAddr)
	if err != nil {
		return fqdns, err
	}
	if len(in.Answer) < 1 {
		return fqdns, errors.New("no answer")
	}
	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns, a.Target)
			return fqdns, nil
		}
	}
	return fqdns, nil
}

func Lookup(fqdn, serverAddr string) []Result {
	var results []Result
	var cfqdn = fqdn //keeping the original
	for {
		cnames, err := LookupCNAME(cfqdn, serverAddr )
		if err == nil && len(cnames) > 0 {
			cfqdn = cnames[0]
			continue // Process the next CNAME
		}
		ips, err := LookupA(cfqdn, serverAddr)
		if err != nil {
			break // There are no A records for this hostname.
		}
		for _, ip := range ips {
			results = append(results, Result{IPAdress: ip, Hostname: fqdn})
		}
		break // All the results will be processed till this step
	}
	return results	
}