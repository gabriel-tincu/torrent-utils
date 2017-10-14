package bencode

const (
	magnetPrefix        = "magnet:?"
	trackerKey          = "tr"
	topicKey            = "xt"
	displayNameKey      = "dn"
	lengthKey           = "xl"
	acceptableSourceKey = "as"
	exactSourceKey      = "xs"
	keywordTopicKey     = "kt"
	manifestTopiKey     = "mt"
)

type Trackers []string
type Topic string
type DisplayName string
