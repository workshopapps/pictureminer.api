package mineservice

import (
	"bytes"
	"fmt"
  "encoding/json"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
  "github.com/workshopapps/pictureminer.api/internal/model"
)

type UrlOne struct{
  Url string `bson:"url" json:"url"`
}

//main
type TagOne struct{
  Tag string `bson:"tag" json:"tag"`
  Data []UrlOne `bson:"data" json:"data"`
}


func GetbatchImages(userId string, batchId string) (*bytes.Buffer, error) {
	//  userId , ok := userId.(string)
	// if !ok {
	// 	return TagOne{}, errors.New("invalid userId")
	// }

 //  batchId , ok = batchId.(string)
 // if !ok {
 //   return TagOne{}, errors.New("invalid batchId")
 // }

var response *bytes.Buffer

  tags, length , err := mongodb.GetUserTags(userId ,batchId)
  if err != nil {
    return response, err
  }
  fmt.Println(tags)
  fmt.Println(length)

  image_collection, tag, length1 , err := mongodb.GetImageTags(batchId)
  if err != nil {
    return response, err
  }

  fmt.Println(tag)
  fmt.Println(length1)

response = filterTags(length,image_collection,tag,tags)
	return response, nil

}

func filterTags(length int, image_collection []model.ImageCollection ,tag []string,tags []string) *bytes.Buffer{

  var tagone TagOne
  var urlone UrlOne
  var dd bytes.Buffer
  for k := 0; k < length ; k++ {
  for i , test:= range image_collection{

    tagone.Tag = tag[k]
    if tag[i] == tags[k]{
      urlone.Url = test.Url
      tagone.Data = append(tagone.Data , urlone )
    }
  }
  json.NewEncoder(&dd).Encode(tagone)
  }
return &dd
}
