package level

import "math/rand"

const (
	defaultMaxLevel = 16
	defaultP        = 0.25
)

type Generator interface {
	GetLevel() int
	MaxLevel() int
}

type RandomGenerator struct {
	lvUpP  float64
	maxLvl int
}

type RandomOption func(rlg *RandomGenerator)

func WithMaxLevel(maxLevel int) RandomOption {
	return func(sl *RandomGenerator) {
		sl.maxLvl = maxLevel
	}
}

func WithCustomP(p float64) RandomOption {
	return func(sl *RandomGenerator) {
		sl.lvUpP = p
	}
}

func NewRandomLevelGenerator(opts ...RandomOption) *RandomGenerator {
	rg := &RandomGenerator{
		lvUpP:  defaultP,
		maxLvl: defaultMaxLevel,
	}

	for _, op := range opts {
		op(rg)
	}

	return rg

}

func (rg *RandomGenerator) GetLevel() int {
	lvl := 1
	r := rand.Float64()
	for r < rg.lvUpP && lvl < (rg.maxLvl-1) {
		lvl++
		r = rand.Float64()
	}

	return lvl
}

func (rg *RandomGenerator) MaxLevel() int {
	return rg.maxLvl
}
