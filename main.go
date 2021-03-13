package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/septemhill/translator/logger"
	"github.com/septemhill/translator/logger/stdoutlog"
	"github.com/septemhill/translator/repository/boltdb"
	"github.com/septemhill/translator/service"
	"github.com/septemhill/translator/service/azure"
)

var fromLang = flag.String("from", "en", "Language support list: https://docs.microsoft.com/en-gb/azure/cognitive-services/translator/language-support")
var toLang = flag.String("to", "zh-Hans", "Language support list: https://docs.microsoft.com/en-gb/azure/cognitive-services/translator/language-support")

func streamOut(tr *service.Translation, w io.Writer) {
	for _, s := range tr.Speech {
		w.Write([]byte(fmt.Sprintf("%s\n", s.Name)))
		for _, wd := range s.Result {
			w.Write([]byte(fmt.Sprintf("\t%s\n", wd)))
		}
		w.Write([]byte("\n"))
		for _, ex := range s.Example {
			w.Write([]byte(fmt.Sprintf("\t%s\n", ex)))
		}
	}
}

func main() {
	ctx := context.Background()
	ctx = logger.NewContext(ctx, stdoutlog.StdoutLogger())

	if len(os.Args) < 2 {
		logger.ContextFatalln(ctx, "No enough parameter")
	}
	flag.Parse()

	repo, err := boltdb.NewBoltDBRepository()
	if err != nil {
		logger.ContextErrorln(ctx, "failed to create boltdb repository")
	}
	defer repo.Close(ctx)

	proxy := NewDictionaryProxy(NewDictionary(azure.NewWordTranslateService()), repo)

	l := len(os.Args)
	tr, err := proxy.Lookup(ctx, os.Args[l-1], *fromLang, *toLang, true)

	if err != nil {
		logger.ContextFatalf(ctx, "failed to lookup word: %s\n", err.Error())
	}
	streamOut(tr, os.Stdout)
}
