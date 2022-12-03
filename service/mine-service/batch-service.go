package mineservice

import (
	"fmt"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
  "github.com/workshopapps/pictureminer.api/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type UrlOne struct{
  Url string `bson:"url" json:"url"`
}

//main
type TagOne struct{
  Tag string `bson:"tag" json:"tag"`
  Data []UrlOne `bson:"data" json:"data"`
}

func GetbatchImages(userId string, batchId string) ([]TagOne, error) {

batch_id_primitive , _ := primitive.ObjectIDFromHex(batchId)
var response []TagOne

  tags, length , err := mongodb.GetUserTags(userId ,batch_id_primitive)
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

func filterTags(length int, image_collection []model.ImageCollection ,tag []string,tags []string) []TagOne{

  var tagone TagOne
  var urlone UrlOne
  var str []TagOne
  for k := 0; k < length ; k++ {
  for i , test:= range image_collection{

    tagone.Tag = tag[k]
    if tag[i] == tags[k]{
      urlone.Url = test.Url
      tagone.Data = append(tagone.Data , urlone )
    }
  }
	str = append(str, tagone)
	fmt.Println(str)

  }
return str
}
