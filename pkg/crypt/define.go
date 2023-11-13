package crypt

type Hash uint

const (
	// Bcrypt is the default hashing algorithm
	MD5 Hash = iota
	SHA1
	SHA224
	SHA256
	SHA384
	SHA512
	SHA512_224
	SHA512_256
)

type Encode uint

const (
	String Encode = iota
	Hex
	Base64
)

type Secret uint

const (
	PKCS1 Secret = iota
	PKCS8
)

type Crypt uint

const (
	RSA Crypt = iota
)

type Padding uint

const (
	PaddingPKCS5 Padding = iota
	PaddingPKCS7
)

type Cipher uint

const (
	ECB = iota
	CBC
	CFB
	OFB
)
