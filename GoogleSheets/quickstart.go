package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/mgo.v2/bson"

	sheets "google.golang.org/api/sheets/v4"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	b64 "encoding/base64"
)

// StateParams is used to represent the extra query parameters passed to the redirect url
type StateParams struct {
	Authorization string `json:"Authorization"`
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config, extraQueryParams []byte) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	log.Printf("getclient tok err: %v \n\n", err)

	if err != nil || true {
		fmt.Printf("en get tkoen from web\n")
		tok = getTokenFromWeb(config, extraQueryParams)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config, extraQueryParams []byte) *oauth2.Token {
	qpEnc := b64.URLEncoding.EncodeToString(extraQueryParams)
	authURL := config.AuthCodeURL(qpEnc, oauth2.AccessTypeOffline)
	log.Printf("stateQueryParams: %v\nextraQueryParams: %v\n", qpEnc, extraQueryParams)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	fmt.Printf("tok: %v \nerr: %v\n", tok, err)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("sheets.googleapis.com-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	log.Printf("tokenFromFile() file: %v \n\n", file)
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)

	log.Printf("tokenFromFile()\nt.AccessToken: %v \nt.tokenType: %v\nt.RefreshToken: %v \nt.expiry: %v\n\n", t.AccessToken, t.TokenType, t.RefreshToken, t.Expiry)

	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {

	fmt.Printf("saveToken() oauth2.Token: %v \nType: %T\n", token, token)
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	for i := 0; i < 5; i++ {
		fmt.Printf("bsonid: %v \n", bson.NewObjectId())
	}
	sp := StateParams{}
	sp.Authorization = "token blah"
	spJSONMarshal, _ := json.Marshal(sp)
	spEnc := b64.URLEncoding.EncodeToString([]byte(spJSONMarshal))
	log.Printf("spEnc: %v \n\n", spEnc)
	// verify: uEnc := "eyJBdXRob3JpemF0aW9uIjoidG9rZW4gYmxhaCJ9"
	decodedstate, _ := b64.URLEncoding.DecodeString(spEnc)
	log.Printf("decodedState: %v\nstring: %v\n", decodedstate, string(decodedstate))
	log.Printf("spJSON Marshal: %v type: %T\n", spJSONMarshal, spJSONMarshal)
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")

	sEnc := b64.StdEncoding.EncodeToString(b)
	fmt.Printf("Client secret base64 string : %v \n", sEnc)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	log.Printf("b: %v \n\n", b)

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config, []byte(spJSONMarshal))

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	//spreadsheetId := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"

	// https://docs.google.com/spreadsheets/d/1zFjra05ZGfaVgKNorPdvAU-bh0QDkOn-CVoXjWtiw2w/edit
	spreadsheetID := "1zFjra05ZGfaVgKNorPdvAU-bh0QDkOn-CVoXjWtiw2w"

	//readRange := "Class Data!A2:E"
	readRange := "A1:F6"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

	log.Printf("\nresp: %v\n\nresp.Values: %v\n\nresp.Values type: %T \n\n", resp, resp.Values, resp.Values)

	if len(resp.Values) > 0 {
		fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			// fmt.Printf("%s, %s\n", row[0], row[4])
			fmt.Printf("Row Type: %T\n", row)
			for _, item := range row {
				fmt.Printf("%s \t", item)
			}
			fmt.Printf("\n")
		}
	} else {
		fmt.Print("No data found.")
	}

	file, err := os.Create("gsheet_result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range resp.Values {
		cur := []string{}
		for _, item := range value {
			cur = append(cur, item.(string))
		}
		err := writer.Write(cur)
		checkError("Cannot write to file", err)
	}

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
