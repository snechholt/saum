package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const directoryExecFlag = "--saum-exec-in-directory"

// exec.Cmd cannot fchdir, and macOS /dev/fd directory entries cannot be used
// for chdir. This tiny exec trampoline pins the child to the checked inode.
func init() {
	if len(os.Args) < 3 || os.Args[1] != directoryExecFlag {
		return
	}
	if err := syscall.Fchdir(3); err != nil {
		fmt.Fprintln(os.Stderr, "pinned workspace:", err)
		os.Exit(126)
	}
	if cwd, err := os.Getwd(); err == nil {
		_ = os.Setenv("PWD", cwd)
	}
	if err := syscall.Exec(os.Args[2], os.Args[2:], os.Environ()); err != nil {
		fmt.Fprintln(os.Stderr, "workspace exec:", err)
		os.Exit(126)
	}
}
func directoryCommand(ctx context.Context, directory *os.File, program string, args ...string) (*exec.Cmd, error) {
	executable, err := exec.LookPath(program)
	if err != nil {
		return nil, err
	}
	self, err := os.Executable()
	if err != nil {
		return nil, err
	}
	command := exec.CommandContext(ctx, self, append([]string{directoryExecFlag, executable}, args...)...)
	command.ExtraFiles = []*os.File{directory}
	return command, nil
}
func isDirectoryTrampoline(identity string) bool {
	return strings.Contains(identity, directoryExecFlag)
}
