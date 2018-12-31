// go run crawler_para.go <file> <threads>
// 	file = the input file name(place the file to the folder of crawler_seq.go
// 	before running the program)
// 	threads = the number of threads(including the main thread)
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
	"strconv"
	"sync"
)
const QuestionUrl =  "https://stackoverflow.com/questions"
const TagUrl = "/questions/tagged"
const NextUrl = "/questions"
const mainThreadBank = 0 
const blockSize = 20   //every thread process blocksize urls every time
type Link struct {
	url   string
	text  string
}
type BlockInfo struct{
	large int
	small int
	numberOflarge int
	total int
}
type goContext struct{
	tagChan chan []string
	questionChan chan []string
	workers int
	questionLock *sync.Mutex
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
	return fmt.Sprintf(" %s - %s", self.text, self.url)
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


// worker function to parse urls and generate tags for each url
func analyze(tags map[string]int, questions map[string]bool, context *goContext, i int) {
	// check if in the main thread
	if i == mainThreadBank{
		for {
			oneTagMap, more := <- context.tagChan
			if !more{
				return
			}else{
				for k := 0; k < len(oneTagMap); k++{
					if _, ok := tags[oneTagMap[k]]; !ok{
						tags[oneTagMap[k]] = 1
					}else{
						tags[oneTagMap[k]] += 1
					}
				}
			}
		}
	}else{
		var tagParts []string
		for {
			questionNum, more := <- context.questionChan
			if !more{
				context.tagChan <- tagParts
				context.workers -= 1
				if context.workers == 0 {
					close(context.tagChan)
				}
				return
			}
			var validQuestionNum []string
			context.questionLock.Lock()
			for j:=0; j< len(questionNum) ;j++{
				if _,ok := questions[questionNum[j]]; !ok{
					validQuestionNum = append(validQuestionNum,questionNum[j])
					questions[questionNum[j]] = true
				}
			}
			context.questionLock.Unlock()
			for j:=0; j< len(validQuestionNum) ;j++{
				url := QuestionUrl + "/" + validQuestionNum[j]
				resp, err := parse(url)
				if err != nil {
					fmt.Printf("Error getting page %s %s\n", url, err)
					return
				}
				tagsPerLink := pageLink(resp)
				for k, _ := range tagsPerLink{
					tagParts = append(tagParts,k)
				}
			}

		}
	}

	
}
func main(){
	fileName := os.Args[1]
	numOfThreads,_:= strconv.Atoi(os.Args[2])
	var tags = make(map[string]int)
	var questions = make(map[string]bool)
	var lock = sync.Mutex{}
	var i = mainThreadBank + 1
	var questionSet []string

	context := goContext{ tagChan: make(chan []string), questionChan: make(chan []string),workers: numOfThreads - 1, questionLock:&lock}
	for ; i < numOfThreads; i++{
		go analyze(tags, questions, &context, i)
	}
	file,err := os.Open(fileName)
	if err != nil{
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	var counter = 0;
	for scanner.Scan(){
		questionNum := scanner.Text()
		if questionNum == "" {
			fmt.Println("Invalid Url")
		}else if counter < blockSize{
			counter++;
			questionSet = append(questionSet,scanner.Text())
		}else{
			context.questionChan <- questionSet
			counter = 0
		}
	}
	if counter >0{
		context.questionChan <- questionSet
	}
	close(context.questionChan)
	analyze(tags, questions, &context,mainThreadBank)
	for k, v := range tags{
		fmt.Printf("%s : %d\n",k,v)
	}
}