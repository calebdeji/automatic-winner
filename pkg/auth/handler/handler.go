package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"zate/constants"
	"zate/database"
	"zate/generics"
	"zate/services"

	"zate/pkg/auth/models"
	"zate/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var authDumpCollection *mongo.Collection = database.OpenCollection(database.Client, constants.AUTH_DUMP_COLLECTION)
var userCollection *mongo.Collection = database.OpenCollection(database.Client, constants.USERS_COLLECTION)
var validate = validator.New()

func SignUp(c *gin.Context) {
	defer utils.HandlePanic(c)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var payload models.User

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, generics.GenericResponse{
			Status:  "failed",
			Message: err.Error(),
		})

		return
	}

	defer cancel()

	validationError := validate.Struct(payload)

	if validationError != nil {
		c.JSON(http.StatusBadRequest, generics.GenericResponse{
			Status:  "failed",
			Message: validationError.Error(),
		})

		return
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": payload.Email})

	defer cancel()

	if err != nil {
		log.Panic("Unable to retrieve email")
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, generics.GenericResponse{
			Status:  "failed",
			Message: "Email exists already",
		})
		return
	}

	defer cancel()

	count, err = authDumpCollection.CountDocuments(ctx, bson.M{"email": payload.Email})

	if err != nil {
		log.Panic("Unable to retrieve email")
	}

	if count > 0 {
		authDumpCollection.DeleteOne(ctx, bson.M{"email": payload.Email})
	}

	hashedPassword := utils.HashPassword(*payload.Password)

	payload.ID = primitive.NewObjectID()
	payload.User_id = payload.ID.Hex()
	payload.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	payload.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	payload.Password = &hashedPassword

	defer cancel()

	//TODO :reduce token lifetime
	token, refreshToken, err := utils.GenerateNewTokenAndRefreshToken(*payload.Email, payload.User_id)

	if err != nil {
		log.Println("Unable to generate token")
		log.Panic(err.Error())
	}

	otp := utils.GenerateRandomOTP(6)

	email_body := fmt.Sprintln("<html><head></head><body><p>Hello,</p>Your verification code is " + otp + ".</p></body></html>")

	err = services.SendEmailV2("Verify email address", *payload.Email, email_body)

	if err != nil {
		log.Print(err.Error())
		log.Panic(err.Error())
	}

	payload.Token = &token
	payload.Refresh_token = &refreshToken
	payload.Generated_SignUp_code = &otp

	response, err := authDumpCollection.InsertOne(ctx, payload)

	if err != nil {
		log.Println("Unable to save record")
		log.Panic(err.Error())
	}

	log.Println(response.InsertedID)

	defer cancel()

	c.JSON(http.StatusOK, SignUpResponse{
		GenericResponse: generics.GenericResponse{
			Status:  "Success",
			Message: "An OTP was sent to your email address, you need to verify your email address before you can use our service",
		},
		Token: *payload.Token,
	})

}

func VerifyEmail(c *gin.Context) {
	defer utils.HandlePanic(c)

	var user models.User
	var payload models.VerifyEmailPayload

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	email := c.GetString("email")
	uid := c.GetString("uid")

	if email == "" || uid == "" {
		log.Panic("Service unavailable")
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadGateway, generics.GenericResponse{
			Status:  "Failed",
			Message: "OTP is required",
		})
		return
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})

	defer cancel()

	if err != nil {
		log.Panic(err.Error())
	}

	if count > 0 {
		c.JSON(http.StatusBadGateway, generics.GenericResponse{
			Status:  "Failed",
			Message: "Email is in use, please sign up with another email",
		})
		return
	}

	err = authDumpCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	defer cancel()

	if err != nil {
		log.Println(err.Error())
		log.Panic("Unable to retrieve data")
	}

	if *user.Generated_SignUp_code != *payload.Otp {

		c.JSON(http.StatusUnauthorized, generics.GenericResponse{
			Status:  "Failed",
			Message: "invalid code",
		})
		return
	}

	token, refreshToken, err := utils.GenerateNewTokenAndRefreshToken(email, uid)

	if err != nil {
		log.Panic("Unable to generate credentials")
	}

	defer cancel()

	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	user.Verified = true
	user.Token = &token
	user.Refresh_token = &refreshToken

	response, err := userCollection.InsertOne(ctx, user)

	defer cancel()

	if err != nil {
		log.Panic("Unable to save data, please try again later")
	}

	result, err := authDumpCollection.DeleteOne(ctx, bson.M{"email": email})

	if err != nil {
		log.Println(err.Error())
	}

	log.Println(result.DeletedCount)

	email_body := fmt.Sprintln("<html><head></head><body><p>Hello,</p>Hi, you have successfully verified your email on zate app.</p></body></html>")

	err = services.SendEmailV2("Verification Successful", email, email_body)

	if err != nil {
		log.Println("Unable to send success message")
	}

	defer cancel()

	c.JSON(http.StatusOK, VerifyEmailResponse{
		GenericResponse: generics.GenericResponse{
			Status:  "success",
			Message: "Email verified successfully",
		},
		Token:         token,
		Refresh_token: refreshToken,
		User_id:       response.InsertedID,
	})

}
