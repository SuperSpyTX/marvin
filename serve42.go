package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

var killText = []byte(`#!/bin/sh

[ -d .git ] || git rev-parse --git-dir > /dev/null 2>&1 && (
	git commit --allow-empty -m 'Ran a curl | sh command'
)

SHELL_NAME=${0:-sh}
echo "Don't pipe curl into $SHELL_NAME. Someone could run naughty commands."
echo "Always save the script to a file and inspect it before running."

osascript -e 'say "Don\'t pipe curl into sh"'
exit 3

`)

func curlKiller(wrap http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := httptest.NewRecorder()
		wrap.ServeHTTP(rec, r)

		rec.HeaderMap.Del("Content-Length")
		for k, v := range rec.HeaderMap {
			w.Header()[k] = v
		}
		r.Header.Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(rec.Code)

		if strings.Contains(r.Header.Get("User-Agent"), "curl") {
			w.Write(killText)
		}

		rec.Body.WriteTo(w)
	})
}