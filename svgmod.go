package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yofu/svgmod/svgmod"
)

func tmpcopy(fn string) (*os.File, error) {
	f, err := os.Open(fn)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(fn)
	t, err := ioutil.TempFile(dir, fmt.Sprintf("tmp_%s_", filepath.Base(fn)))
	if err != nil {
		t.Close()
		return nil, err
	}
	_, err = io.Copy(t, f)
	t.Seek(0, os.SEEK_SET)
	return t, nil
}

func parseCommand(e, font string) ([]*svgmod.Command, error) {
	ctxt := strings.Split(e, ";")
	commands := make([]*svgmod.Command, len(ctxt))
	nc := 0
	for _, ct := range ctxt {
		c, err := svgmod.Parse(strings.TrimLeft(ct, " "), font)
		if err != nil {
			return commands, err
		}
		commands[nc] = c
		nc++
	}
	return commands[:nc], nil
}

func parseScript(fn, font string) ([]*svgmod.Command, error) {
	f, err := os.Open(fn)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	commands := make([]*svgmod.Command, 0)
	nc := 0
	s := bufio.NewScanner(f)
	for s.Scan() {
		ct := s.Text()
		if strings.HasPrefix(ct, "#") {
			continue
		}
		c, err := svgmod.Parse(strings.TrimLeft(ct, " "), font)
		if err != nil {
			return commands, err
		}
		commands = append(commands, c)
		nc++
	}
	return commands[:nc], s.Err()
}

func main() {
	e := flag.String("e", "", "command")
	s := flag.String("s", "", "script file name")
	rmtmp := flag.Bool("rmtmp", false, "remove tmp file")
	verbose := flag.Bool("v", false, "verbose")
	fontfamily := flag.String("ff", "", "font family")

	flag.Parse()

	if *e == "" && *s == "" {
		log.Fatal("no command")
		os.Exit(1)
	}

	n := flag.NArg()
	if n == 0 {
		log.Fatal("no input")
		os.Exit(1)
	}

	filenames := make([]string, 0)

	for _, a := range flag.Args() {
		fs, err := filepath.Glob(a)
		if err != nil {
			continue
		}
		filenames = append(filenames, fs...)
	}

	if len(filenames) == 0 {
		log.Fatal("no input")
		os.Exit(1)
	}

	var font string
	switch strings.ToLower(*fontfamily) {
	default:
		font="FreeSerif"
	case "serif", "mincho":
		font="FreeSerif"
	case "sans", "gothic":
		font="FreeSans"
	}

	var commands []*svgmod.Command
	if *e != "" {
		cs, err := parseCommand(*e, font)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		commands = cs
	} else if *s != "" {
		cs, err := parseScript(*s, font)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		commands = cs
	}

	for _, c := range commands {
		if *verbose {
			fmt.Println(c.Statement)
		}
		for _, fn := range filenames {
			if *verbose {
				fmt.Printf("\t%s\n", fn)
			}
			fdin, err := tmpcopy(fn)
			defer func() {
				fdin.Close()
				if *rmtmp {
					os.Remove(fdin.Name())
				}
			}()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fdout, err := os.Create(fn)
			defer fdout.Close()
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = c.Exec(fdin, fdout)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
