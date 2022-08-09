package main

import "time"

func measureTime(f func()) time.Duration {
	tic := time.Now()
	f()
	toc := time.Now()
	return toc.Sub(tic)
}
