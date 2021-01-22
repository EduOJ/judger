package main

import (
	"fmt"
	"github.com/suntt2019/Judger/judger"
)

func main() {

	config := judger.Config{
		MaxCPUTime:           1000,
		MaxRealTine:          2000,
		MaxMemory:            128 * 1024 * 1024,
		MaxStack:             32 * 1024 * 1024,
		MaxProcessNumber:     200,
		MaxOutputSize:        10000,
		MemoryLimitCheckOnly: 0,
		ExePath:              "test_programs/a+b/a+b",
		InputPath:            "test_programs/a+b/a+b.in",
		OutputPath:           "test_programs/a+b/a+b.out",
		ErrorPath:            "test_programs/a+b/a+b.out",
		Args:                 []string{},
		Env:                  []string{},
		LongPath:             "test.log",
		SeccompRuleName:      "c_cpp",
		Uid:                  0,
		Gid:                  0,
	}
	fmt.Print(judger.Run(config))
}
