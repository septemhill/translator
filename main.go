package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/translatortext"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

var cfg *TslConfig

func init() {
	var err error
	cfg, err = NewTslConfig()
	if err != nil {
		panic("tslconfig init failed")
	}
}

type DictionaryClient struct {
	client translatortext.TranslatorClient
}

func NewDictionaryClient() *DictionaryClient {
	key := os.Getenv("AZURE_TRANSLATOR_TEXT_KEY")
	endpoint := os.Getenv("AZURE_TRANSLATOR_TEXT_ENDPOINT")

	client := translatortext.NewTranslatorClient(endpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(key)

	return &DictionaryClient{
		client: client,
	}
}

func (cli *DictionaryClient) displayDictionaryLookupResult(dres *translatortext.DictionaryLookupResultItem, ex func(translatortext.DictionaryLookupResultItem) *translatortext.DictionaryExampleResultItem) {
	var eres *translatortext.DictionaryExampleResultItem
	if dres == nil {
		return
	}

	if ex != nil {
		eres = ex(*dres)
	}

	for _, t := range *dres.Translations {
		fmt.Println(to.String(t.PosTag))
		fmt.Printf("\t%s\n", to.String(t.DisplayTarget))

		if ex != nil {
			for _, e := range *eres.Examples {
				fmt.Printf("\t\t%s %s %s\n", to.String(e.SourcePrefix), to.String(e.SourceTerm), to.String(e.SourceSuffix))
			}
		}
	}
}

func (cli *DictionaryClient) dictionaryLookup(word string) *translatortext.DictionaryLookupResultItem {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	lookupInput := []translatortext.DictionaryLookupTextInput{
		translatortext.DictionaryLookupTextInput{
			Text: to.StringPtr(word),
		},
	}

	result, err := cli.client.DictionaryLookup(ctx, cfg.DefaultFrom, cfg.DefaultTo, lookupInput, "")
	if err != nil {
		log.Println("translator: ", err)
	}

	if len(*result.Value) == 0 {
		return nil
	}

	return &(*result.Value)[0]
}

func (cli *DictionaryClient) dictionaryExample(dlresult translatortext.DictionaryLookupResultItem) *translatortext.DictionaryExampleResultItem {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	inputs := make([]translatortext.DictionaryExampleTextInput, 0)

	for _, word := range *dlresult.Translations {
		i := translatortext.DictionaryExampleTextInput{
			Text:        (*word.BackTranslations)[0].NormalizedText,
			Translation: word.NormalizedTarget,
		}
		inputs = append(inputs, i)
	}

	result, err := cli.client.DictionaryExamples(ctx, cfg.DefaultFrom, cfg.DefaultTo, inputs, "")
	if err != nil {
		log.Println("DictionaryExample failed: ", err)
	}

	return &(*result.Value)[0]
}

func (cli *DictionaryClient) DictionaryLookup(word string, ex bool) {
	if ex {
		cli.displayDictionaryLookupResult(cli.dictionaryLookup(word), cli.dictionaryExample)
		return
	}

	cli.displayDictionaryLookupResult(cli.dictionaryLookup(word), nil)
}

func main() {
	client := NewDictionaryClient()
	if client == nil {
		log.Printf("internal unknown error")
	}

	client.DictionaryLookup("document", false)
}
