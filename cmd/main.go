package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	courseSdk "github.com/iGuessImaDev/go_course_sdk/course"
	userSdk "github.com/iGuessImaDev/go_course_sdk/user"
	"github.com/iGuessImaDev/gocourse_enrollment/internal/enrollment"
	"github.com/iGuessImaDev/gocourse_enrollment/pkg/bootstrap"
	"github.com/iGuessImaDev/gocourse_enrollment/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	pagLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pagLimDef == "" {
		l.Fatal("paginator limit default is required")
	}

	courseTrans := courseSdk.NewHTTPClient(os.Getenv("API_COURSE_URL"), "")
	userTrans := userSdk.NewHTTPClient(os.Getenv("API_USER_URL"), "")

	ctx := context.Background()
	enrollmentRepo := enrollment.NewRepo(db, l)
	enrollmentSrv := enrollment.NewService(l, userTrans, courseTrans, enrollmentRepo)

	h := handler.NewEnrollmentHTTPServer(ctx, enrollment.MakeEndpoints(enrollmentSrv, enrollment.Config{LimPageDef: pagLimDef}))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		l.Println("listen in ", address)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		log.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
