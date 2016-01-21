package crashreport

import "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/docker/machine/libmachine/log"

type logger struct{}

func (d *logger) Printf(fmtString string, args ...interface{}) {
	log.Debugf(fmtString, args)
}
