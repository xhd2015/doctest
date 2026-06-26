package assert

type Option interface {
	apply(*matchConfig)
}

type matchConfig struct {
	mode              matchMode
	normalizeNewlines bool
}

type matchMode int

const (
	matchExact matchMode = iota
	matchContains
)

type containsOption struct{}

func (containsOption) apply(cfg *matchConfig) {
	cfg.mode = matchContains
}

// Contains enables MatchContains mode: the template must appear as a contiguous
// ordered subregion of actual output.
func Contains() Option {
	return containsOption{}
}

type normalizeNewlinesOption struct {
	v bool
}

func (o normalizeNewlinesOption) apply(cfg *matchConfig) {
	cfg.normalizeNewlines = o.v
}

// NormalizeNewlines controls \r\n → \n normalization (default true).
func NormalizeNewlines(v bool) Option {
	return normalizeNewlinesOption{v: v}
}

func defaultMatchConfig() matchConfig {
	return matchConfig{
		mode:              matchExact,
		normalizeNewlines: true,
	}
}

func applyOptions(opts ...Option) matchConfig {
	cfg := defaultMatchConfig()
	for _, o := range opts {
		o.apply(&cfg)
	}
	return cfg
}