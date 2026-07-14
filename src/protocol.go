package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

// --- Функции для работы с VarInt протокола Minecraft ---
func ReadVarInt(r io.Reader) (int, error) {
	var num int
	var shift uint
	buf := make([]byte, 1)
	for {
		_, err := r.Read(buf)
		if err != nil {
			return 0, err
		}
		b := buf[0]
		num |= int(b&0x7F) << shift
		if (b & 0x80) == 0 {
			break
		}
		shift += 7
		if shift >= 32 {
			return 0, errors.New("VarInt is too big")
		}
	}
	return num, nil
}

func WriteVarInt(w io.Writer, val int) error {
	for {
		if (val & ^0x7F) == 0 {
			_, err := w.Write([]byte{byte(val)})
			return err
		}
		_, err := w.Write([]byte{byte((val & 0x7F) | 0x80)})
		if err != nil {
			return err
		}
		val = int(uint(val) >> 7)
	}
}

func WriteString(w io.Writer, s string) error {
	if err := WriteVarInt(w, len(s)); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}

// PingServer отправляет Handshake и Status Request, возвращая распарсенный JSON
func PingServer(address string, timeout time.Duration) (*StatusResponse, error) {
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		host = address
		portStr = "25565"
		address = net.JoinHostPort(host, portStr)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("неверный порт: %w", err)
	}

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться: %w", err)
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(timeout))

	// 1. Handshake Packet
	hsBuf := new(bytes.Buffer)
	_ = WriteVarInt(hsBuf, 0x00)
	_ = WriteVarInt(hsBuf, 47) // Протокол 47 подходит для большинства SLP опросов
	_ = WriteString(hsBuf, host)
	_, _ = hsBuf.Write([]byte{byte(port >> 8), byte(port)})
	_ = WriteVarInt(hsBuf, 1) // Next State = 1 (Status)

	if err := WriteVarInt(conn, hsBuf.Len()); err != nil {
		return nil, err
	}
	if _, err := conn.Write(hsBuf.Bytes()); err != nil {
		return nil, err
	}

	// 2. Status Request Packet
	reqBuf := new(bytes.Buffer)
	_ = WriteVarInt(reqBuf, 0x00)

	if err := WriteVarInt(conn, reqBuf.Len()); err != nil {
		return nil, err
	}
	if _, err := conn.Write(reqBuf.Bytes()); err != nil {
		return nil, err
	}

	// 3. Status Response Packet
	packetLength, err := ReadVarInt(conn)
	if err != nil {
		return nil, err
	}
	if packetLength <= 0 {
		return nil, errors.New("неверная длина пакета")
	}

	packetID, err := ReadVarInt(conn)
	if err != nil {
		return nil, err
	}
	if packetID != 0x00 {
		return nil, fmt.Errorf("неожиданный ID пакета: 0x%02x", packetID)
	}

	jsonLength, err := ReadVarInt(conn)
	if err != nil {
		return nil, err
	}

	jsonBytes := make([]byte, jsonLength)
	_, err = io.ReadFull(conn, jsonBytes)
	if err != nil {
		return nil, err
	}

	var response StatusResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
