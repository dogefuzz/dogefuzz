package fuzz

import (
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

type fuzzerLeader struct {
	blackboxFuzzer        interfaces.Fuzzer
	greyboxFuzzer         interfaces.Fuzzer
	directedGreyboxFuzzer interfaces.Fuzzer
}

func NewFuzzerLeader(e env) *fuzzerLeader {
	return &fuzzerLeader{
		blackboxFuzzer:        e.BlackboxFuzzer(),
		greyboxFuzzer:         e.GreyboxFuzzer(),
		directedGreyboxFuzzer: e.DirectedGreyboxFuzzer(),
	}
}

func (l *fuzzerLeader) GetFuzzerStrategy(typ common.FuzzingType) (interfaces.Fuzzer, error) {
	switch typ {
	case common.BLACKBOX_FUZZING:
		return l.blackboxFuzzer, nil
	case common.GREYBOX_FUZZING:
		return l.greyboxFuzzer, nil
	case common.DIRECTED_GREYBOX_FUZZING:
		return l.directedGreyboxFuzzer, nil
	default:
		return nil, ErrFuzzerTypeNotFound
	}
}
