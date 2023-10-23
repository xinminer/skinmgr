package exec

import (
	"bufio"
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"os/exec"
	"skinmgr/internal/log"
	"strings"
)

func PlotCopy(f string, svr string) error {
	cmd := fmt.Sprintf("chia_plot_copy -d -t %s -- %s", svr, f)
	r, err := gproc.ShellExec(gctx.New(), cmd)
	if err != nil {
		return err
	}
	log.Log.Infof("chia plot copy result: %s", r)
	return nil
}

func PlotCopyCount(proc string) (int, error) {
	// 执行 ps 命令获取进程信息
	output, err := exec.Command("ps", "-ef").Output()
	if err != nil {
		return 0, err
	}
	// 逐行扫描 ps 命令的输出结果，查找目标进程名称
	count := 0
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, proc) {
			count = count + 1
		}
	}
	return count, nil
}
