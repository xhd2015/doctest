package oshim

import "example.com/realmod/http/internal/leaf"

func Bridge() string { return leaf.Hello() }
