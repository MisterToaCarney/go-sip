package headers

import "fmt"

func Via(transport, host, rport, branch string) Header {
	header := NewHeader("Via")
	if branch == "" {
		branch = GenBranch()
	}
	header.Parameters["branch"] = fmt.Sprint("z9hG4bK", branch)
	if rport == "" {
		header.Value = fmt.Sprint(transport, " ", host)
	} else {
		header.Value = fmt.Sprint(transport, " ", host, ":", rport)
	}
	return header
}

func To(user, host, tag string) Header {
	header := NewHeader("To")
	header.Value = fmt.Sprintf("<sip:%s@%s>", user, host)

	if tag == "generate_tag" {
		tag = GenTag()
	}

	if tag == "" {
		return header
	} else {
		header.Parameters["tag"] = tag
		return header
	}
}

func From(user, host, tag string) Header {
	header := NewHeader("From")
	header.Value = fmt.Sprintf("<sip:%s@%s>", user, host)
	if tag == "generate_tag" {
		tag = GenTag()
	}
	if tag == "" {
		return header
	} else {
		header.Parameters["tag"] = tag
		return header
	}
}

func CallID(callId string) Header {
	header := NewHeader("Call-ID")
	if callId == "" {
		callId = GenCallID()
	}
	header.Value = callId
	return header
}

func CSeq(seq int, method string) Header {
	header := NewHeader("CSeq")
	header.Value = fmt.Sprintf("%v %s", seq, method)
	return header
}

func MaxForwards(forwards int) Header {
	header := NewHeader("Max-Forwards")
	header.Value = fmt.Sprint(forwards)
	return header
}

func Contact(user, destAddr, destPort string) Header {
	header := NewHeader("Contact")
	header.Value = fmt.Sprintf("<sip:%s@%s:%s>", user, destAddr, destPort)
	header.Parameters["expires"] = "3600"
	header.Parameters["q"] = "0.5"
	return header
}

func Allow() Header {
	header := NewHeader("Allow")
	header.Value = "INVITE, ACK, CANCEL, OPTIONS, BYE, REFER, SUBSCRIBE, NOTIFY, INFO, PUBLISH, MESSAGE"
	return header
}

func ContentLength(length int) Header {
	header := NewHeader("Content-Length")
	header.Value = fmt.Sprint(length)
	return header
}

func ContentType(contentType string) Header {
	header := NewHeader("Content-Type")
	header.Value = contentType
	return header
}

func UserAgent() Header {
	header := NewHeader("User-Agent")
	header.Value = "TPhone v0.1"
	return header
}

func Accept(contentType string) Header {
	header := NewHeader("Accept")
	header.Value = contentType
	return header
}

func DigestAuthorization(username, realm, nonce, uri, response, algorithm string) Header {
	header := NewHeader("Authorization")
	header.Value = fmt.Sprintf(
		`Digest username="%s", realm="%s", nonce="%s", uri="%s", response="%s", algorithm=%s`,
		username, realm, nonce, uri, response, algorithm,
	)
	return header
}
