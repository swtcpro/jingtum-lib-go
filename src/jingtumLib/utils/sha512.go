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

//Sha512 Sha512
type Sha512 struct {
	h hash.Hash
}

//NewSha512 NewSha512
func NewSha512() *Sha512 {
	s := new(Sha512)
	s.h = sha512.New()
	return s
}

//Finish Finish
func (s *Sha512) Finish() []byte {
	return s.h.Sum(nil)
}

//Add Add
func (s *Sha512) Add(bytes []byte) (int, error) {
	return s.h.Write(bytes)
}

//Add32 Add32
func (s *Sha512) Add32(i uint32) (int, error) {
	var b []byte
	b = append(b, byte(((i >> 24) & 0xFF)))
	b = append(b, byte(((i >> 16) & 0xFF)))
	b = append(b, byte(((i >> 8) & 0xFF)))
	b = append(b, byte((i & 0xFF)))
	return s.h.Write(b)
}

//Finish256 32字节
func (s *Sha512) Finish256() []byte {
	return s.h.Sum(nil)[0:32]
}

//Finish128 16字节
func (s *Sha512) Finish128() []byte {
	return s.h.Sum(nil)[0:16]
}
