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

package encoder_test

import (
	"strings"
	"testing"

	"github.com/booster-proj/lsaddr/encoder"
)

func TestEncode_BPF(t *testing.T) {
	l := netFiles0
	var w strings.Builder
	if err := encoder.NewBPF(&w).Encode(l); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expOut := "host 192.168.0.61 and port 54104 or host ::1 and port 60051\n"
	if expOut != w.String() {
		t.Fatalf("Unexpected output: wanted \"%s\", found \"%s\"", expOut, w.String())
	}
}
