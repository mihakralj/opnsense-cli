/*
The MIT License (MIT)

Copyright (c) 2020 Alexander Staubo.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package internal

type Operation int

const (
	OpDelete Operation = iota
	OpInsert
	OpUnchanged
)

type Hunk struct {
	LineNum int
	Line    string
	Operation
}

func Diff(s1, s2 []string) []Hunk {
	return newDiffer(s1, s2).computeHunks()
}

type differ struct {
	s1, s2      []string
	pairs       [][]int
	chunks      []Hunk
	prevLineNum int
}

func newDiffer(s1, s2 []string) *differ {
	lcs := computeLCS(s1, s2)

	pairs := [][]int{}
	for i := 0; i < len(lcs); i++ {
		pairs = append(pairs, []int{-1, -1})
	}

	for i, j := 0, 0; i < len(s1) && j < len(lcs); i++ {
		if s1[i] == lcs[j] {
			pairs[j][0] = i
			j++
		}
	}
	for i, j := 0, 0; i < len(s2) && j < len(lcs); i++ {
		if s2[i] == lcs[j] {
			pairs[j][1] = i
			j++
		}
	}

	return &differ{s1, s2, pairs, nil, 0}
}

func (d *differ) computeHunks() []Hunk {
	i1, i2 := 0, 0
	for i := 0; i < len(d.pairs); i++ {
		i1, i2 = d.checkInterval(i1, d.pairs[i][0], i2, d.pairs[i][1])
	}

	d.checkInterval(i1, len(d.s1), i2, len(d.s2))

	for i := d.prevLineNum; i < len(d.s1); i++ {
		d.chunks = append(d.chunks, Hunk{
			Operation: OpUnchanged,
			LineNum:   i,
			Line:      d.s1[i],
		})
	}

	if len(d.chunks) == 0 {
		return []Hunk{}
	}
	return d.chunks
}

func (d *differ) checkInterval(i1 int, n1 int, i2 int, n2 int) (int, int) {
	for i1 < n1 || i2 < n2 {
		p1 := i1 < n1
		p2 := i2 < n2

		if p1 || p2 {
			if skip := i1 - d.prevLineNum; skip > 0 {
				for i := d.prevLineNum; i < i1; i++ {
					d.chunks = append(d.chunks, Hunk{
						Operation: OpUnchanged,
						LineNum:   i,
						Line:      d.s1[i],
					})
				}
				d.prevLineNum = i1
			}
		}

		if p1 {
			d.chunks = append(d.chunks, Hunk{
				Operation: OpDelete,
				LineNum:   i1,
				Line:      d.s1[i1],
			})
			i1++
		}

		if p2 {
			d.chunks = append(d.chunks, Hunk{
				Operation: OpInsert,
				LineNum:   i2,
				Line:      d.s2[i2],
			})
			i2++
		}

		if p1 || p2 {
			d.prevLineNum = i1
		}
	}

	return i1 + 1, i2 + 1
}

func computeLCS(a, b []string) []string {
	aLen := len(a)
	bLen := len(b)
	lengths := make([][]int, aLen+1)
	for i := 0; i <= aLen; i++ {
		lengths[i] = make([]int, bLen+1)
	}

	for i := 0; i < aLen; i++ {
		for j := 0; j < bLen; j++ {
			switch {
			case a[i] == b[j]:
				lengths[i+1][j+1] = lengths[i][j] + 1
			case lengths[i+1][j] > lengths[i][j+1]:
				lengths[i+1][j+1] = lengths[i+1][j]
			default:
				lengths[i+1][j+1] = lengths[i][j+1]
			}
		}
	}

	s := make([]string, 0, lengths[aLen][bLen])
	for x, y := aLen, bLen; x != 0 && y != 0; {
		switch {
		case lengths[x][y] == lengths[x-1][y]:
			x--
		case lengths[x][y] == lengths[x][y-1]:
			y--
		default:
			s = append(s, a[x-1])
			x--
			y--
		}
	}

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func PruneContext(hunks []Hunk, context int) []Hunk {
	result := make([]Hunk, 0, len(hunks))

	var buf []Hunk
	for _, h := range hunks {
		if h.Operation == OpUnchanged {
			buf = append(buf, h)
		} else {
			if c := min(len(buf), context); c > 0 {
				result = append(result, buf[len(buf)-c:]...)
			}
			buf = buf[0:0]
			result = append(result, h)
		}
	}
	if c := min(len(buf), context); c > 0 {
		result = append(result, buf[0:c]...)
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}