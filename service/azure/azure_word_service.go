package azure

import (
	"context"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/translatortext"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/septemhill/translator/service"
)

var (
	key      = os.Getenv("AZURE_COGNITIVE_SERVICE_KEY")
	endpoint = os.Getenv("AZURE_COGNITIVE_SERVICE_ENDPOINT")
)

type wordService struct {
	client *translatortext.TranslatorClient
}

func NewWordTranslateService() *wordService {
	client := translatortext.NewTranslatorClient(endpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(key)

	return &wordService{
		client: &client,
	}
}

func (s *wordService) lookup(ctx context.Context, words []string, src, dst string) (*translatortext.ListDictionaryLookupResultItem, error) {
	var texts []translatortext.DictionaryLookupTextInput

	for _, word := range words {
		texts = append(texts, translatortext.DictionaryLookupTextInput{Text: to.StringPtr(word)})
	}

	results, err := s.client.DictionaryLookup(ctx, src, dst, texts, "")
	if err != nil {
		return nil, err
	}

	return &results, nil
}

func (s *wordService) queryExample(ctx context.Context, word, translation, src, dst string) ([]string, error) {
	examples := []translatortext.DictionaryExampleTextInput{
		{Text: to.StringPtr(word), Translation: to.StringPtr(translation)},
	}

	exresults, err := s.client.DictionaryExamples(ctx, src, dst, examples, "")
	if err != nil {
		return nil, err
	}

	var results []string
	for _, v := range *exresults.Value {
		for _, trans := range *v.Examples {
			results = append(results, strings.Join([]string{*trans.SourcePrefix, *trans.SourceTerm, *trans.SourceSuffix}, ""))
		}
	}

	return results, nil
}

func (s *wordService) WordLookup(ctx context.Context, word, src, dst string, ex bool) (*service.Translation, error) {
	result, err := s.lookup(ctx, []string{word}, src, dst)
	if err != nil {
		return nil, err
	}

	m := make(map[string]*service.Speech)
	for _, v := range *result.Value {
		for _, trans := range *v.Translations {
			if m[*trans.PosTag] == nil {
				m[*trans.PosTag] = new(service.Speech)
			}
			m[*trans.PosTag].Result = append(m[*trans.PosTag].Result, *trans.DisplayTarget)
		}
	}

	if ex {
		for _, sp := range m {
			if len(sp.Result) > 0 {
				res, err := s.queryExample(ctx, word, sp.Result[0], src, dst)
				if err != nil {
					return nil, err
				}
				sp.Example = res
			}
		}
	}

	var trans service.Translation
	trans.Word = word
	for k, v := range m {
		v.Name = k
		trans.Speech = append(trans.Speech, *v)
	}

	return &trans, nil
}

// func (s *wordService) WordsLookup(ctx context.Context, words []string, src, dst string) error {
// 	result, err := s.lookup(ctx, words, src, dst)
// 	if err != nil {
// 		return err
// 	}

// 	//TODO: parse `results` and return expected data
// 	_ = result

// 	fmt.Println(result)
// 	return nil
// }

var _ service.WordTranslateService = (*wordService)(nil)
