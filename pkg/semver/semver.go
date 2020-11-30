package semver

import (
	msemver "github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5/plumbing"
	"sort"
)

type TagsSemver struct {
	msvers  []*msemver.Version
	sverMap map[string]string
}

func NewTagsSemver() *TagsSemver {
	return &TagsSemver{
		sverMap: make(map[string]string),
	}
}

func (s *TagsSemver) AddTag(tagAnnotation string) {
	v, err := msemver.NewVersion(tagAnnotation)
	if err == nil {
		s.msvers = append(s.msvers, v)
		s.sverMap[v.String()] = tagAnnotation
	}
}

func (s *TagsSemver) GetLatestTag() string {
	if len(s.msvers) == 0 {
		return plumbing.Master.String()
	}
	sort.Sort(msemver.Collection(s.msvers))
	latestVersion := s.msvers[len(s.msvers)-1]
	return s.sverMap[latestVersion.String()]
}
