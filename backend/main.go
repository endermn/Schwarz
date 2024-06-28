package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/stoyan-kukev/team-project/backend/middleware"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage:", os.Args[0], "<products.csv> <store.csv>")
		os.Exit(1)
	}
	productsCSVPath, storeCSVPath := os.Args[1], os.Args[2]

	box, err := objectbox.NewBuilder().Model(ObjectBoxModel()).Build()
	if err != nil {
		log.Fatal(err)
	}
	defer box.Close()

	userBox := BoxForuser(box)
	productBox := BoxForproduct(box)
	storeBox := BoxForstore(box)

	makeDefaultUser(userBox)
	loadProducts(productBox, productsCSVPath)
	defaultStoreID := loadDefaultStore(storeBox, storeCSVPath)

	mux := http.NewServeMux()
	registerMainEndpoints(mux, userBox, productBox, storeBox, defaultStoreID)

	handler := middleware.EnableCORS(mux)
	server := &http.Server{
		Addr:    ":12345",
		Handler: handler,
	}

	terminationChan := make(chan os.Signal, 1)
	signal.Notify(terminationChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-terminationChan
		server.Shutdown(context.Background())
	}()

	go runAuth(userBox)

	log.Println("Server started")
	err = server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
