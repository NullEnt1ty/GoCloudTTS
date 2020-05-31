package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] text\n\nOPTIONS:\n\n", os.Args[0])
	flag.PrintDefaults()
}

func ensureCacheDirExists(cacheDir string) {
	info, err := os.Stat(cacheDir)

	if os.IsNotExist(err) {
		err := os.MkdirAll(cacheDir, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}

		return
	}

	if !info.IsDir() {
		log.Fatalf("Cache dir at '%s' exists but is not a directory", cacheDir)
	}
}

func getCacheFilename(text string, language string, voiceName string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	textHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return fmt.Sprintf("%s-%s-%s.wav", textHash, language, voiceName)
}

func findCachedAudioData(cacheDir string, text string, languageCode string, voiceName string) ([]byte, string) {
	cacheFilename := getCacheFilename(text, languageCode, voiceName)
	cacheFilePath := path.Join(cacheDir, cacheFilename)

	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return nil, cacheFilePath
	}

	audioData, err := ioutil.ReadFile(cacheFilePath)

	if err != nil {
		log.Fatal(err)
	}

	return audioData, cacheFilePath
}

func cacheAudioData(audioData []byte, cacheDir string, text string, languageCode string, voiceName string) {
	cacheFilename := getCacheFilename(text, languageCode, voiceName)
	cacheFilePath := path.Join(cacheDir, cacheFilename)

	log.Printf("Caching audio data. File: %s\n", cacheFilePath)
	ioutil.WriteFile(cacheFilePath, audioData, os.ModePerm)
}

func main() {
	languageCode := flag.String("language", "en-US", "set the language code of the input text")
	voiceName := flag.String("voice-name", "en-US-Standard-C", "set the name of the voice which should be used")
	cacheDir := flag.String("cache-dir", "/tmp/gocloudtts/cache", "set the cache directory for voice files")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	text := strings.TrimSpace(flag.Arg(0))

	ensureCacheDirExists(*cacheDir)
	audioData, cacheFilePath := findCachedAudioData(*cacheDir, text, *languageCode, *voiceName)
	if audioData != nil {
		log.Printf("Cache hit! File: %s\n", cacheFilePath)
		os.Stdout.Write(audioData)
		os.Exit(0)
	}

	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx)

	if err != nil {
		log.Fatal(err)
	}

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: *languageCode,
			Name:         *voiceName,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)

	if err != nil {
		log.Fatal(err)
	}

	cacheAudioData(resp.AudioContent, *cacheDir, text, *languageCode, *voiceName)
	os.Stdout.Write(resp.AudioContent)
}
