package main

import (
	"fmt"
	"os"

	"github.com/zapomnij/firecraft/pkg/auth"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	ms, err := auth.NewMsAuth(os.Args[1])
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	mc, err := auth.NewMinecraftAuthentication(ms.AccessToken)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Println(mc.MinecraftToken, mc.Userhash)
}
