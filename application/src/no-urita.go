//go:build nowebview

package main

func openWv(url string) {
	panic("App is not built with webview support!")
}
