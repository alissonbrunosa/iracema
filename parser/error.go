package parser

import "fmt"

type Error struct {
	Msg string
}

func (e Error) Error() string {
	return e.Msg
}

type ErrorList []*Error

func (list ErrorList) Err() error {
	if len(list) == 0 {
		return nil
	}

	return list
}

func (list ErrorList) Error() string {
	switch len(list) {
	case 0:
		return "no errors"
	case 1:
		return list[0].Error()
	default:
		for _, e := range list {
			fmt.Println(e)
		}
		return fmt.Sprintf("%s (and %d more errors)", list[0], len(list)-1)
	}
}
