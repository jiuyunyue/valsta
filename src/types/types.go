package types

import (
	"fmt"
)

type VoterInfo struct {
	Address       string
	VoteTimes     uint64
	VoteProposals []uint64
}

func NewVoterInfo(address string, voteTimes uint64, voteProposals []uint64) VoterInfo {
	return VoterInfo{
		Address:       address,
		VoteTimes:     voteTimes,
		VoteProposals: voteProposals,
	}
}

func (v VoterInfo) String() string {
	return fmt.Sprintf("Address: %v, Times : %v, Proposals : %v \n", v.Address, v.VoteTimes, v.VoteProposals)
}

type ValidatorInfo struct {
	Address    string
	AccAddress string
	Jailed     bool
	Times      int64
	SurRate    string
}

func (v ValidatorInfo) String() string {
	return fmt.Sprintf("Address: %s, Times : %d, ValAddress: %s, Jailed : %v, SurRate : %v\n", v.Address, v.Times, v.AccAddress, v.Jailed, v.SurRate)
}

type Uptime map[string]ValidatorInfo

func (s Uptime) String() string {
	str := ""
	for k, v := range s {
		str += fmt.Sprintf(
			"Validator Address: %s , Times : %d  , ValAddress: %s ,Jailed : %t ,SurRate: %v \n",
			k,
			v.Times,
			v.AccAddress,
			v.Jailed,
			v.SurRate,
		)
	}
	return str
}

// ValidatorList true : unJail ,false : jailed
type ValidatorList map[string]bool

func (j ValidatorList) String() string {
	str := ""
	for k, v := range j {
		str += fmt.Sprintf("address: %s ,has jailed : %t \n", k, !v)
	}
	return str
}

// Jailed , true : Jailed ,false : unJailed
type Jailed map[string]bool

func (j Jailed) String() string {
	str := ""
	for k, v := range j {
		str += fmt.Sprintf("address: %s ,has jailed : %t \n", k, v)
	}
	return str
}
