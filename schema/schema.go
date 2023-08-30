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
}

type TodoList []Todo

func initTodoList() TodoList {
  var list TodoList

  todo1 := Todo{"a", "asdasdad", false}
  todo2 := Todo{"b", "asdasdasdasdasd", false}
  todo3 := Todo{"c", "Hahahahasdasds", false}

  list = append(list, todo1, todo2, todo3)
  fmt.Println("List: ", list)

  return list
}

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
  Name: "RootQuery",
  Fields: graphql.Field{
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
        for _, todo := range TodoList {
          if todo.id == idQuery {
            return todo, nil
          }
        }
      }

      return Todo{}, errors.New("Todo does not exist")
    },
  },
})

var TodoSchema, err = graphql.NewSchema(graphql.SchemaConfig{
  Query: rootQuery,
  // Mutation: rootMutation,
})

