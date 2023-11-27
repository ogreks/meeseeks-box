package feishu

import (
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
)

type UserMessageOption func(u *UserMessage)

type UserMessage struct {
	log *zap.Logger
	db  orm.Repo

	EnableOnP2MessageReceiveV1 bool
	EnableOnP2MessageReadV1    bool
	EnableOnP2UserCreatedV3    bool
}

func (u *UserMessage) WithOnP2MessageReceiveV1() UserMessageOption {
	return func(u *UserMessage) {
		u.EnableOnP2MessageReceiveV1 = true
	}
}

func (u *UserMessage) WithOnP2MessageReadV1() UserMessageOption {
	return func(u *UserMessage) {
		u.EnableOnP2MessageReadV1 = true
	}
}

func (u *UserMessage) WithOnP2UserCreatedV3() UserMessageOption {
	return func(u *UserMessage) {
		u.EnableOnP2UserCreatedV3 = true
	}
}

func NewUserMessage(log *zap.Logger, db orm.Repo, opts ...UserMessageOption) *UserMessage {
	u := &UserMessage{
		log: log,
		db:  db,
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

func (u *UserMessage) Register() {

}
