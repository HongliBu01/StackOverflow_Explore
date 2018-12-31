// go run crawler_seq.go <file>
// 	file = the input file name(place the file to the folder of crawler_seq.go
// 	before running the program)
package main
import(
	"os"
	"bufio"
	"log"
	"fmt"
	"strings"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
)
const QuestionUrl =  "https://stackoverflow.com/questions"
const TagUrl = "/questions/tagged"
const NextUrl = "/questions"

type Link struct {
	url   string
	text  string
}
//  extract <a href></a> from htmltoken
func NewLink(tag html.Token, text string) Link {
	link := Link{text: strings.TrimSpace(text)}

	for i := range tag.Attr {
		if tag.Attr[i].Key == "href" {
			link.url = strings.TrimSpace(tag.Attr[i].Val)
		}
	}
	return link
}
func (self Link) String() string {
	return fmt.Sprintf("%s - %s", self.text,  self.url)
}
// check if link valid or not
func (self Link) Valid() bool {
	if len(self.text) == 0 {
		return false
	}
	if len(self.url) == 0 || strings.Contains(strings.ToLower(self.url), "javascript") {
		return false
	}

	return true
}
// decode c%2b%2b c%23 to c++ and c#
func decodeSpecialSym(input string) string{
	var output string
	if strings.Contains(input, "%2b"){
		output = strings.Replace(input, "%2b", "+", -1)
	}else if strings.Contains(input, "%23"){
		output = strings.Replace(input, "%23", "#", -1)
	}else{
		output = input
	}
	return output
}

// parse url through net/http
func parse(url string) (resp *http.Response, err error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Cannot get page")
	}

	return r, err
}

// traverse whole webpage and find tag link
func pageLink(resp *http.Response) (map[string]int) {
	tags := make(map[string]int)
	page := html.NewTokenizer(resp.Body)
	var start *html.Token
	var text string

	for{
		_ = page.Next()
		token := page.Token()
		if token.Type == html.ErrorToken {
			break
		}
		if start != nil && token.Type == html.TextToken {
			text = fmt.Sprintf("%s%s", text, token.Data)
		}
		if token.DataAtom == atom.A{
			switch token.Type {
			case html.StartTagToken:
				if len(token.Attr) > 0{
					start = &token
				}

			case html.EndTagToken:
				if start == nil {
					continue
				}
				link := NewLink(*start, text)
				if link.Valid() {
					if strings.HasPrefix(link.url, NextUrl){
						if strings.HasPrefix(link.url, TagUrl){
							tagparts := strings.Split(link.url, "/")
							if len(tagparts) == 4{
								decodedTag := decodeSpecialSym(tagparts[3])
								if _,ok := tags[decodedTag]; !ok{
									tags[decodedTag] =1
								}
							}
						}
					}
				}
				start = nil
				text = ""
			}
		}
	}
	return tags
}

// parse urls and generate tags for each url
func analyze(tags map[string]int, questions map[string]bool, questionNum string) {
	if _, ok := questions[questionNum]; !ok{
		questions[questionNum] = true
	}else{
		return
	}
	url := QuestionUrl + "/" + questionNum
	resp, err := parse(url)
	if err != nil {
		fmt.Printf("Error getting page %s %s\n", url, err)
		return
	}
	tagsPerLink := pageLink(resp)
	for k,_ := range tagsPerLink{
		if _, ok := tags[k]; !ok{
			tags[k] = 1
		}else{
			tags[k] += 1
		}
	}

}

func main(){
	fileName := os.Args[1]
	file,err := os.Open(fileName)
	tags := make(map[string]int)
	questions := make(map[string]bool)
	if err != nil{
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		questionNum := scanner.Text()
		if questionNum == "" {
			fmt.Println("Invalid Url")
		}else{
			analyze(tags, questions, questionNum)
		}
	}
	for k, v := range tags{
		fmt.Printf("%s : %d\n",k,v)
	}
}