package headers

type WWWAuthenticateHeader struct {
	Method     string
	Parameters map[string]string
}

type CSeqHeader struct {
	CSeq   int
	Method string
}
