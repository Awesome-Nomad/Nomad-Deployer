package config

const gitlabMaxSlugLength = 63

type GitProvider string

type GitConfig struct {
	DefaultRef      string      `hcl:"default_ref,optional"`
	DefaultProvider GitProvider `hcl:"default_provider,optional"`
}

func (p GitProvider) NormalizeSlug(slug string) string {
	switch string(p) {
	case "gitlab":
		return gitlabTruncateSlug(slug)
	default:
		return slug
	}
}

// gitlabTruncateSlug truncate slug. Maximum length of gitlab slug is 63 character
func gitlabTruncateSlug(slug string) string {
	if len(slug) > gitlabMaxSlugLength {
		return slug[:gitlabMaxSlugLength]
	}
	return slug
}
