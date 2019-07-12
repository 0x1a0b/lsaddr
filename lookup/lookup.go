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

	"github.com/booster-proj/lsaddr/lookup/internal"
)

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
	rgx, err := internal.BuildNFFilter(s)
	if err != nil {
		return []NetFile{}, err
	}

	log.Printf("regexp built: \"%s\"", rgx.String())

	ll, err := internal.OpenNetFiles(rgx)
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

// Hosts returns the list of source and destination addresses contained
// in `ff`.
func Hosts(ff []NetFile) (src, dst []net.Addr) {
	for _, v := range ff {
		src = append(src, v.Src)
		dst = append(dst, v.Dst)
	}
	return
}
