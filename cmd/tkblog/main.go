// Copyright 2024 Innkeeper tkane(童昌泰) &lt;452212568@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/tkane/tkblog.

package main

import (
	"os"

	"github.com/tkane/tkblog/internal/tkblog"
)

func main() {
	command := tkblog.NewBlogCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
