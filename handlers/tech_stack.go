package handlers

import (
	"io"
	"net/http"

	wappalyzer "github.com/projectdiscovery/wappalyzergo"
)

func getTechStack(domain string) (techStackList []string, err error) {

	// Use http client to get the received domain
	resp, err := http.DefaultClient.Get("https://" + domain)
	if err != nil {
		return nil, err
	}
	data, _ := io.ReadAll(resp.Body)

	// Fingerprint the technology on a page
	wappalyzerClient, err := wappalyzer.New()
	fingerprints := wappalyzerClient.Fingerprint(resp.Header, data)

	// Make empty list
	techStackList = make([]string, 0)

	// Add the received map keys to list
	for value := range fingerprints {
		techStackList = append(techStackList, value)
	}

	return
}

func HandleTechStack() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawUrl, err := extractURL(r)
		if err != nil {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}

		techStack, err := getTechStack(rawUrl.Hostname())

		if err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}

		JSON(w, KV{
			"tech-stack": techStack,
		}, http.StatusOK)

	})
}