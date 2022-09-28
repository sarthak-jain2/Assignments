package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// function to close for deadline process
func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	defer cancel()

	defer func() {

		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// function to connect to database
func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// function to insert one result
func insertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) {

	collection := client.Database(dataBase).Collection(col)

	collection.InsertOne(ctx, doc)

}

// function to find all and one user
func query(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {

	collection := client.Database(dataBase).Collection(col)

	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

// function to update the user details using email
func UpdateOne(client *mongo.Client, ctx context.Context, dataBase, col string, filter, update interface{}) {

	collection := client.Database(dataBase).Collection(col)

	collection.UpdateOne(ctx, filter, update)
	fmt.Println("Details updated successfully")
	return
}

// function to delete the data using email
func deleteOne(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}) {

	collection := client.Database(dataBase).Collection(col)

	collection.DeleteOne(ctx, query)
	fmt.Println("Data deleted successfully")

}

// our mainfunction
func main() {
	//first we will make the connection with database
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)

	// our code to add the data
	var choice int
	fmt.Println("Hi! Welcome to our database")
	fmt.Println("Please enter your choice:")

	for true {
		fmt.Print("1.Create DB\n2.Find all Users\n3.Find User by email\n4.Update User\n5.Delete User\n6.Exit\n")
		fmt.Scanln(&choice)
		if choice == 1 {
			var data interface{}
			data = bson.D{
				{"Firstname", "Shaun"},
				{"Lastname", "Jain"},
				{"Email", "shaun_jain@persistent.com"},
				{"State", "MP"},
				{"Company", "Persistent"},
			}

			insertOne(client, ctx, "User", "users", data)
			fmt.Println("User added successfully")

		} else if choice == 2 {
			println("The details of all users are:")
			var filter, option interface{}
			filter = bson.D{
				{"Company", "Persistent"},
			}
			option = bson.D{{"_id", 0}}
			cursor, err := query(client, ctx, "User", "users", filter, option)
			// handle the errors.
			if err != nil {
				panic(err)
			}
			var results []bson.D
			if err := cursor.All(ctx, &results); err != nil {
				panic(err)
			}
			for _, data := range results {
				fmt.Println(data)
			}

		} else if choice == 3 {
			var email string
			fmt.Println("Enter the email of user to find the record:")
			fmt.Scanln(&email)
			var filter, option interface{}
			filter = bson.D{
				{"Email", email},
			}
			option = bson.D{{"_id", 0}}
			cursor, err := query(client, ctx, "User", "users", filter, option)
			// handle the errors.
			if err != nil {
				panic(err)
			}
			var results []bson.D
			if err := cursor.All(ctx, &results); err != nil {
				panic(err)
			}
			fmt.Println("The requested data is:")
			for _, data := range results {
				fmt.Println(data)
			}

		} else if choice == 4 {
			var email string
			fmt.Println("Enter the email of user to update the details:")
			fmt.Scanln(&email)
			filter := bson.D{
				{"Email", email},
			}
			var ch int
			fmt.Println("Select the detail you want to update")
			fmt.Print("1.Firstname\n2.Lastname\n3.Email\n4.State\n")
			fmt.Scanln(&ch)
			if ch == 1 {
				var detail string
				fmt.Println("Enter Firstname to update:")
				fmt.Scanln(&detail)
				update := bson.D{
					{"$set", bson.D{
						{"Firstname", detail},
					}},
				}
				UpdateOne(client, ctx, "User", "users", filter, update)
			} else if ch == 2 {
				var detail string
				fmt.Println("Enter Lastname to update:")
				fmt.Scanln(&detail)
				update := bson.D{
					{"$set", bson.D{
						{"Lastname", detail},
					}},
				}
				UpdateOne(client, ctx, "User", "users", filter, update)

			} else if ch == 3 {
				var detail string
				fmt.Println("Enter Email to update:")
				fmt.Scanln(&detail)
				update := bson.D{
					{"$set", bson.D{
						{"Email", detail},
					}},
				}
				UpdateOne(client, ctx, "User", "users", filter, update)
			} else if ch == 4 {
				var detail string
				fmt.Println("Enter State to update:")
				fmt.Scanln(&detail)
				update := bson.D{
					{"$set", bson.D{
						{"State", detail},
					}},
				}
				UpdateOne(client, ctx, "User", "users", filter, update)

			} else {
				fmt.Println("Invalid choice")
			}
		} else if choice == 5 {
			var email string
			fmt.Println("Enter the email of user to delete the data:")
			fmt.Scanln(&email)
			filter := bson.D{
				{"Email", email},
			}
			deleteOne(client, ctx, "User", "users", filter)
		} else if choice == 6 {
			fmt.Println("Thankyou!")
			break
		} else {
			fmt.Println("Invalid choice! Enter Again")
		}
	}
}
