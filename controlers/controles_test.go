package controlers_test

import (
	"fmt"
	"testing"

	"github.com/sanzhang007/crawl_nodes/controlers"
)

func Test(t *testing.T) {
	b, err := controlers.GetPotocolServerAndPortAndPassword(`ss://YWVzLTI1Ni1jZmI6cXdlclJFV1FAQEAyMjEuMTUwLjEwOS41Ojk1NTU=#%F0%9F%87%B0%F0%9F%87%B7KR_439
	`)
	fmt.Println(string(b))
	fmt.Println(err)

}
