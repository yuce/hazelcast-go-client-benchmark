/*
 * Copyright (c) 2008-2022, Hazelcast, Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func bye(msg string) {
	// ignoring the error
	_, _ = fmt.Fprintf(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		bye("Usage: parbench config.json")
	}
	config, err := LoadConfigFromPath(os.Args[1])
	if err != nil {
		panic(err)
	}
	svc, err := StartNewService(context.Background(), config)
	if err != nil {
		panic(err)
	}
	took := measureTime(func() {
		for i := 0; i < config.OperationCount; i++ {
			key := fmt.Sprintf("key-%d")
			if err := svc.Do(context.Background(), key); err != nil {
				panic(err)
			}
		}

	})
	// ignoring the error
	_, _ = fmt.Fprintf(os.Stderr, "Took: %d ms for %d operations\n", took.Milliseconds(), config.OperationCount)
}

func measureTime(f func()) time.Duration {
	tic := time.Now()
	f()
	toc := time.Now()
	return toc.Sub(tic)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustValue(v interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return v
}
