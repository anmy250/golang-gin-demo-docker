package Controllers
import (
    "TodoApp/Database"
    models "TodoApp/Models"
    "bytes"
    "github.com/minio/minio-go/v7"
)
func GetAllTodos(Client *minio.Client, bucketName string) (res models.TodoList) {
    todoList := Config.GetDataTodoList(Client, bucketName)
    res = models.GetAllTodos(todoList)
    return res
}
func AddTodo(Client *minio.Client, bucketName, objectName, id, name string)(res models.Todo){
    jsonData,res := models.CreateJson(id, name)
    data := bytes.NewReader(jsonData)
    err := Config.UploadData(Client, bucketName,objectName, data)
    if err !=nil{
        panic(err)
    }
    return res
}
func UploadJson(Client *minio.Client, bucketName, objectName, id, name string){
    jsonFile,_:=models.CreateJson(id, name)
    data := bytes.NewReader(jsonFile)
    err := Config.UploadData(Client, bucketName,objectName, data)
    if err !=nil{
        panic(err)
    }
}