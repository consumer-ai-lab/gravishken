package helper

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"backend/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Add_Model_To_DB(Collection *mongo.Collection, Model types.ModelInterface) error {

	fmt.Println("Adding " + Model.GetCollectionName() + " to database...")
	_, err := Collection.InsertOne(context.TODO(), Model)

	if err != nil {
		fmt.Println("Error in adding Model to database: ", err)
		return err
	}

	log.Default().Println(Model.GetCollectionName() + " added successfully")

	return nil
}

func Get_All_Models(collection *mongo.Collection, modelType types.ModelInterface) ([]types.ModelInterface, error) {
	fmt.Println("Fetching all models from database...")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("Error in fetching all models: ", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []types.ModelInterface

	// Create a new slice of the concrete type
	sliceType := reflect.SliceOf(reflect.TypeOf(modelType).Elem())
	concreteSlice := reflect.MakeSlice(sliceType, 0, 0)
	concreteSlicePtr := reflect.New(concreteSlice.Type())
	concreteSlicePtr.Elem().Set(concreteSlice)

	// Decode into the concrete slice
	if err := cursor.All(context.TODO(), concreteSlicePtr.Interface()); err != nil {
		fmt.Println("Error in decoding models: ", err)
		return nil, err
	}

	// Convert concrete slice to []ModelInterface
	concreteSliceVal := concreteSlicePtr.Elem()
	for i := 0; i < concreteSliceVal.Len(); i++ {
		item := concreteSliceVal.Index(i).Addr().Interface().(types.ModelInterface)
		results = append(results, item)
	}

	fmt.Println("Models in DB:", results)
	return results, nil
}

func GetModelById(collection *mongo.Collection, ID string, modelType types.ModelInterface) (types.ModelInterface, error) {
	fmt.Println("Fetching model from database...")
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	var result types.ModelInterface
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		fmt.Println("Error in fetching model: ", err)
		return nil, err
	}

	fmt.Println("Model in DB:", result)
	return result, nil
}

func GetModelByBatchId(collection *mongo.Collection, ID string, modelType types.ModelInterface) ([]types.ModelInterface, error) {
	fmt.Println("Fetching model from database...")
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	var results []types.ModelInterface
	cursor, err := collection.Find(context.TODO(), bson.M{"batch": objectID})
	if err != nil {
		fmt.Println("Error in fetching model: ", err)
		return nil, err
	}

	sliceType := reflect.SliceOf(reflect.TypeOf(modelType).Elem())
	concreteSlice := reflect.MakeSlice(sliceType, 0, 0)
	concreteSlicePtr := reflect.New(concreteSlice.Type())
	concreteSlicePtr.Elem().Set(concreteSlice)

	// Decode into the concrete slice
	if err := cursor.All(context.TODO(), concreteSlicePtr.Interface()); err != nil {
		fmt.Println("Error in decoding models: ", err)
		return nil, err
	}

	// Convert concrete slice to []ModelInterface
	concreteSliceVal := concreteSlicePtr.Elem()
	for i := 0; i < concreteSliceVal.Len(); i++ {
		item := concreteSliceVal.Index(i).Addr().Interface().(types.ModelInterface)
		results = append(results, item)
	}

	fmt.Println("Models in DB:", results)
	return results, nil
}

func Delete_Model_By_ID(Collection *mongo.Collection, ID string) error {

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	result, err := Collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		log.Printf("Error in deleting Model from database: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no document found with ID: %s", ID)
	}

	log.Printf("Model with ID %s deleted successfully", ID)
	return nil
}

func Delete_ALL_Model(Collection *mongo.Collection) error {

	_, err := Collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("Error in deleting Model from database: ", err)
		return err
	}

	log.Default().Println("Model deleted successfully")

	return nil
}

func Update_Model_By_ID(Collection *mongo.Collection, ID string, Model types.ModelInterface) error {

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	result, err := Collection.ReplaceOne(context.TODO(), bson.M{"_id": objectID}, Model)
	if err != nil {
		log.Printf("Error in updating Model from database: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with ID: %s", ID)
	}

	log.Printf("Model with ID %s updated successfully", ID)
	return nil
}
