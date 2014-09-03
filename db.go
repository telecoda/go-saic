package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/HouzuoGuo/tiedot/db"
)

var imagesDBDir = "/tmp/go-saic/db"
var imagesDB db.DB

var discoveredCollectionName = "DiscoveredImages"
var thumbnailCollectionName = "ThumbnailImages"

var discoveredImagesColl *collection
var thumbnailImagesColl *collection

type collection struct {
	name  string
	dbCol *db.Col
}

func initDB() {
	// ****************** Collection Management ******************

	//os.RemoveAll(myDBDir)
	//defer os.RemoveAll(myDBDir)

	// (Create if not exist) open a database
	imagesDB, err := db.OpenDB(imagesDBDir)
	if err != nil {
		panic(err)
	}

	if optClearDB {
		// drop collections
		fmt.Println("Info: clearing existing db content")
		if err := imagesDB.Drop(discoveredCollectionName); err != nil {
			fmt.Printf("Info: Dropping discoveredImagesColl - %s\n", err)

		}
		if err := imagesDB.Drop(thumbnailCollectionName); err != nil {
			fmt.Printf("Info: Dropping thumbnailImagesColl - %s\n", err)
		}

	}

	// Create two collections:
	if err := imagesDB.Create(discoveredCollectionName); err != nil {
		// collections already exists - ignore error
		fmt.Printf("Info: Creating discoveredImagesColl - %s\n", err)

	}
	if err := imagesDB.Create(thumbnailCollectionName); err != nil {
		// collection already exists - ignore error
		fmt.Printf("Info: Creating thumbnailImagesColl - %s\n", err)
	}

	if optScrubDB {
		// Scrub (repair and compact) collections
		fmt.Printf("Scrubbing: %s collection\n", discoveredCollectionName)
		if err := imagesDB.Scrub(discoveredCollectionName); err != nil {
			panic(err)
		}
		fmt.Printf("Scrubbing: %s collection\n", thumbnailCollectionName)
		if err := imagesDB.Scrub(thumbnailCollectionName); err != nil {
			panic(err)
		}
	}

	discoveredImagesColl = &collection{
		name:  discoveredCollectionName,
		dbCol: imagesDB.Use(discoveredCollectionName),
	}

	thumbnailImagesColl = &collection{
		name:  thumbnailCollectionName,
		dbCol: imagesDB.Use(thumbnailCollectionName),
	}

	// Create indexes
	if err := discoveredImagesColl.dbCol.Index([]string{"filepath"}); err != nil {
		// index already exists - ignore error
		fmt.Printf("Info: Building index on discoveredImagesColl - %s\n", err)
	}
	if err := thumbnailImagesColl.dbCol.Index([]string{"filepath"}); err != nil {
		// index already exists - ignore error
		fmt.Printf("Info: Building index on thumbnailImagesColl - %s\n", err)
	}

}

func listDB() {

	fmt.Println("DiscoveredImages")
	fmt.Println("================")

	discoveredImagesColl.dbCol.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		fmt.Println("Document", id, "is", string(docContent))
		return true  // move on to the next document OR
		return false // do not move on to the next document
	})

	fmt.Println("ThumbnailImages")
	fmt.Println("===============")

	thumbnailImagesColl.dbCol.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		fmt.Println("Document", id, "is", string(docContent))
		return true  // move on to the next document OR
		return false // do not move on to the next document
	})

}

func (c *collection) saveImages(images []ImageDetail) error {

	fmt.Println("Saving images to collection:")
	for _, imageDetail := range images {

		fmt.Printf("Info: Saving image - %s\n", imageDetail.FilePath)
		docId, err := c.saveImage(imageDetail)

		if err != nil {

			fmt.Printf("Info: %s\n", err)

		} else {

			fmt.Printf("Saved doc: %d for image: %s\n", docId, imageDetail.Filename)

		}
	}

	return nil
}

func (c *collection) saveImage(image ImageDetail) (int, error) {

	// check if image already exists for filepath
	imageFound := c.findImageByPath(image.FilePath)
	if imageFound != nil {
		return 0, errors.New("Image already exists.")
	}
	// Insert document (afterwards the docID uniquely identifies the document and will never change)
	return c.dbCol.Insert(map[string]interface{}{
		"id":       image.Id,
		"filename": image.Filename,
		"filepath": image.FilePath})
}

func (c *collection) findImageByPath(filePath string) *ImageDetail {
	var query interface{}
	queryString := fmt.Sprintf(`{"eq": "%s", "in": ["filepath"]}`, filePath)

	json.Unmarshal([]byte(queryString), &query)

	queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

	if err := db.EvalQuery(query, c.dbCol, &queryResult); err != nil {
		panic(err)
	}

	// Query result are document IDs
	for id := range queryResult {
		// To get query result document, simply read it
		result, err := c.dbCol.Read(id)
		if err != nil {
			panic(err)
		}
		imageDetail := &ImageDetail{
			Id:       result["id"].(string),
			FilePath: result["filepath"].(string),
			Filename: result["filename"].(string),
		}
		return imageDetail
	}

	return nil
}

/*
In embedded usage, you are encouraged to use all public functions concurrently.
However please do not use public functions in "data" package by yourself - you most likely will not need to use them directly.

To compile and run the example:
    go build && ./tiedot -mode=example

It may require as much as 1.5GB of free disk space in order to run the example.
*/

func embeddedExample() {
	// ****************** Collection Management ******************

	myDBDir := "/tmp/MyDatabase"
	os.RemoveAll(myDBDir)
	defer os.RemoveAll(myDBDir)

	// (Create if not exist) open a database
	myDB, err := db.OpenDB(myDBDir)
	if err != nil {
		panic(err)
	}

	// Create two collections: Feeds and Votes
	if err := myDB.Create("Feeds"); err != nil {
		panic(err)
	}
	if err := myDB.Create("Votes"); err != nil {
		panic(err)
	}

	// What collections do I now have?
	for _, name := range myDB.AllCols() {
		fmt.Printf("I have a collection called %s\n", name)
	}

	// Rename collection "Votes" to "Points"
	if err := myDB.Rename("Votes", "Points"); err != nil {
		panic(err)
	}

	// Drop (delete) collection "Points"
	if err := myDB.Drop("Points"); err != nil {
		panic(err)
	}

	// Scrub (repair and compact) "Feeds"
	if err := myDB.Scrub("Feeds"); err != nil {
		panic(err)
	}

	// ****************** Document Management ******************

	// Start using a collection (the reference is valid until DB schema changes or Scrub is carried out)
	feeds := myDB.Use("Feeds")

	// Insert document (afterwards the docID uniquely identifies the document and will never change)
	docID, err := feeds.Insert(map[string]interface{}{
		"name": "Go 1.2 is released",
		"url":  "golang.org"})
	if err != nil {
		panic(err)
	}

	// Read document
	readBack, err := feeds.Read(docID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Document", docID, "is", readBack)

	// Update document
	err = feeds.Update(docID, map[string]interface{}{
		"name": "Go is very popular",
		"url":  "google.com"})
	if err != nil {
		panic(err)
	}

	// Process all documents (note that document order is undetermined)
	feeds.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		fmt.Println("Document", id, "is", string(docContent))
		return true  // move on to the next document OR
		return false // do not move on to the next document
	})

	// Delete document
	if err := feeds.Delete(docID); err != nil {
		panic(err)
	}

	// ****************** Index Management ******************
	// Indexes assist in many types of queries
	// Create index (path leads to document JSON attribute)
	if err := feeds.Index([]string{"author", "name", "first_name"}); err != nil {
		panic(err)
	}
	if err := feeds.Index([]string{"Title"}); err != nil {
		panic(err)
	}
	if err := feeds.Index([]string{"Source"}); err != nil {
		panic(err)
	}

	// What indexes do I have on collection A?
	for _, path := range feeds.AllIndexes() {
		fmt.Printf("I have an index on path %v\n", path)
	}

	// Remove index
	if err := feeds.Unindex([]string{"author", "name", "first_name"}); err != nil {
		panic(err)
	}

	// ****************** Queries ******************
	// Prepare some documents for the query
	feeds.Insert(map[string]interface{}{"Title": "New Go release", "Source": "golang.org", "Age": 3})
	feeds.Insert(map[string]interface{}{"Title": "Kitkat is here", "Source": "google.com", "Age": 2})
	feeds.Insert(map[string]interface{}{"Title": "Good Slackware", "Source": "slackware.com", "Age": 1})

	var query interface{}
	json.Unmarshal([]byte(`[{"eq": "New Go release", "in": ["Title"]}, {"eq": "slackware.com", "in": ["Source"]}]`), &query)

	queryResult := make(map[int]struct{}) // query result (document IDs) goes into map keys

	if err := db.EvalQuery(query, feeds, &queryResult); err != nil {
		panic(err)
	}

	// Query result are document IDs
	for id := range queryResult {
		// To get query result document, simply read it
		readBack, err := feeds.Read(id)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Query returned document %v\n", readBack)
	}

	// Gracefully close database
	if err := myDB.Close(); err != nil {
		panic(err)
	}
}
