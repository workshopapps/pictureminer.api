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

func filterTags(length int, image_collection []model.BatchImage ,tag []string,tags []string) []TagOne{

  var str []TagOne

  for k := 0; k < length ; k++ {
			var tagone TagOne
  for i , test:= range image_collection{
			tagone.Tag = tags[k]
  if tag[i] == tags[k]{
			fmt.Println(tag[i]+tags[k])
			fmt.Println(test.URL)
			 var urlone UrlOne
	     urlone.Url = test.URL
	     tagone.Data = append(tagone.Data , urlone )
			fmt.Println(tagone.Data)
    	}
  	}
	str = append(str, tagone)
  }

return str
}
