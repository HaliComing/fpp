package sudaili

import (
	"fmt"
	"testing"
)

// 测试提取
func TestExtract(t *testing.T) {
	extract, err := New().Extract()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, ip := range extract {
			fmt.Printf("%s:%s\n", ip.IP, ip.Port)
		}
	}
}
