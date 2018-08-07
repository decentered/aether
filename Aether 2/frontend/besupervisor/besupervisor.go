// Frontend > BackendSupervisor

// This package handles the supervisory tasks related to the backend this frontend is the admin (admin) of.

package besupervisor

import (
	"aether-core/services/globals"
	"aether-core/services/logging"
	"fmt"
	"os"
	"os/exec"
	// "time"
)

func StartLocalBackend() {
	// todo - this needs to be replaced with running the binary instead of using "go ..." command.
	// time.Sleep(1 * time.Second)
	// return // todo
	cmd := exec.Command("go", constructArgs()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "../../aether-core/backend"
	logging.Log(1, "Local backend being started")
	err := cmd.Run()
	if err != nil {
		// return err
		logging.Logf(1, "Local backend had an error. Err: %v", err)
	}
	logging.Log(1, "Local backend exited.")
	// return nil
}

func constructArgs() []string {
	fesrvaddr := "127.0.0.1"
	// fesrvport := globals.FrontendTransientConfig.FrontendServerPort
	fesrvport := globals.FrontendConfig.GetFrontendAPIPort()
	fePublicKey := globals.FrontendConfig.GetMarshaledFrontendPublicKey()
	backendLogginglevel := 1
	baseCmd := []string{"run", "main.go", "run"}
	baseCmd = append(baseCmd, fmt.Sprintf("--logginglevel=%d", backendLogginglevel))
	baseCmd = append(baseCmd, fmt.Sprintf("--adminfeaddr=%s:%d", fesrvaddr, fesrvport))
	baseCmd = append(baseCmd, fmt.Sprintf("--adminfepk=%s", fePublicKey))
	return baseCmd
}
