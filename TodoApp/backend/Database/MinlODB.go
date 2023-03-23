package Config
import (
    "context"
    "fmt"
    "io"
    "log"
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
    "github.com/sirupsen/logrus"
)
const endpoint = "db:9000"
const accessKeyID = "minioadmin"
const secretAccessKey = "minioadmin"
//Create connect
func ConnectDB()(c *minio.Client,err error){
    useSSL := false
    minioClient, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
        Secure: useSSL,
    })
    if err != nil {
        log.Fatalln(err)
    }
    return minioClient, err
}
//Set permission
func SetPermission(client *minio.Client, bucketName string) error{
    policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::`+ bucketName +`/*"],"Sid": ""}]}`
    err := client.SetBucketPolicy(context.Background(), bucketName, policy)
    if err != nil {
        fmt.Println(err)
        return err
    }
    return err
}
//Create bucket
func CreateBucket(client *minio.Client, bucketName string) error{
    ctx := context.Background()
    err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "ap-northeast-1"})
    if err != nil{
        exists, errBucketExists :=client.BucketExists(ctx, bucketName)
        if errBucketExists != nil {
            logrus.Errorf("[UploadImage] check bucket exists error: %s", err)
            return err
        }
        if !exists {
            logrus.Errorf("[UploadImage] make bucket error: %s", err)
            return err
        }
    }
    return err
}
//Upload data to MinIO
func UploadData(client *minio.Client, bucketName, objectName string, data io.Reader) error {
    _, err := client.GetBucketPolicy(context.Background(), bucketName)
    if err != nil {
        log.Fatalln(err)
    }
    n, err := client.PutObject(context.Background(), bucketName, objectName, data, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
    if err != nil {
        fmt.Println(err)
        return err
    }
    fmt.Println("Successfully uploaded bytes: ", n)
    return err
}
//Get a data from MinIO
func GetDataTodo(client *minio.Client, bucketName, objectName string) (file io.Reader) {
    _, err := client.GetBucketPolicy(context.Background(), bucketName)
    if err != nil {
        log.Fatalln(err)
    }
    file, err = client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
    if err != nil {
        fmt.Println(err)
        return 
    }
    return file
}
//Get all data from MinIO
func GetDataTodoList(client *minio.Client, bucketName string) (file []io.Reader) {
    _, err := client.GetBucketPolicy(context.Background(), bucketName)
    if err != nil {
        log.Fatalln(err)
    }
    objectCh := client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
        Recursive: true,
     })
     for object := range objectCh {
        file = append(file, GetDataTodo(client, bucketName, object.Key))
    }
    return file
}