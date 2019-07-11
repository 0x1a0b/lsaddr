// Copyright © 2019 booster authors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package lookup

import (
	"log"
	"net"
	"os"
	"io"
	"regexp"
	"bytes"

	"github.com/booster-proj/lsaddr/lookup/internal"
	"gopkg.in/pipe.v2"
)

var Logger = log.New(os.Stderr, "[lookup] ", 0)

// NetFile contains some information obtained from a network file.
type NetFile struct {
	Command string   // command owning the file
	Src     net.Addr // source address
	Dst     net.Addr // destination address
}

// OpenNetFiles compiles a regular expression out of "s". Some manipulation
// may be performed on "s" before it is compiled, depending on the hosting
// operating system: on macOS for example, if "s" ends with ".app", it
// will be trated as the root path to an application, otherwise "s" will be
// compiled untouched.
// It then uses ``lsof'' (or its platform dependent equivalent) tool to find
// the list of open files, filtering the list taking only the lines that
// match against the regular expression built.
func OpenNetFiles(s string) ([]NetFile, error) {
	rgx, err := buildRgx(s)
	if err != nil {
		return []NetFile{}, err
	}
	Logger.Printf("regexp built: \"%s\"", rgx.String())

	ll, err := openNetFiles(rgx)
	if err != nil {
		return []NetFile{}, err
	}

	// map ``internal.OpenFile'' to ``NetFile''
	ff := make([]NetFile, len(ll))
	for i, v := range ll {
		src, dst := v.UnmarshalName()
		ff[i] = NetFile{
			Command: v.Command,
			Src:     src,
			Dst:     dst,
		}
	}
	return ff, nil
}

// HostsDedup returns the list of source and destination addresses contained
// in `ff`.
func Hosts(ff []NetFile) (src, dst []net.Addr) {
	for _, v := range ff {
		src = append(src, v.Src)
		dst = append(dst, v.Dst)
	}
	return
}

type lsofDecoderFunc func(io.Reader) ([]*internal.OpenFile, error)

// Private helpers

// buildRgx compiles a regular expression out of "s". Some manipulation
// may be performed on "s" before it is compiled, depending on the hosting
// operating system: on macOS for example, if "s" ends with ".app", it
// will be trated as the root path to an application.
func buildRgx(s string) (*regexp.Regexp, error) {
	expr, err := prepareExpr(s)
	if err != nil {
		return nil, err
	}
	rgx, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	return rgx, nil
}

// openNetFiles uses ``lsof'' (or its platform dependent equivalent) to find
// the list of open network files. It then filters the result using "rgx":
// each line that does not match is discarded.
func openNetFiles(rgx *regexp.Regexp) ([]*internal.OpenFile, error) {
	dec := lsofDecoder()

	p := pipe.Line(
		lsofCmd(),
		pipe.Filter(func(line []byte) bool {
			return rgx.Match(line)
		}),
	)
	output, err := pipe.Output(p)
	if err != nil {
		return []*internal.OpenFile{}, err
	}

	buf := bytes.NewBuffer(output)
	return dec(buf)
}
