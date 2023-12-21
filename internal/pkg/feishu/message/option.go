package message

type MessageHandleOptions func(m *MessageHandle)

func WithEncryptKey(encrypt string) MessageHandleOptions {
	return func(m *MessageHandle) {
		m.EncryptKey.Store(encrypt)
	}
}

// WithVerificationToken
func WithVerificationToken(token string) MessageHandleOptions {
	return func(m *MessageHandle) {
		m.VerificationToken.Store(token)
	}
}
