package bencode

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

type Bencoded map[string]interface{}

func decodeList(buff []byte) (result []interface{}, remaining []byte, err error) {
	if rune(buff[0]) != listStartChar {
		return nil, nil, fmt.Errorf("start character for list should be %c", listStartChar)
	}
	next := rune(buff[1])
	buff = buff[1:]
	for {
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
	if rune(buff[0]) != dictStartChar {
		return nil, nil, fmt.Errorf("start character for list should be %c", dictStartChar)
	}
	var key string
	var list []interface{}
	var dic Bencoded
	var integer int
	var bits string
	buff = buff[1:]
	for {
		key, buff, err = decodeByte(buff)
		if err != nil {
			return nil, nil, err
		}
		next, buff := buff[0], buff[1:]
		switch next {
		case listStartChar:
			list, buff, err = decodeList(buff)
			if err != nil {
				return nil, nil, err
			}
			result[key] = list
		case intStartChar:
			integer, buff, err = decodeInt(buff)
			if err != nil {
				return nil, nil, err
			}
			result[key] = integer
		case dictStartChar:
			dic, buff, err = decodeDict(buff)
			if err != nil {
				return nil, nil, err
			}
			result[key] = dic
		case endChar:
			return result, buff[1:], nil
		default:
			bits, buff, err = decodeByte(buff)
			if err != nil {
				return nil, nil, err
			}
			result[key] = bits
		}
	}
}

func decodeByte(buff []byte) (result string, remaining []byte, err error) {
	size := []byte{}
	index := 0
	for {
		if rune(result[index]) != ':' {
			size = append(size, result[index])
			index++
		} else {
			index++
			break
		}
	}
	byteSize, err := strconv.Atoi(string(size))
	if err != nil {
		return "", nil, err
	}
	return string(buff[index : index+byteSize]), buff[:index+byteSize], nil
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
		index++
		chars = append(chars, buff[index])
	}
	result, err = strconv.Atoi(string(chars))
	return result, buff[index:], err
}

func BDecode(reader io.Reader) (*Bencoded, error) {
	bits, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	result, remaining, err := decodeDict(bits)
	if err != nil {
		return nil, err
	}
	if len(remaining) != 0 {
		return nil, fmt.Errorf("remaining buffer should be empty")
	}
	return &result, nil
}
