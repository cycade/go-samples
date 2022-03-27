package main

import (
	"encoding/binary"
	"fmt"
)

// FROM https://github.com/Terry-Mao/goim/blob/master/api/protocol/protocol.go
const (
	// size
	_packSize      = 4 // 包长度
	_headerSize    = 2 // 包头长度
	_verSize       = 2 // 协议版本号
	_opSize        = 4 // 操作码
	_seqSize       = 4 // jeager-id ???
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = int32(1<<12) + int32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_heartOffset  = _seqOffset + _seqSize
)

type DecodeResult struct {
	PacketSize uint32
	HeaderSize uint16
	Version    uint16
	Op         uint32
	SequenceId uint32
	Body       string
	Finalized  bool
}

func Decoder(raw []byte) (result *DecodeResult, err error) {
	result.PacketSize = binary.BigEndian.Uint32(raw[:_headerOffset])
	result.HeaderSize = binary.BigEndian.Uint16(raw[_headerOffset:_verOffset])
	result.Version = binary.BigEndian.Uint16(raw[_verOffset:_opOffset])
	result.SequenceId = binary.BigEndian.Uint32(raw[_opOffset:_seqOffset])

	if int32(result.PacketSize) > _maxPackSize {
		return nil, fmt.Errorf("message 长度超过最大限制 %d", _maxPackSize)
	}

	if result.HeaderSize != _rawHeaderSize {
		return nil, fmt.Errorf("实际 Header 长度与元信息不符")
	}

	result.Body = string(rune(binary.BigEndian.Uint32(raw[_rawHeaderSize:])))
	result.Finalized = true
	if int(result.PacketSize) > len(raw) {
		result.Finalized = false
	}

	return result, nil
}
