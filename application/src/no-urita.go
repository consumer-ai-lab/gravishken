//go:build nowebview

package main

func uritaOpenWv(url string) {
	panic("App is not built with webview support!")
}
