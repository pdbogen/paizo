package paizo

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/http/cookiejar"
	"github.com/pdbogen/paizo/dom"
)

const index_page_url = "https://paizo.com"
const signin_page_url = "https://secure.paizo.com/cgi-bin/WebObjects/Store.woa/wa/DirectAction/signIn?path=paizo"

func NewSession() (sess Session, err error) {
	if sess.jar == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			return Session{}, fmt.Errorf("creating cookie jar?! %s", err)
		}
		sess.jar = jar
	}

	client := http.Client{}
	client.Jar = sess.jar
	index, err := client.Get(index_page_url)
	if err != nil {
		return Session{}, fmt.Errorf("retrieving index page: %s", err)
	}
	if index.StatusCode/100 != 2 {
		return Session{}, fmt.Errorf("status code from / was %d, not 2XX", index.StatusCode)
	}

	return sess, nil
}

type Session struct {
	jar *cookiejar.Jar
}


func (s Session) Authenticate(username, password string) (err error) {
	client := http.Client{}
	client.Jar = s.jar
	signinPage, err := client.Get(signin_page_url)
	if err != nil {
		return fmt.Errorf("retrieving signin page: %s", err)
	}
	defer signinPage.Body.Close()

	if signinPage.StatusCode/100 != 2 {
		return fmt.Errorf("status code from GET %s was %d, not 2XX", signin_page_url, signinPage.StatusCode)
	}

	//body, err := ioutil.ReadAll(signinPage.Body)
	//_, err = html.Parse(bytes.NewBuffer(body))
	doc, err := html.Parse(signinPage.Body)
	if err != nil {
		return fmt.Errorf("parsing signin page: %s", err)
	}

	// Look for form with "Sign In" button
	oInput := dom.Find(doc, dom.WithTag("input").And(dom.WithAttribute("name", "o")))
	if oInput == nil {
		return errors.New("could not find input with name `o` in document")
	}

	//oVal := dom.Attribute(oInput, "value")

	form := dom.FindParent(oInput, dom.WithTag("form"))

	if form == nil {
		return errors.New("couldn't find parent form of `o` input?")
	}

	fmt.Println("Got form ", form)
	return nil
}
