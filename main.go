package main

import (
	"beego-project/models"
	_ "beego-project/routers"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/beego/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func init() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("DATABASE_URL environment variable is not set")
	}
	orm.RegisterDriver("postgresql", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", dbURL)

}

func main() {
	orm.RunSyncdb("default", true, false)
	db, err := orm.GetDB()
	if err != nil {
		log.Println("get default DataBase")
	}
	orm.AddAliasWthDB("default", "postgres", db)

	// create bucket if you're going to use aws
	storageType := beego.AppConfig.DefaultString("storagetype", "local")

	if storageType == "s3" {
		awsBucketName, awsEndpoint, awsRegion := models.GetAWSInfo()

		awsCfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(awsRegion),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "dummy")),
		)
		if err != nil {
			log.Fatalf("Cannot load the AWS configs: %s", err)
		}

		// Create the resource client
		client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.UsePathStyle = true
			o.BaseEndpoint = aws.String(awsEndpoint)
		})

		_, err = client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
			Bucket: aws.String(awsBucketName),
		})
		if err != nil {
			log.Fatalf("Failed to create bucket: %s", err)
		}
	}

	fmt.Println("hello from", beego.BConfig.AppName)

	beego.Run()

}
