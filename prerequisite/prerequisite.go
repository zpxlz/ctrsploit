package prerequisite

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/util"
)

type Interface interface {
	Check() error
	Output()
	GetSatisfied() bool
}
type Prerequisites []Interface

func (ps Prerequisites) Satisfied() (satisfied bool, err error) {
	satisfied = true
	for _, p := range ps {
		err = p.Check()
		if err != nil {
			return
		}
		p.Output()
		if err != nil {
			return
		}
		if !p.GetSatisfied() {
			satisfied = false
		}
	}
	return
}

type BasePrerequisite struct {
	Name      string
	Info      string
	Satisfied bool
}

func (p BasePrerequisite) GetSatisfied() bool {
	return p.Satisfied
}

// Output print prerequisite with colorful; should be used after p.Check().
func (p BasePrerequisite) Output() {
	jsonOutput := true
	if jsonOutput {

	} else {
		if true {
			log.Logger.Infof(
				"%s %s: %s",
				util.TitleWithBgWhiteBold(p.Name),
				util.ColorfulTickOrBallot(p.Satisfied),
				p.Info,
			)
		} else {
			log.Logger.Infof(
				"%s %t: %s",
				p.Name, p.Satisfied, p.Info,
			)
		}
	}
	return
}