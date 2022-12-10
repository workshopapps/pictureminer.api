package batchservice

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseSynonymAPI struct {
	Results []struct {
		HasCategories []string `json:"hasCategories,omitempty"`
		MemberOf      string   `json:"memberOf,omitempty"`
		Synonyms      []string `json:"synonyms,omitempty"`
		TypeOf        []string `json:"typeOf,omitempty"`
		HasTypes      []string `json:"hasTypes,omitempty"`
		Derivation    []string `json:"derivation,omitempty"`
		VerbGroup     []string `json:"verbGroup,omitempty"`
	} `json:"results,omitempty"`
}

type Result struct {
	Name  string
	Score float64
}

type Label struct {
	URL     string
	Results []Result
}

type (
	APIResponse struct {
		Result APIResponseResult `json:"result"`
	}

	APIResponseResult struct {
		Tags []Tag `json:"tags"`
	}

	Tag struct {
		Confidence float64           `json:"confidence"`
		Tag        map[string]string `json:"tag"`
	}
)

const (
	UntaggedThreshold = 5
	Untagged          = "untagged"
	statusComplete    = "completed"
)

func processBatch(email, bName, desc, userID string, batchID primitive.ObjectID, tags, urls []string) {
	// labels for each url
	labels := fetchLabelsForURLS(urls)

	// classify label to matching tag
	for i := 0; i < len(tags); i++ {
		tags[i] = strings.TrimSpace(tags[i])
	}
	batchImgs := classifyLabels(batchID.Hex(), labels, tags)

	// save to db
	err := saveToDB(batchImgs)
	if err != nil {
		warnUser(email, bName, err.Error())
	}

	// on complete, notify user through email
	notifyUser(email, bName)

	// update batch status
	updateBatchStatusDB(batchID, email, bName)
}

func updateBatchStatusDB(batchID primitive.ObjectID, email, bName string) {
	// update batch status in db
	database := config.GetConfig().Mongodb.Database
	userCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.BatchCollection)
	filter := bson.M{"_id": batchID}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: statusComplete}}}}
	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		warnUser(email, bName, err.Error())
	}
}

func fetchLabelsForURLS(urls []string) []Label {
	var labels []Label
	ImaggaURL := config.GetConfig().ImaggaAPI.URL
	AuthToken := config.GetConfig().ImaggaAPI.Auth
	httpClient := &http.Client{}
	for _, url := range urls {
		label, err := getLabel(httpClient, ImaggaURL, url, AuthToken)
		if err == nil {
			labels = append(labels, label)
		}
	}
	return labels
}

func getLabel(client *http.Client, imaggaURL, url, authToken string) (Label, error) {
	req, _ := http.NewRequest("GET", imaggaURL, nil)
	req.Header.Set("Authorization", authToken)
	q := req.URL.Query()
	q.Add("image_url", url)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return Label{}, err
	}
	defer res.Body.Close()

	var apiRes APIResponse
	err = json.NewDecoder(res.Body).Decode(&apiRes)
	if err != nil {
		return Label{}, err
	}

	var results []Result
	for _, tag := range apiRes.Result.Tags {
		results = append(results, Result{Name: tag.Tag["en"], Score: tag.Confidence})
	}

	var label = Label{
		URL:     url,
		Results: results,
	}

	return label, nil
}

func classifyLabels(batchID string, labels []Label, tags []string) []model.BatchImage {
	// for O(1) lookups
	tagsMap := make(map[string][]string)
	fieldTags := []string{}
	for _, tag := range tags {
		fieldTags = append(fieldTags, strings.Fields(tag)...)
	}
	for _, tag := range fieldTags {
		tagsMap[strings.ToLower(tag)] = getSynonymsAPI(tag)
	}

	// get best matching tag for each label
	var batchImgs []model.BatchImage
	for _, label := range labels {
		batchImgs = append(batchImgs, classifyLabelImproved(batchID, label, tags, tagsMap))
	}

	return batchImgs
}

func classifyLabelImproved(batchID string, label Label, tags []string, tagsMap map[string][]string) model.BatchImage {
	type MatchResult struct {
		tag      string
		score    float64
		scoreAux float64
	}
	matchResults := []MatchResult{}

	for _, tag := range tags {
		score, scoreAux, tagTokens := 0.0, 0.0, strings.Fields(tag)
		for _, token := range tagTokens {
			a, b := matchTag(label.Results, token, tagsMap)
			score, scoreAux = score+a, scoreAux+b
		}
		// fmt.Println(MatchResult{tag: tag, score: score, scoreAux: scoreAux})
		matchResults = append(matchResults, MatchResult{tag: tag, score: score, scoreAux: scoreAux})
	}

	choice, score := Untagged, 0.0
	for _, res := range matchResults {
		if res.score > score {
			choice = res.tag
			score = res.score
		}
	}

	// consider second choice: synonyms score
	if choice == Untagged || score < UntaggedThreshold {
		choice = Untagged
		for _, res := range matchResults {
			if res.score > score {
				choice = res.tag
				score = res.score
			}
		}
	}

	batchImage := model.BatchImage{
		ID:      primitive.NewObjectID(),
		BatchID: batchID,
		URL:     label.URL,
		Tag:     choice,
	}
	return batchImage
}

func matchTag(results []Result, tag string, tagsMap map[string][]string) (float64, float64) {
	mscore := 0.0
	// simple match criteria: equality, substring, superstring
	for _, res := range results {
		name, tag := strings.ToLower(res.Name), strings.ToLower(tag)
		if name == tag {
			mscore = math.Max(mscore, res.Score)
		}
		if strings.Contains(name, tag) || strings.Contains(tag, name) {
			mscore = math.Max(mscore, res.Score)
		}
	}

	// complex match criteria: synonym with equality, substring, superstring
	mscoreSynonym := 0.0
	for _, res := range results {
		name := strings.ToLower(res.Name)
		for _, synonym := range tagsMap[tag] {
			ok := synonym != "" && (name == synonym || strings.Contains(name, synonym) || strings.Contains(synonym, name))
			if ok {
				mscoreSynonym = math.Max(mscoreSynonym, res.Score)
			}
		}
	}
	return mscore, mscoreSynonym
}

func classifyLabel(batchID string, label Label, tagsMap map[string][]string) model.BatchImage {
	var filtered []Result

	// filter for results that match with atleast one tag
	for _, res := range label.Results {
		name := smoothify(res.Name)
		if _, ok := tagsMap[name]; ok {
			filtered = append(filtered, res)
		}
	}

	// get result with higest confidence/score
	bestTag, maxScore := Untagged, -1.0
	for _, res := range filtered {
		if res.Score > maxScore {
			bestTag, maxScore = res.Name, res.Score
		}
	}

	batchImage := model.BatchImage{
		ID:      primitive.NewObjectID(),
		BatchID: batchID,
		URL:     label.URL,
		Tag:     bestTag,
	}
	return batchImage
}

func getSynonymsAPI(word string) []string {
	synonyms := []string{}
	URL := config.GetConfig().WordsAPIConfig.URL
	trailer := config.GetConfig().WordsAPIConfig.Trailer

	resp, err := http.Get(URL + word + trailer)
	if err != nil || resp.StatusCode != http.StatusOK {
		return synonyms
	}

	respAPI := ResponseSynonymAPI{}
	json.NewDecoder(resp.Body).Decode(&respAPI)

	for _, val := range respAPI.Results {
		synonyms = append(synonyms, val.Synonyms...)
		synonyms = append(synonyms, val.HasTypes...)
		synonyms = append(synonyms, val.TypeOf...)
		synonyms = append(synonyms, val.HasCategories...)
		synonyms = append(synonyms, val.MemberOf)
		synonyms = append(synonyms, val.VerbGroup...)
	}

	return synonyms
}

func smoothify(str string) string {
	return strings.Join(strings.Fields(str), "-")
}

func saveToDB(batchImgs []model.BatchImage) error {
	database := config.GetConfig().Mongodb.Database
	batchImgCol := mongodb.GetCollection(mongodb.Connection(), database, constants.BatchImageCollection)

	// mongodb insertMany supports []interface{} only
	imgs := make([]interface{}, len(batchImgs))
	for i := 0; i < len(imgs); i++ {
		imgs[i] = batchImgs[i]
	}

	_, err := batchImgCol.InsertMany(context.TODO(), imgs)
	if err != nil {
		return err
	}

	return nil
}

func notifyUser(email, batchName string) {
	from := config.GetConfig().NotifyEmail.Email
	password := config.GetConfig().NotifyEmail.Email
	body := fmt.Sprintf("Hello, batch <b>%v</b> processing is complete!", batchName)

	utility.EmailSender(from, password, []string{email}, "Process Batch Complete", body)
}

func warnUser(email, batchName, msg string) {
	from := config.GetConfig().NotifyEmail.Email
	password := config.GetConfig().NotifyEmail.Email
	body := fmt.Sprintf("Hello, batch <b>%v</b> processing is failed!<br>error: %v", batchName, msg)

	utility.EmailSender(from, password, []string{email}, "Process Batch Failed", body)
}

/*--------------------*/

func classifyAlgorithm(resultId string, label Label, tagMap map[string]bool) model.BatchImage {
	var filterResult []Result

	for _, res := range label.Results {
		name := smoothify(res.Name)
		if tagMap[name] {
			filterResult = append(filterResult, res)
		}
	}

	bestTag, maxScore := Untagged, -1.0
	for _, res := range filterResult {
		if res.Score > maxScore {
			bestTag, maxScore = res.Name, res.Score
		}
	}

	batchImage := model.BatchImage{
		ID:      primitive.NewObjectID(),
		BatchID: resultId,
		URL:     label.URL,
		Tag:     bestTag,
	}
	return batchImage

}

func labelClassifier(resultId string, labels []Label, tags []string) []model.BatchImage {
	// for O(1) lookups
	tagsMap := make(map[string]bool)
	for _, tag := range tags {
		tagsMap[smoothify(tag)] = true
	}

	// get best matching tag for each label
	var batchImgs []model.BatchImage
	for _, label := range labels {
		batchImgs = append(batchImgs, classifyAlgorithm(resultId, label, tagsMap))
	}

	return batchImgs
}
