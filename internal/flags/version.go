package flags

import (
	"github.com/jc-lab/intel-amt-host-api/pkg/utils"
)

func (f *Flags) handleVersionCommand() utils.ReturnCode {
	if err := f.versionCommand.Parse(f.commandLineArgs[2:]); err != nil {
		return utils.IncorrectCommandLineParameters
	}
	// runs locally
	f.Local = true
	return utils.Success
}
