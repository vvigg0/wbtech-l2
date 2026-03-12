package shell

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

type redirs struct {
	in  io.Reader
	out io.Writer
	cl  func()
}

func setupBuiltinRedirs(c command, defaultIn io.Reader, defaultOut io.Writer) (redirs, error) {
	r := redirs{
		in:  defaultIn,
		out: defaultOut,
		cl:  func() {},
	}

	var files []*os.File
	r.cl = func() {
		for _, f := range files {
			_ = f.Close()
		}
	}

	if c.in != "" {
		f, err := os.Open(c.in)
		if err != nil {
			r.cl()
			return redirs{}, err
		}
		files = append(files, f)
		r.in = f
	}

	if c.out != "" {
		f, err := os.OpenFile(c.out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			r.cl()
			return redirs{}, err
		}
		files = append(files, f)
		r.out = f
	}

	return r, nil
}

func runExpression(expr expression, out, errw io.Writer) (int, *exitRequest) {
	if len(expr.jobs) == 0 {
		return 0, nil
	}

	status, exit := runJob(expr.jobs[0], out, errw)
	if exit != nil {
		return 0, exit
	}
	for i, op := range expr.ops {
		next := expr.jobs[i+1]
		switch op {
		case opAnd:
			if status == 0 {
				status, exit = runJob(next, out, errw)
				if exit != nil {
					return 0, exit
				}
			}
		case opOr:
			if status != 0 {
				status, exit = runJob(next, out, errw)
				if exit != nil {
					return 0, exit
				}
			}
		}
	}
	return status, nil
}

func runJob(job job, out, errw io.Writer) (int, *exitRequest) {
	if len(job.pipeline) == 0 {
		return 1, nil
	}
	if len(job.pipeline) == 1 {
		c := job.pipeline[0]

		rd, err := setupBuiltinRedirs(c, os.Stdin, out)
		if err != nil {
			fmt.Fprintln(errw, err)
			return 1, nil
		}
		defer rd.cl()

		handled, code, exit := runBuiltin(c, rd.out, errw)
		if exit != nil {
			return 0, exit
		}
		if handled {
			return code, nil
		}

		code, e := runExternal(c, out, errw, nil, nil)
		if e != nil {
			fmt.Fprintln(errw, e)
		}
		return code, nil
	}
	code, err := runPipeline(job.pipeline, out, errw)
	if err != nil {
		fmt.Fprintln(errw, err)
	}
	return code, nil
}

func runPipeline(cmds []command, out, errw io.Writer) (int, error) {
	n := len(cmds)
	if n == 0 {
		return 1, errors.New("empty pipeline")
	}

	procs := make([]*exec.Cmd, n)
	for i := range n {
		if len(cmds[i].argv) == 0 {
			return 1, errors.New("empty command")
		}
		procs[i] = exec.Command(cmds[i].argv[0], cmds[i].argv[1:]...)
		procs[i].Stderr = errw
	}

	type pp struct{ r, w *os.File }
	pipes := make([]pp, 0, n-1)
	for range n - 1 {
		r, w, err := os.Pipe()
		if err != nil {
			for _, p := range pipes {
				_ = p.r.Close()
				_ = p.w.Close()
			}
			return 1, err
		}
		pipes = append(pipes, pp{r: r, w: w})
	}

	var files []*os.File
	closeAll := func() {
		for _, p := range pipes {
			_ = p.r.Close()
			_ = p.w.Close()
		}
		for _, f := range files {
			_ = f.Close()
		}
	}

	for i := range n {
		if cmds[i].in != "" {
			f, err := os.Open(cmds[i].in)
			if err != nil {
				closeAll()
				return 1, err
			}
			files = append(files, f)
			procs[i].Stdin = f
		} else if i == 0 {
			procs[i].Stdin = os.Stdin
		} else {
			procs[i].Stdin = pipes[i-1].r
		}

		if cmds[i].out != "" {
			f, err := os.OpenFile(cmds[i].out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				closeAll()
				return 1, err
			}
			files = append(files, f)
			procs[i].Stdout = f
		} else if i == n-1 {
			procs[i].Stdout = out
		} else {
			procs[i].Stdout = pipes[i].w
		}
	}

	started := 0
	for i := range n {
		if err := procs[i].Start(); err != nil {
			for j := range started {
				_ = procs[j].Process.Kill()
			}
			closeAll()
			return 1, err
		}
		started++
	}

	closeAll()

	lastStatus := 0
	for i := range n {
		err := procs[i].Wait()
		if i == n-1 {
			if err == nil {
				lastStatus = 0
			} else {
				var ee *exec.ExitError
				if errors.As(err, &ee) {
					if ws, ok := ee.Sys().(syscall.WaitStatus); ok {
						lastStatus = ws.ExitStatus()
					} else {
						lastStatus = 1
					}
				} else {
					lastStatus = 1
				}
			}
		}
	}

	return lastStatus, nil
}

func runExternal(c command, out, errw io.Writer, stdin io.Reader, stdout io.Writer) (int, error) {
	cmd := exec.Command(c.argv[0], c.argv[1:]...)
	cmd.Stderr = errw

	if stdin != nil {
		cmd.Stdin = stdin
	} else {
		cmd.Stdin = os.Stdin
	}
	if stdout != nil {
		cmd.Stdout = stdout
	} else {
		cmd.Stdout = out
	}

	if c.in != "" {
		f, err := os.Open(c.in)
		if err != nil {
			return 1, err
		}
		defer f.Close()
		cmd.Stdin = f
	}

	if c.out != "" {
		f, err := os.OpenFile(c.out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return 1, err
		}
		defer f.Close()
		cmd.Stdout = f
	}

	err := cmd.Run()
	if err == nil {
		return 0, nil
	}

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		if ws, ok := ee.Sys().(syscall.WaitStatus); ok {
			return ws.ExitStatus(), nil
		}
		return 1, nil
	}

	return 1, err
}
