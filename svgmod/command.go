package svgmod

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Command struct {
	Name     string
	Statement string
	function func(*os.File, *os.File) error
}

func (c *Command) Exec(inp, otp *os.File) error {
	return c.function(inp, otp)
}

func CommandSubstitute(before, after string) (*Command, error) {
	c := new(Command)
	c.Name = "substitute"
	c.Statement = fmt.Sprintf("substitute: %s → %s", before, after)
	c.function = func(inp, otp *os.File) error {
		s := bufio.NewReader(inp)
		wr := bufio.NewWriter(otp)
		for {
			buf := make([]byte, 0)
			var err error
			for {
				var l []byte
				var pf bool
				l, pf, err = s.ReadLine()
				buf = append(buf, l...)
				if !pf {
					break
				}
			}
			r := []byte(strings.Replace(string(buf), before, after, -1))
			for {
				nn, er := wr.Write(r)
				if er == nil {
					break
				}
				r = r[nn:]
			}
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			wr.WriteByte('\n')
			wr.Flush()
		}
		return nil
	}
	return c, nil
}
