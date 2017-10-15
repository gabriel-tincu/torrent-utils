package bcode

import (
	"fmt"
	"strconv"
)

const (
	intStartChar  = 'i'
	endChar       = 'e'
	listStartChar = 'l'
	dictStartChar = 'd'
)

var digits = map[byte]bool{
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
	'0': true,
	':': true,
}

type Bencoded map[string]interface{}

func decodeList(buff []byte) (result []interface{}, remaining []byte, err error) {
	if rune(buff[0]) != listStartChar {
		return nil, nil, fmt.Errorf("start character for list should be %c", listStartChar)
	}
	buff = buff[1:]
	next := buff[0]
	for {
		if buff[0] == endChar {
			return result, buff[1:], nil
		}
		switch next {
		case intStartChar:
			resultInt, remaining, err := decodeInt(buff)
			if err != nil {
				return nil, nil, err
			}
			buff = remaining
			result = append(result, resultInt)
		case listStartChar:
			resultList, remaining, err := decodeList(buff)
			if err != nil {
				return nil, nil, err
			}
			buff = remaining
			result = append(result, resultList)
		case dictStartChar:
			resultDict, remaining, err := decodeDict(buff)
			if err != nil {
				return nil, nil, err
			}
			buff = remaining
			result = append(result, resultDict)
		case endChar:
			return result, buff[1:], nil
		default:
			if !digits[next] {
				return nil, nil, fmt.Errorf("should receive an integer")
			}
			resultBytes, remaining, err := decodeByte(buff)
			if err != nil {
				return nil, nil, err
			}
			buff = remaining
			result = append(result, resultBytes)
		}
	}
}

func decodeDict(buff []byte) (result Bencoded, remaining []byte, err error) {
	if buff[0] != dictStartChar {
		return nil, nil, fmt.Errorf("start character for dict should be %c", dictStartChar)
	}
	result = make(map[string]interface{})
	var key string
	var list []interface{}
	var dic Bencoded
	var integer int
	var bits string
	buff = buff[1:]
	for {
		if buff[0] == endChar {
			return result, buff[1:], nil
		}
		key, buff, err = decodeByte(buff)
		if err != nil {
			return nil, nil, err
		}
		next := buff[0]
		switch next {
		case listStartChar:
			list, buff, err = decodeList(buff)
			if err != nil {
				return nil, nil, err
			}
			result[key] = list
			continue
		case intStartChar:
			integer, buff, err = decodeInt(buff)
			if err != nil {
				return nil, nil, err
			}
			result[key] = integer
			continue
		case dictStartChar:
			dic, buff, err = decodeDict(buff)
			if err != nil {
				return nil, nil, err
			}
			result[key] = dic
			continue
		case endChar:
			return result, buff[1:], nil
		default:
			if !digits[next] {
				return nil, nil, fmt.Errorf("should receive an integer")
			}
			bits, buff, err = decodeByte(buff)
			if err != nil {
				return nil, nil, err
			}
			result[key] = bits
			continue
		}
	}
}

func decodeByte(buff []byte) (result string, remaining []byte, err error) {
	size := []byte{}
	index := 0
	if !digits[buff[index]] {
		return "", nil, fmt.Errorf("decoded char should be a number or colon :%c", buff[index])
	}

	for {
		if buff[index] != ':' {
			size = append(size, buff[index])
			index++
		} else {
			index++
			break
		}
	}
	sizeStr := string(size)
	byteSize, err := strconv.Atoi(sizeStr)
	if err != nil {
		return "", nil, err
	}
	return string(buff[index : index+byteSize]), buff[index+byteSize:], nil
}

func decodeInt(buff []byte) (result int, remaining []byte, err error) {
	chars := []byte{}
	if buff[0] != intStartChar {
		return 0, nil, fmt.Errorf("input buffer does not contain an integer signaling char")
	}
	index := 1
	for {
		bit := buff[index]
		if bit == endChar {
			index++
			break
		}
		chars = append(chars, buff[index])
		index++
	}
	result, err = strconv.Atoi(string(chars))
	return result, buff[index:], err
}

func BDecode(bits []byte) (Bencoded, error) {
	result, remaining, err := decodeDict(bits)
	if err != nil {
		return nil, err
	}
	if len(remaining) != 0 {
		return nil, fmt.Errorf("remaining buffer should be empty")
	}
	return result, nil
}
