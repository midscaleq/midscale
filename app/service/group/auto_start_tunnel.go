package group

import (
	"fmt"
	"midscale/midscale/app/data/mgr/wg/conf"
	"midscale/midscale/app/data/model/define"
	"midscale/midscale/deps/mskit/exec"
	"strings"
	"time"
)

func StopTunnel(groupName string) error {

	// wireguard.exe /uninstalltunnelservice ms-0704
	p := []string{define.WgExeUninstall, groupName}
	stdoutput, stderrout, err := exec.Run(define.WgExe, p)
	if err != nil {
		errstr0 := "does not exist as an installed service"
		if !strings.Contains(string(stderrout), errstr0) {
			return fmt.Errorf("exec.Run %v %v, stdoutput:%v, stderrout:%v, err: %v",
				define.WgExe, strings.Join(p, " "), string(stdoutput), string(stderrout), err)
		}
	} else {
		time.Sleep(time.Second * 6) // uninstall tunnel need some time.
	}

	return nil
}

func StartTunnel(groupName, nickName string) error {
	// e.g. wireguard.exe /installtunnelservice %appdata%+"/midscale/ms-0704.conf"
	confFileName, err := conf.ConfFileName(groupName, nickName)
	if err != nil {
		return err
	}
	p := []string{define.WgExeInstall, confFileName}
	stdoutput, stderrout, err := exec.Run(define.WgExe, p)
	if err != nil {
		return fmt.Errorf("exec.Run %v %v, stdoutput:%v, stderrout:%v, err: %v",
			define.WgExe, strings.Join(p, " "), string(stdoutput), string(stderrout), err)
	} else {
		time.Sleep(time.Second * 3) // starting tunnel need some time.
	}
	return nil
}
