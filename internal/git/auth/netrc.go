package auth

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io/ioutil"
	"strings"
)

type netrcProvider struct {
	netrcPath string
	auths     map[string]*http.BasicAuth
}

func NewNetRCProvider(netrcPath string) (Provider, error) {
	p := &netrcProvider{
		netrcPath: netrcPath,
		auths:     make(map[string]*http.BasicAuth),
	}
	netrcLines, err := readNetrc(netrcPath)
	if err != nil {
		return nil, err
	}
	for _, line := range netrcLines {
		p.auths[line.machine] = &http.BasicAuth{
			Username: line.login,
			Password: line.password,
		}
	}
	return p, nil
}

func (n *netrcProvider) Get(host string) (*http.BasicAuth, error) {
	basicAuth, ok := n.auths[host]
	if !ok {
		return nil, ErrAuthNotFound
	}
	return basicAuth, nil
}

type netrcLine struct {
	machine  string
	login    string
	password string
}

func parseNetrc(data string) []netrcLine {
	// See https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html
	// for documentation on the .netrc format.
	var nrc []netrcLine
	var l netrcLine
	inMacro := false
	for _, line := range strings.Split(data, "\n") {
		if inMacro {
			if line == "" {
				inMacro = false
			}
			continue
		}

		f := strings.Fields(line)
		i := 0
		for ; i < len(f)-1; i += 2 {
			// Reset at each "machine" token.
			// “The auto-login process searches the .netrc file for a machine token
			// that matches […]. Once a match is made, the subsequent .netrc tokens
			// are processed, stopping when the end of file is reached or another
			// machine or a default token is encountered.”
			switch f[i] {
			case "machine":
				l = netrcLine{machine: f[i+1]}
			case "default":
				break
			case "login":
				l.login = f[i+1]
			case "password":
				l.password = f[i+1]
			case "macdef":
				// “A macro is defined with the specified name; its contents begin with
				// the next .netrc line and continue until a null line (consecutive
				// new-line characters) is encountered.”
				inMacro = true
			}
			if l.machine != "" && l.login != "" && l.password != "" {
				nrc = append(nrc, l)
				l = netrcLine{}
			}
		}

		if i < len(f) && f[i] == "default" {
			// “There can be only one default token, and it must be after all machine tokens.”
			break
		}
	}

	return nrc
}

func readNetrc(netrcPath string) ([]netrcLine, error) {
	data, err := ioutil.ReadFile(netrcPath)
	if err != nil {
		return nil, err
	}

	return parseNetrc(string(data)), nil
}
