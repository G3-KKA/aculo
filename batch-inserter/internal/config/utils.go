package config

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

func Print(cfg Config) {
	ybytes, err := yaml.Marshal(cfg)
	if err == nil {
		bbuf := bytes.Buffer{}
		bbuf.WriteString("[INFO] Config: \n\t")
		for _, v := range ybytes {
			if v == '\n' {
				bbuf.WriteByte('\n')
				bbuf.WriteByte('\t')
				continue
			}
			bbuf.WriteByte(v)

		}
		fmt.Println(bbuf.String())
	}
}
