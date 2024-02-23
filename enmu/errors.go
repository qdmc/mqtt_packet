package enmu

import "errors"

var TypeError = errors.New("message type error")
var TopicsEmpty = errors.New("message topics is empty")
var TopicError = errors.New("topic name is error")
var FixedEmpty = errors.New("message fixed header is empty")
var RemainingLengthErr = errors.New("FixedHeader.RemainingLength is error")
