package service

import "context"

type Speech struct {
	Name    string
	Result  []string
	Example []string
}

type Translation struct {
	Word   string
	Speech []Speech
}

type WordTranslateService interface {
	WordLookup(ctx context.Context, word, src, dst string, example bool) (*Translation, error)
	//WordsLookup(ctx context.Context, word []string, src, dst string) error
}
