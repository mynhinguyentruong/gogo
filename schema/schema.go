package schema

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
)

type Todo struct {
  ID string `json:"id"`
  Text string `json:"text"`
  Completed bool `json:"completed"` 
  User_email string `json:"user_email,omitempty"` 
}

type User struct {
  ID string 
  Method string
  Email string
  Image_url string

}

type TodoList []Todo

var list TodoList = InitTodoList()

func InitTodoList() TodoList {
  var list TodoList

  todo1 := Todo{"a", "asdasdad", false, "Hahahahasdasds@sgm.ca"}
  todo2 := Todo{"b", "asdasdasdasdasd", false, "Hahahahasdasds@sgm.ca"}
  todo3 := Todo{"c", "Hahahahasdasds", false, "Hahahahasdasds@sgm.ca"}

  list = append(list, todo1, todo2, todo3)
  fmt.Println("List: ", list)

  return list
}

var todoType = graphql.NewObject(graphql.ObjectConfig{
  Name: "Todo",
  Fields: graphql.Fields{
    "id": &graphql.Field{
      Type: graphql.String,
    },
    "text": &graphql.Field{
      Type: graphql.String,
    },
    "completed": &graphql.Field{
      Type: graphql.Boolean,
    },
    "user_email": &graphql.Field{
      Type: graphql.String,
    },
  },
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
  Name: "RootQuery",
  Fields: graphql.Fields{
      "todo": &graphql.Field{
        Type: todoType,
        Description: "Get a single todo",
        Args: graphql.FieldConfigArgument{
          "id": &graphql.ArgumentConfig{
            Type: graphql.String,
          },
        },
        Resolve: func (params graphql.ResolveParams) (interface{}, error) {
          idQuery, isOk := params.Args["id"].(string)
          
          if isOk {
            for _, todo := range list {
              if todo.ID == idQuery {
                return todo, nil
              }
            }
          }

          return Todo{}, errors.New("Todo does not exist")
      },
    },
  },
})

var TodoSchema, err = graphql.NewSchema(graphql.SchemaConfig{
  Query: rootQuery,
  // Mutation: rootMutation,
})

