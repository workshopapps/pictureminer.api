package mineservice

import (
	"fmt"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UrlOne struct {
	Url string `bson:"url" json:"url"`
}

// main
type TagOne struct {
	Tag  string   `bson:"tag" json:"tag"`
	Data []UrlOne `bson:"data" json:"data"`
}

func GetbatchImages(userId string, batchId string) ([]TagOne, error) {

	batch_id_primitive, _ := primitive.ObjectIDFromHex(batchId)
	var response = make([]TagOne, 0)

	tags, length, err := mongodb.GetUserTags(userId, batch_id_primitive)
	if err != nil {
		return response, err
	}
	fmt.Println(tags)
	fmt.Println(length)

	image_collection, tag, length1, err := mongodb.GetImageTags(batchId)
	if err != nil {
		return response, err
	}

	fmt.Println(tag)
	fmt.Println(length1)

	response = filterTags(length, image_collection, tag, tags)

	return response, nil

}


func filterTags(length int,image_collection []model.BatchImage,tag []string,tags []string) []TagOne{

	str := make([]TagOne, 0)

  for k := 0; k < length ; k++ {
  var tagone TagOne
  tagone.Tag = tags[k]
  for i , test:= range image_collection{
    if tag[i] == tags[k]{
      var urlone UrlOne
      urlone.Url = test.URL
      tagone.Data = append(tagone.Data , urlone )
    }
  }
  str = append(str, tagone)
  }


  DistagsBatch := unique(tags)
  DisttagImage := unique(tag)
  diff := missing(DistagsBatch, DisttagImage)

  fmt.Println(diff)
  var tagone TagOne
  tagone.Tag = "untagged"
  for p := 0; p < len(diff) ; p++ {
      for _ , test:= range image_collection{
          if diff[p] == test.Tag {
            var urlTwo UrlOne
            urlTwo.Url = test.URL
            tagone.Data = append(tagone.Data , urlTwo)
            fmt.Println(test.URL)
          }
        }
      }
      str = append(str, tagone)
      fmt.Println(str)
return str
}


func unique(intSlice []string) []string {
    keys := make(map[string]bool)
    list := []string{}
    for _, entry := range intSlice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}


// missing compares two slices and returns slice of differences
func missing(a, b []string) []string {
    type void struct{}
	// create map with length of the 'a' slice
	ma := make(map[string]void, len(a))
	diffs := []string{}
	// Convert first slice to map with empty struct (0 bytes)
	for _, ka := range a {
		ma[ka] = void{}
	}
	// find missing values in a
	for _, kb := range b {
		if _, ok := ma[kb]; !ok {
			diffs = append(diffs, kb)
		}
	}
	return diffs
}
