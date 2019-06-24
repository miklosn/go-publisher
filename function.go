package function

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/miklosn/pubsub-publisher/pkg/publisher"
)

var p *publisher.Publisher

func PublishMessage(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	p.Send("learn-242515", "asd", []byte{0, 0, 0, 0, 0})
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
	start = time.Now()
	p.Send("learn-242515", "asd", []byte{0, 0, 0, 0, 0})
	elapsed = time.Since(start)
	log.Printf("Binomial took %s", elapsed)
}

func init() {
	p, _ = publisher.New()
	for _, pair := range os.Environ() {
		fmt.Println(pair)
	}

}
