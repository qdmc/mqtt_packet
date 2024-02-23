package enmu

import "errors"

var TypeError = errors.New("message type error")
var TopicsEmpty = errors.New("message topics is empty")
var TopicError = errors.New("topic name is error")
var FixedEmpty = errors.New("message fixed header is empty")
var RemainingLengthErr = errors.New("FixedHeader.RemainingLength is error")
var QosError = errors.New("message qos must be 0,1 or 2")
