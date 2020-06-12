# GoCloudTTS

Translate text to speech using Google Cloud on the command line.

## Installation

wip

## Usage

Make sure to provide your Google Cloud credentials by setting the environment variable `GOOGLE_APPLICATION_CREDENTIALS`
as described [on this page](https://cloud.google.com/docs/authentication/getting-started).

```
$ GoCloudTTS --help

Usage: GoCloudTTS [OPTIONS] text

OPTIONS:

  -cache-dir string
        set the cache directory for voice files (default "/tmp/gocloudtts/cache")
  -language string
        set the language code of the input text (default "en-US")
  -voice-name string
        set the name of the voice which should be used (default "en-US-Standard-C")
```

You can find a list of supported languages and voices [here](https://cloud.google.com/text-to-speech/docs/voices).

### Example

The following command will translate the text "Peter Piper picked a peck of pickled peppers" and play it through `aplay`.

```
$ GoCloudTTS "Peter Piper picked a peck of pickled peppers" | aplay
```

To change the language use the parameters `-language` and `-voice-name`.

```
$ GoCloudTTS -language ja-JP -voice-name ja-JP-Standard-A "おまえ わ もう しんでいる! なに?!"
```
