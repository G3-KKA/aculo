package aculo

// TODO: ftp://

// File file://
// Address http(s)://
type Destination string

const schemaHTTP = "http://"
const schemaHTTPS = "https://"
const schemaFILE = "file://"

func (dst Destination) ValidHTTP() bool {
	if len(dst) < len(schemaHTTP) {
		return false
	}
	if dst[0:len(schemaHTTP)] != schemaHTTP {
		return false
	}
	return true
}
func (dst Destination) ValidHTTPS() bool {
	if len(dst) < len(schemaHTTPS) {
		return false
	}
	if dst[0:len(schemaHTTPS)] != schemaHTTPS {
		return false
	}
	return true
}
func (dst Destination) ValidFile() bool {
	if len(dst) < len(schemaFILE) {
		return false
	}
	if dst[0:len(schemaFILE)] != schemaFILE {
		return false
	}
	return true

}
