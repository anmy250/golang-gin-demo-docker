package main
import (
    "TodoApp/Database"
    "TodoApp/Controllers"
    models "TodoApp/Models"
    "log"
    "time"
    "github.com/gin-gonic/gin"
)
func main() {
    time.Sleep(3 * time.Second)
    Client, err := Config.ConnectDB()
    if err != nil {
        log.Println(err)
    }
    err = Config.CreateBucket(Client, "todolist")
    if err != nil {
        log.Println(err)
    }
    //example data
    Controllers.UploadJson(Client, "todolist","todo1.json","1","go to school")
    Controllers.UploadJson(Client, "todolist","todo2.json","2","go to canteen")
    Controllers.UploadJson(Client, "todolist","todo3.json","3","come back home")
    //route
    router := gin.Default()
    router.GET("/api/todos", func(c *gin.Context) {
        ab := Controllers.GetAllTodos(Client, "todolist")
        c.JSON(200, ab)
        })
    router.POST("/api/todos/:uid/:id", func(c *gin.Context) {
        var todo models.Todo
        err := c.BindJSON(&todo)
        if err != nil {
            log.Fatal(err)
        }
        res := Controllers.AddTodo(Client, "todolist", c.Param("id")+".json", c.Param("id"), todo.Name)
        c.JSON(200, res)
    })
    router.Run(":8080")
}