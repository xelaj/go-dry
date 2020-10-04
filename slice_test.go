package dry

import (
	"testing"

	"github.com/k0kubun/pp/v3"
)

func TestDeleteIndex(t *testing.T) {
	res := DeleteIndex([]string{"1", "2", "3"}, 2).([]string)
	pp.Println(res)

	pp.Println(SliceExpand([]string{"a", "b", "c"}, 0, 6))
}
