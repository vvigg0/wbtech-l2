package shell

import "os"

func expand(expr *expression) {
	for ji := range expr.jobs {
		for ci := range expr.jobs[ji].pipeline {
			cmd := &expr.jobs[ji].pipeline[ci]
			for ai := range cmd.argv {
				cmd.argv[ai] = expandVars(cmd.argv[ai])
			}
			if cmd.in != "" {
				cmd.in = expandVars(cmd.in)
			}
			if cmd.out != "" {
				cmd.out = expandVars(cmd.out)
			}
		}
	}
}

func expandVars(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); {
		if s[i] != '$' {
			out = append(out, s[i])
			i++
			continue
		}
		if i+1 >= len(s) {
			out = append(out, '$')
			i++
			continue
		}
		j := i + 1
		if !isVarStart(s[j]) {
			out = append(out, '$')
			i++
			continue
		}
		j++
		for j < len(s) && isVarChar(s[j]) {
			j++
		}
		name := s[i+1 : j]
		out = append(out, []byte(os.Getenv(name))...)
		i = j
	}
	return string(out)
}

func isVarStart(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || b == '_'
}

func isVarChar(b byte) bool {
	return isVarStart(b) || (b >= '0' && b <= '9')
}
