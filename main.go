package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/prometheus/common/log"
	"github.com/septemhill/translator/repository/boltdb"
	"github.com/septemhill/translator/service"
	"github.com/septemhill/translator/service/azure"
)

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
	if len(os.Args) < 2 {
		log.Errorln("No enough parameter")
	}

	ctx := context.Background()
	repo, _ := boltdb.NewBoltDBRepository()
	defer repo.Close(ctx)

	proxy := NewDictionaryProxy(NewDictionary(azure.NewWordTranslateService()), repo)

	st := time.Now()
	tr, err := proxy.Lookup(ctx, os.Args[1], "en", "zh-Hans", true)
	et := time.Now()

	if err != nil {
		log.Fatal("Failed to lookup word: ", err)
	}
	streamOut(tr, os.Stdout)
	fmt.Println(et.Sub(st))
}
