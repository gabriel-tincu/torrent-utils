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
	colonChar     = ':'
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
	var bits []byte
	buff = buff[1:]
	for {
		if buff[0] == endChar {
			return result, buff[1:], nil
		}
		key, buff, err = decodeString(buff)
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

func decodeString(buff []byte) (result string, remaining []byte, err error) {
	size := []byte{}
	index := 0
	if !digits[buff[index]] {
		return "", nil, fmt.Errorf("decoded char should be a number or colon :%c", buff[index])
	}

	for {
		if buff[index] != colonChar {
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


func decodeByte(buff []byte) (result []byte, remaining []byte, err error) {
	size := []byte{}
	index := 0
	if !digits[buff[index]] {
		return nil, nil, fmt.Errorf("decoded char should be a number or colon :%c", buff[index])
	}

	for {
		if buff[index] != colonChar {
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
		return nil, nil, err
	}
	return buff[index : index+byteSize], buff[index+byteSize:], nil
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

func Marshal(data Bencoded) ([]byte, error) {
	result := []byte{}
	return encodeDict(data, result)
}

func encodeDict(data Bencoded, bits []byte) ([]byte, error) {
	var err error
	bits = append(bits, dictStartChar)
	for key, val := range data {
		bits, err = encodeBytes([]byte(key), bits)
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case int:
			bits, err = encodeInt(v, bits)
			if err != nil {
				return nil, err
			}
		case string:
			bits, err = encodeBytes([]byte(v), bits)
			if err != nil {
				return nil, err
			}
		case []byte:
			bits, err = encodeBytes(v, bits)
			if err != nil {
				return nil, err
			}
		case []interface{}:
			bits, err = encodeList(v, bits)
			if err != nil {
				return nil, err
			}
		case map[string]interface{}:
			bits, err = encodeDict(v, bits)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("data in a dict should be one of <<int>> <<string>> <<list>> <<dict>>:\n%T", v)
		}
	}
	bits = append(bits, endChar)
	return bits, nil

}

func encodeList(data []interface{}, bits []byte) ([]byte, error) {
	var err error
	bits = append(bits, listStartChar)
	for _, val := range data {
		switch v := val.(type) {
		case int:
			bits, err = encodeInt(v, bits)
			if err != nil {
				return nil, err
			}
		case string:
			bits, err = encodeBytes([]byte(v), bits)
			if err != nil {
				return nil, err
			}
		case []byte:
			bits, err = encodeBytes(v, bits)
			if err != nil {
				return nil, err
			}
		case []interface{}:
			bits, err = encodeList(v, bits)
			if err != nil {
				return nil, err
			}
		case map[string]interface{}:
			bits, err = encodeDict(v, bits)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("data in a list should be one of <<int>> <<string>> <<list>> <<dict>> :\n%T", v)
		}
	}
	bits = append(bits, endChar)
	return bits, nil
}

func encodeInt(data int, bits []byte) ([]byte, error) {
	bits = append(bits, intStartChar)
	bits = append(bits, []byte(strconv.Itoa(data))...)
	bits = append(bits, endChar)
	return bits, nil
}

func encodeString(data string, bits []byte) (string, error) {
	sizeStr := strconv.Itoa(len(data))
	bits = append(bits, []byte(sizeStr)...)
	bits = append(bits, colonChar)
	bits = append(bits, []byte(data)...)
	return string(bits), nil
}

func encodeBytes(data []byte, bits []byte) ([]byte, error) {
	sizeStr := strconv.Itoa(len(data))
	bits = append(bits, []byte(sizeStr)...)
	bits = append(bits, colonChar)
	bits = append(bits, data...)
	return bits, nil
}

func Unmarshal(bits []byte) (Bencoded, error) {
	result, remaining, err := decodeDict(bits)
	if err != nil {
		return nil, err
	}
	if len(remaining) != 0 {
		return nil, fmt.Errorf("remaining buffer should be empty")
	}
	return result, nil
}
