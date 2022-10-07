package main

import (
	"context"
	"fmt"
	"regexp"
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
		300*time.Second)

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
		fmt.Print("1.Action(create_user, find_user, findAll_user)\n2.Exit\n")
		fmt.Scanln(&choice)
		if choice == 1 {
			var ch int
			fmt.Println("Select the detail you want to update")
			fmt.Print("1.Create User\n2.Find User\n3.Find All User\n")
			fmt.Scanln(&ch)
			if ch == 1 {
				var fname, lname, email, state string
				var regexpEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

				fmt.Println("Enter the Firstname of user:")
				fmt.Scanln(&fname)
				fmt.Println("Enter the Lastname of user:")
				fmt.Scanln(&lname)
				fmt.Println("Enter the Email of user:")
				fmt.Scanln(&email)
				fmt.Println("Enter the State of user:")
				fmt.Scanln(&state)
				// validating input parameters
				// fname, lname and state will be on same parameters
				// checking on the basis of empty
				if fname == "" || lname == "" || email == "" || state == "" {
					fmt.Println("Data is not validated, all inputs feilds are required.")
				} else if len(fname) < 2 || len(fname) > 40 || len(lname) < 2 || len(lname) > 40 || len(email) < 2 || len(email) > 40 || len(state) < 1 || len(state) > 40 {
					fmt.Println("Data is not validated, all input fields must be entered in Full.")
				} else if !regexpEmail.MatchString(email) {
					fmt.Println("Data is not validated, please enter a valid email address.")
				} else {
					fmt.Println("Data is validated successfully")
					var data interface{}
					data = bson.D{
						{"Firstname", fname},
						{"Lastname", lname},
						{"Email", email},
						{"State", state},
						{"Company", "MYday"},
					}

					insertOne(client, ctx, "User", "users", data)
					fmt.Println("User added successfully")
				}
			} else if ch == 2 {
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

			} else if ch == 3 {
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

			} else {
				fmt.Println("Invalid choice")
			}

		} else if choice == 2 {
			fmt.Println("Thankyou!")
			break
		} else {
			fmt.Println("Invalid choice! Enter Again")
		}
	}
}
