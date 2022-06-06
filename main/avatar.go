package main

import (
	"errors"
)

// AvatarインスタンスがアバターURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

// ユーザーのプロフィール画像を表す型
type Avatar interface {
	// 指定されたクライアントのアバターURLを返す
	// 問題が発生した場合はエラーを返す
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
