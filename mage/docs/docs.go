package docs

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	// Triggers a build of the v1 website
	porterV1Webhook = "https://api.netlify.com/build_hooks/60ca5ba254754934bce864b1"
)

// Deploy triggers a Netlify build for the website.
func Deploy() error {
	// Put up a page on the preview that redirects to the live site
	os.MkdirAll("docs/public", 0755)
	err := copy("hack/website-redirect.html", "docs/public/index.html")
	if err != nil {
		return err
	}

	return TriggerNetlifyDeployment(porterV1Webhook)
}

func copy(src string, dest string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}

	destF, err := os.Create(dest)
	if err != nil {
		return err
	}

	_, err = io.Copy(destF, srcF)
	return err
}

// TriggerNetlifyDeployment builds a netlify site using the specified webhook
func TriggerNetlifyDeployment(webhook string) error {
	emptyMsg := "{}"
	data := strings.NewReader(emptyMsg)
	fmt.Println("POST", webhook)
	fmt.Println(emptyMsg)

	r, err := http.Post(webhook, "application/json", data)
	if err != nil {
		return err
	}

	if r.StatusCode >= 300 {
		defer r.Body.Close()
		msg, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("request failed (%d) %s: %s", r.StatusCode, r.Status, msg)
	}

	return nil
}
