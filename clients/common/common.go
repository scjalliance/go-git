package common

import (
	"fmt"
	"net/url"
	"strings"

	"gopkg.in/sourcegraph/go-vcsurl.v1"
	"gopkg.in/src-d/go-git.v2/pktline"
)

const GitUploadPackServiceName = "git-upload-pack"

type Endpoint string

func NewEndpoint(url string) (Endpoint, error) {
	vcs, err := vcsurl.Parse(url)
	if err != nil {
		return "", err
	}

	link := vcs.Link()
	if !strings.HasSuffix(link, ".git") {
		link += ".git"
	}

	return Endpoint(link), nil
}

func (e Endpoint) Service(name string) string {
	return fmt.Sprintf("%s/info/refs?service=%s", e, name)
}

// Capabilities contains all the server capabilities
// https://github.com/git/git/blob/master/Documentation/technical/protocol-capabilities.txt
type Capabilities map[string][]string

func parseCapabilities(line string) Capabilities {
	values, _ := url.ParseQuery(strings.Replace(line, " ", "&", -1))

	return Capabilities(values)
}

// Supports returns true if capability is preent
func (r Capabilities) Supports(capability string) bool {
	_, ok := r[capability]
	return ok
}

// Get returns the values for a capability
func (r Capabilities) Get(capability string) []string {
	return r[capability]
}

// SymbolicReference returns the reference for a given symbolic reference
func (r Capabilities) SymbolicReference(sym string) string {
	if !r.Supports("symref") {
		return ""
	}

	for _, symref := range r.Get("symref") {
		parts := strings.Split(symref, ":")
		if len(parts) != 2 {
			continue
		}

		if parts[0] == sym {
			return parts[1]
		}
	}

	return ""
}

type GitUploadPackInfo struct {
	Capabilities Capabilities
	Branches     map[string]string
}

func NewGitUploadPackInfo(d *pktline.Decoder) (*GitUploadPackInfo, error) {
	info := &GitUploadPackInfo{}
	if err := info.read(d); err != nil {
		return nil, err
	}

	return info, nil
}

func (r *GitUploadPackInfo) read(d *pktline.Decoder) error {
	lines, err := d.ReadAll()
	if err != nil {
		return err
	}

	r.Branches = map[string]string{}
	for _, line := range lines {
		if !r.isValidLine(line) {
			continue
		}

		if r.Capabilities == nil {
			r.Capabilities = parseCapabilities(line)
			continue
		}

		r.readLine(line)
	}

	return nil
}

func (r *GitUploadPackInfo) isValidLine(line string) bool {
	if line[0] == '#' {
		return false
	}

	return true
}

func (r *GitUploadPackInfo) readLine(line string) {
	commit, branch := r.getCommitAndBranch(line)
	r.Branches[branch] = commit
}

func (r *GitUploadPackInfo) getCommitAndBranch(line string) (string, string) {
	parts := strings.Split(strings.Trim(line, " \n"), " ")
	if len(parts) != 2 {
		return "", ""
	}

	return parts[0], parts[1]
}