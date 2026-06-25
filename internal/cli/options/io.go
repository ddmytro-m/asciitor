package options

import "io"

const stdioMarker = "-"

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }
