/**
 *
 * 文件功能介绍
 *
 * @FileName: sha512.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-11 10:44:32
 * @UpdateTime: 2018-07-11 10:44:54
 * Copyright@2018 版权所有
 */

package utils

import (
	"crypto/sha512"
	"hash"
)

type Sha512 struct {
	h hash.Hash
}

func NewSha512() *Sha512 {
	s := new(Sha512)
	s.h = sha512.New()
	return s
}

func (s *Sha512) Finish() []byte {
	return s.h.Sum(nil)
}

func (s *Sha512) Add(bytes []byte) (int, error) {
	return s.h.Write(bytes)
}

func (s *Sha512) Add32(i uint32) (int, error) {
	var b []byte
	b = append(b, byte(((i >> 24) & 0xFF)))
	b = append(b, byte(((i >> 16) & 0xFF)))
	b = append(b, byte(((i >> 8) & 0xFF)))
	b = append(b, byte((i & 0xFF)))
	return s.h.Write(b)
}

func (s *Sha512) Finish256() []byte {
	return s.h.Sum(nil)[0:32]
}

func (s *Sha512) Finish128() []byte {
	return s.h.Sum(nil)[0:16]
}
