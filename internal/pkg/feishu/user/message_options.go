package feishu

func (u *UserMessage) GetEnableOnP2MessageReceiveV1() bool {
	return u.EnableOnP2MessageReceiveV1.Load()
}

func (u *UserMessage) WithEnableOnP2MessageReadV1(enable bool) bool {
	u.EnableOnP2MessageReadV1.Store(enable)

	return u.EnableOnP2MessageReadV1.Load()
}

func (u *UserMessage) WithEnableOnP2MessageReceiveV1(enable bool) bool {
	u.EnableOnP2MessageReceiveV1.Store(enable)

	return u.EnableOnP2MessageReceiveV1.Load()
}

func (u *UserMessage) WithEnableOnP2UserCreatedV3(enable bool) bool {
	u.EnableOnP2UserCreatedV3.Store(enable)

	return u.EnableOnP2UserCreatedV3.Load()
}

func (u *UserMessage) WithVerificationToken(token string, isCallback ...bool) string {
	u.VerificationToken.Store(token)

	if len(isCallback) <= 0 {
		isCallback[0] = true
	}

	if isCallback[0] {
		u.init()
	}

	return u.VerificationToken.Load()
}

func (u *UserMessage) WithEncryptKey(secret string, isCallback ...bool) string {
	u.EncryptKey.Store(secret)

	if len(isCallback) <= 0 {
		isCallback[0] = true
	}

	if isCallback[0] {
		u.init()
	}

	return u.EncryptKey.Load()
}

func WithOnP2MessageReceiveV1(enable ...bool) UserMessageOption {
	return func(u *UserMessage) {
		isEnabled := true
		if len(enable) >= 1 {
			isEnabled = enable[0]
		}

		u.WithEnableOnP2MessageReceiveV1(isEnabled)
	}
}

func WithOnP2MessageReadV1(enable ...bool) UserMessageOption {
	return func(u *UserMessage) {
		isEnabled := true
		if len(enable) >= 1 {
			isEnabled = enable[0]
		}

		u.WithEnableOnP2MessageReadV1(isEnabled)
	}
}

func WithOnP2UserCreatedV3(enable ...bool) UserMessageOption {
	return func(u *UserMessage) {
		isEnabled := true
		if len(enable) >= 1 {
			isEnabled = enable[0]
		}
		u.EnableOnP2UserCreatedV3.Store(isEnabled)
	}
}

func WithEncryptKey(key string) UserMessageOption {
	return func(u *UserMessage) {
		u.EncryptKey.Store(key)
	}
}

func WithVerificationToken(token string, iscallback ...bool) UserMessageOption {
	return func(u *UserMessage) {
		u.VerificationToken.Store(token)

		if len(iscallback) >= 1 && !iscallback[0] {
			return
		}

		u.init()
	}
}
