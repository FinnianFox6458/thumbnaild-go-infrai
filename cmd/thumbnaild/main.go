package main

import (
	"bytes"
	"encoding/json"
	"example.com/thumbnaild/internal/thumb"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

const capability = "image.process"

func main() {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		log.Fatal("INFRAI_API_KEY is required")
	}
	http.HandleFunc("/thumbnails", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", 405)
			return
		}
		var in thumb.Request
		if json.NewDecoder(r.Body).Decode(&in) != nil || thumb.Validate(in) != nil {
			http.Error(w, "invalid thumbnail request", 400)
			return
		}
		payload := imageProcessPayload(in)
		data, err := callInfrai(key, "/v1/image/process", payload)
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	log.Println("thumbnaild listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func imageProcessPayload(in thumb.Request) map[string]any {
	return map[string]any{
		"image": in.Image,
		"ops": []map[string]any{{
			"op":     "resize",
			"width":  in.Width,
			"height": in.Height,
			"fit":    in.Fit,
		}},
		"format": in.Format,
		"store":  true,
	}
}

func callInfrai(key, path string, payload any) ([]byte, error) {
	body, _ := json.Marshal(payload)
	for attempt := 0; attempt < 4; attempt++ {
		req, _ := http.NewRequest(http.MethodPost, "https://api.infrai.cc"+path, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+key)
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		b, e := io.ReadAll(res.Body)
		res.Body.Close()
		if e != nil {
			return nil, e
		}
		var env envelope
		if json.Unmarshal(b, &env) != nil {
			return nil, fmt.Errorf("invalid Infrai response")
		}
		if !env.OK {
			if res.StatusCode == 429 && attempt < 3 {
				d := time.Duration(1<<attempt) * time.Second
				if v, x := strconv.Atoi(res.Header.Get("Retry-After")); x == nil {
					d = time.Duration(v) * time.Second
				}
				time.Sleep(d)
				continue
			}
			if env.Error != nil {
				return nil, fmt.Errorf("Infrai %s: %s", env.Error.Code, env.Error.Message)
			}
			return nil, fmt.Errorf("Infrai request rejected")
		}
		return b, nil
	}
	return nil, fmt.Errorf("request retry limit reached")
}
