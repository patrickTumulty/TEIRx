package core

import (
	"fmt"
	"pen_daemon/src/txlog"
	"pen_daemon/src/utils"
)

/**
 * Prepare environment. This may include creating directories or
 * doing anything else on the system to make sure we are ready to go
 */
func PrepareEnvironment() {
	err := utils.FilesEnsureDirExists("./images")
	if err != nil {
		msg := fmt.Sprintf("Error creating images/ directory: %s", err.Error())
		txlog.TxLogError(msg)
		panic(msg)
	}

	err = utils.FilesEnsureDirExists("./images/movies")
	if err != nil {
		msg := fmt.Sprintf("Error creating images/movies/ directory: %s", err.Error())
		txlog.TxLogError(msg)
		panic(msg)
	}
}
