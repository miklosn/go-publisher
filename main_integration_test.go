package publisher

import (
	"context"
	"log"
	"testing"
)

var p *Publisher

func init() {
	p, _ = New()
}
func BenchmarkSend(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := p.Send(context.Background(), "learn-242515", "asd", []byte{0, 0})
		if err != nil {
			log.Println(err)
			b.Fail()
		}
	}
}
