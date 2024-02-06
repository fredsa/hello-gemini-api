# Hello, Gemini API.

Sample app showing how to access Gemini using the Gemini API.


## Prerequisites

1.  Download and install Go, see https://go.dev/. To determine the version of Go
    that is available on your path run:

    ```sh
    go version
    ```

1.  [Optional] Install `gcloud` CLI, see https://cloud.google.com/sdk/docs/install.


## Prepare sample project

1.  Clone this repo.

    ```shell
    git clone git@github.com:fredsa/hello-gemini-api.git
    ```

1.  Change into the project directory.

    ```shell
    cd hello-gemini-api
    ```

1.  Install the
    [`github.com/google/generative-ai-go/genai`](https://pkg.go.dev/github.com/google/generative-ai-go/genai)
    package.

    ```shell
    go get github.com/google/generative-ai-go/genai
    ```


## Authorize the app

1.  Create API key from [ai.google.dev](https://ai.google.dev/).

2.  Modify `main.go` to use this API key.

    ```go
    const apiKey = "your-api-key"  // Keep safe.
    ```

3. Keep your API key safe.


## Run the sample

1. Compile and run.

    ```shell
    go run main.go
    ```

    The output should look something like this
    ```none
    >> Hello, who are you?
    I am a large language model, trained by Google.
    ```