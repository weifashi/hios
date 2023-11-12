package error

import "errors"

type WithError struct {
	Msg    string
	Detail interface{}
	Map    map[string]interface{}
	Err    error
}

func (e WithError) Error() string {
	content := ""
	if e.Detail != nil {
		content = e.Msg
	} else if e.Map != nil {
		content = e.Msg
	} else {
		content = e.Msg
	}
	if content == "" {
		if e.Err != nil {
			return e.Err.Error()
		}
		return errors.New(e.Msg).Error()
	}
	return content
}

func New(Key string) WithError {
	return WithError{
		Msg:    Key,
		Detail: nil,
		Err:    nil,
	}
}

func WithDetail(Key string, detail any, err error) WithError {
	return WithError{
		Msg:    Key,
		Detail: detail,
		Err:    err,
	}
}

func WithErr(Key string, err error) WithError {
	return WithError{
		Msg:    Key,
		Detail: "",
		Err:    err,
	}
}

func WithMap(Key string, maps map[string]any, err error) WithError {
	return WithError{
		Msg: Key,
		Map: maps,
		Err: err,
	}
}
