package magnet

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	magnetPrefix        = "magnet:?"
	trackerKey          = "tr"
	topicKey            = "xt"
	displayNameKey      = "dn"
	lengthKey           = "xl"
	acceptableSourceKey = "as"
	exactSourceKey      = "xs"
	keywordTopicKey     = "kt"
	manifestTopicKey    = "mt"
)

type Trackers []string

type Magnet struct {
	Trackers         Trackers
	Length           int
	DisplayName      string
	Topic            string
	AcceptableSource string
	ExactSource      string
	KeywordTopic     string
	ManifestTopic    string
}

func ParseMagnet(content []byte) (*Magnet, error) {
	start := string(content[:len(magnetPrefix)])
	if start != magnetPrefix {
		return nil, fmt.Errorf("error parsing magnet link, the string does not start properly :(")
	}
	val, err := url.ParseQuery(string(content[len(magnetPrefix):]))
	if err != nil {
		return nil, err
	}
	response := &Magnet{}
	for k, v := range val {
		switch k {
		case trackerKey:
			response.Trackers = v
		case displayNameKey:
			response.DisplayName = v[0]
		case lengthKey:
			val, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, err
			}
			response.Length = val
		case acceptableSourceKey:
			response.AcceptableSource = v[0]
		case topicKey:
			response.Topic = v[0]
		case exactSourceKey:
			response.ExactSource = v[0]
		case manifestTopicKey:
			response.ManifestTopic = v[0]
		case keywordTopicKey:
			response.KeywordTopic = v[0]
		}
	}
	return response, nil
}
