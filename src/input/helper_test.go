package input

import (
	"bytes"
	"io"
	"strings"

	. "gopkg.in/check.v1"
)

type HelperSuite struct{}

var _ = Suite(&HelperSuite{})

func (s *HelperSuite) TestGetRecord(c *C) {
	h := newHelper(new(MockFormat))
	h.factories = []ReaderFactory{
		func() io.Reader {
			return bytes.NewBuffer([]byte("a\nb\n"))
		},
	}

	c.Assert(h.GetLine(), Equals, "a")
	c.Assert(h.GetLine(), Equals, "b")
	c.Assert(h.GetLine(), Equals, "")
	c.Assert(h.IsEOF(), Equals, true)
}

func (s *HelperSuite) TestGetRecordMultipleFactgories(c *C) {
	h := newHelper(new(MockFormat))
	h.factories = []ReaderFactory{
		func() io.Reader {
			return bytes.NewBuffer([]byte("a\nb\n"))
		},
		func() io.Reader {
			return nil
		},
		func() io.Reader {
			return bytes.NewBuffer([]byte("c\nd\n"))
		},
	}

	c.Assert(h.GetLine(), Equals, "a")
	c.Assert(h.IsEOF(), Equals, false)
	c.Assert(h.GetLine(), Equals, "b")
	c.Assert(h.IsEOF(), Equals, false)
	c.Assert(h.GetLine(), Equals, "c")
	c.Assert(h.IsEOF(), Equals, false)
	c.Assert(h.GetLine(), Equals, "d")
	c.Assert(h.IsEOF(), Equals, false)
	c.Assert(h.GetLine(), Equals, "")
	c.Assert(h.IsEOF(), Equals, true)
}

func (s *HelperSuite) TestGetRecordNonNewLineTerminated(c *C) {
	h := newHelper(new(MockFormat))
	h.factories = []ReaderFactory{
		func() io.Reader {
			return bytes.NewBuffer([]byte("a\nb"))
		},
		func() io.Reader {
			return bytes.NewBuffer([]byte("c\nd"))
		},
	}

	c.Assert(h.GetLine(), Equals, "a")
	c.Assert(h.IsEOF(), Equals, false)
	c.Assert(h.GetLine(), Equals, "b")
	c.Assert(h.IsEOF(), Equals, false)
	c.Assert(h.GetLine(), Equals, "c")
	c.Assert(h.IsEOF(), Equals, false)
	c.Assert(h.GetLine(), Equals, "d")
	c.Assert(h.IsEOF(), Equals, false)
	c.Assert(h.GetLine(), Equals, "")
	c.Assert(h.IsEOF(), Equals, true)
}

func (s *HelperSuite) TestGetRecordLongTokens(c *C) {
	long := strings.Repeat("0", 64*1024)
	h := newHelper(new(MockFormat))
	h.factories = []ReaderFactory{
		func() io.Reader {
			return bytes.NewBuffer([]byte(long + "\nb\n"))
		},
	}

	c.Assert(len(h.GetLine()), Equals, len(long))
	c.Assert(h.GetLine(), Equals, "b")
}
