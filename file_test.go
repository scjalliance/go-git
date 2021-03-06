package git

import (
	"os"

	"gopkg.in/src-d/go-git.v2/core"
	"gopkg.in/src-d/go-git.v2/formats/packfile"

	. "gopkg.in/check.v1"
)

type SuiteFile struct {
	repos map[string]*Repository
}

var _ = Suite(&SuiteFile{})

// create the repositories of the fixtures
func (s *SuiteFile) SetUpSuite(c *C) {
	fixtureRepos := [...]struct {
		url      string
		packfile string
	}{
		{"https://github.com/tyba/git-fixture.git", "formats/packfile/fixtures/git-fixture.ofs-delta"},
		{"https://github.com/cpcs499/Final_Pres_P", "formats/packfile/fixtures/Final_Pres_P.ofs-delta"},
	}
	s.repos = make(map[string]*Repository, 0)
	for _, fixRepo := range fixtureRepos {
		s.repos[fixRepo.url] = NewPlainRepository()

		d, err := os.Open(fixRepo.packfile)
		c.Assert(err, IsNil)

		r := packfile.NewReader(d)
		r.Format = packfile.OFSDeltaFormat

		_, err = r.Read(s.repos[fixRepo.url].Storage)
		c.Assert(err, IsNil)

		c.Assert(d.Close(), IsNil)
	}
}

var contentsTests = []struct {
	repo     string // the repo name as in localRepos
	commit   string // the commit to search for the file
	path     string // the path of the file to find
	contents string // expected contents of the file
}{
	{
		"https://github.com/tyba/git-fixture.git",
		"b029517f6300c2da0f4b651b8642506cd6aaf45d",
		".gitignore",
		`*.class

# Mobile Tools for Java (J2ME)
.mtj.tmp/

# Package Files #
*.jar
*.war
*.ear

# virtual machine crash logs, see http://www.java.com/en/download/help/error_hotspot.xml
hs_err_pid*
`,
	},
	{
		"https://github.com/tyba/git-fixture.git",
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"CHANGELOG",
		`Initial changelog
`,
	},
}

func (s *SuiteFile) TestContents(c *C) {
	for i, t := range contentsTests {
		commit, err := s.repos[t.repo].Commit(core.NewHash(t.commit))
		c.Assert(err, IsNil, Commentf("subtest %d: %v (%s)", i, err, t.commit))

		file, err := commit.File(t.path)
		c.Assert(err, IsNil)
		c.Assert(file.Contents(), Equals, t.contents, Commentf(
			"subtest %d: commit=%s, path=%s", i, t.commit, t.path))
	}
}

var linesTests = []struct {
	repo   string   // the repo name as in localRepos
	commit string   // the commit to search for the file
	path   string   // the path of the file to find
	lines  []string // expected lines in the file
}{
	{
		"https://github.com/tyba/git-fixture.git",
		"b029517f6300c2da0f4b651b8642506cd6aaf45d",
		".gitignore",
		[]string{
			"*.class",
			"",
			"# Mobile Tools for Java (J2ME)",
			".mtj.tmp/",
			"",
			"# Package Files #",
			"*.jar",
			"*.war",
			"*.ear",
			"",
			"# virtual machine crash logs, see http://www.java.com/en/download/help/error_hotspot.xml",
			"hs_err_pid*",
		},
	},
	{
		"https://github.com/tyba/git-fixture.git",
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"CHANGELOG",
		[]string{
			"Initial changelog",
		},
	},
}

func (s *SuiteFile) TestLines(c *C) {
	for i, t := range linesTests {
		commit, err := s.repos[t.repo].Commit(core.NewHash(t.commit))
		c.Assert(err, IsNil, Commentf("subtest %d: %v (%s)", i, err, t.commit))

		file, err := commit.File(t.path)
		c.Assert(err, IsNil)
		c.Assert(file.Lines(), DeepEquals, t.lines, Commentf(
			"subtest %d: commit=%s, path=%s", i, t.commit, t.path))
	}
}

var ignoreEmptyDirEntriesTests = []struct {
	repo   string // the repo name as in localRepos
	commit string // the commit to search for the file
}{
	{
		"https://github.com/cpcs499/Final_Pres_P",
		"70bade703ce556c2c7391a8065c45c943e8b6bc3",
		// the Final dir in this commit is empty
	},
}

// It is difficult to assert that we are ignoring an (empty) dir as even
// if we don't, no files will be found in it.
//
// At least this test has a high chance of panicking if
// we don't ignore empty dirs.
func (s *SuiteFile) TestIgnoreEmptyDirEntries(c *C) {
	for i, t := range ignoreEmptyDirEntriesTests {
		commit, err := s.repos[t.repo].Commit(core.NewHash(t.commit))
		c.Assert(err, IsNil, Commentf("subtest %d: %v (%s)", i, err, t.commit))

		for file := range commit.Tree().Files() {
			_ = file.Contents()
			// this would probably panic if we are not ignoring empty dirs
		}
	}
}
