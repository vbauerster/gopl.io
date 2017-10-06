package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	// q := url.QueryEscape(strings.Join(terms, " "))
	// q := url.Values{"q": {strings.Join(terms, " ")}}
	// q := url.Values{}
	// q.Add("q", strings.Join(terms, " "))
	// fmt.Println(IssueURL + q.Encode())

	req, err := http.NewRequest(http.MethodGet, IssueURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("q", strings.Join(terms, " "))
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	result := new(IssuesSearchResult)
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}
	return result, nil
}
