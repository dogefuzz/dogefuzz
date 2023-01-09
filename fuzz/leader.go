package fuzz

import "github.com/dogefuzz/dogefuzz/pkg/common"

type FuzzerLeader interface {
	GetFuzzerStrategy(typ common.FuzzingType) (Fuzzer, error)
}

type fuzzerLeader struct {
	blackboxFuzzer        Fuzzer
	greyboxFuzzer         Fuzzer
	directedGreyboxFuzzer Fuzzer
}

func NewFuzzerLeader(e env) *fuzzerLeader {
	return &fuzzerLeader{
		blackboxFuzzer:        e.BlackboxFuzzer(),
		greyboxFuzzer:         e.GreyboxFuzzer(),
		directedGreyboxFuzzer: e.DirectedGreyboxFuzzer(),
	}
}

func (l *fuzzerLeader) GetFuzzerStrategy(typ common.FuzzingType) (Fuzzer, error) {
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
