package mineservice

import (
	"fmt"
	// "github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)


func GetMonthlylimit(userId string) (bool, error) {

	var response bool
  //
	plan, err := mongodb.GetUserPlan(userId)
	if err != nil {
		return response, err
	}
	fmt.Println(plan)

	_, time_collection , _ , err := mongodb.GetMinedTime(userId)
	if err != nil {
		return response, err
	}

	// fmt.Println(image_collection)
	// fmt.Println(time_collection)

  _, time_batch_collection , _ , err := mongodb.GetBatchTime(userId)
  if err != nil {
    return response, err
  }

  // fmt.Println(batch_collection)
  // fmt.Println(time_batch_collection)

  totalMined := len(time_collection) + len(time_batch_collection)
  fmt.Println(totalMined)
  if plan == "free" || plan == "" {
    if totalMined >= 10 {
      response = false
    }else{
      response = true
    }
  }else{
    response = true
  }

	return response, nil

}
