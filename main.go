package main

import (
	"context"
	"encoding/json"
	"fmt"

	elastic "gopkg.in/olivere/elastic.v5"
)

type Mobile struct {
	Name    string `json:"name"`
	Camera  string `json:"camera"`
	Storage string
	Battery string
}

const s = `
curl -XPUT http://localhost:9200/products?pretty
curl -XPUT http://localhost:9200/products?pretty
curl -XPUT http://localhost:9200/customers?pretty
curl -XPUT http://localhost:9200/orders?pretty

curl -XPUT http://localhost:9200/products/mobiles/1?pretty -d'
> { "name":"iphone 7", "camera":"12mp", "storage":"256gb", "battery":"included" }
> '

curl -XPUT http://localhost:9200/products/mobiles/2?pretty -d'
{ "name":"samsung galaxy", "camera":"8mp", "storage":"256gb", "battery":"included" }
'

curl -XPUT http://localhost:9200/products/laptops/1?pretty -d'
{ "name":"macbook pro", "camera":"8mp", "storage":"500gb", "battery":"included", "os":"el capitain" }
'

curl -XPUT http://localhost:9200/products/laptops/1?pretty -d'
{ "name":"macbook pro", "camera":"8mp", "storage":"500gb", "battery":"included", "os":"el capitain" }
'

curl -XGET 'http://localhost:9200/products/mobiles/1?pretty&source=name,camera'

curl -XPOST http://localhost:9200/products/laptops/1/_update?pretty -d'
{ "doc":{"color":"black"} }
'

curl -XGET http://localhost:9200/products/laptops/1?pretty
`
const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"tweet":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"image":{
					"type":"keyword"
				},
				"created":{
					"type":"date"
				},
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`

func main() {
	ctx := context.Background()
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}
	exists, err := client.IndexExists("twitter").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("mobilos").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	p := client.Index().Index("mobilos")
	iphone := Mobile{}
	json.Unmarshal([]byte(
		`{ "name":"iphone 7", "camera":"12mp", "storage":"256gb", "battery":"included" }`),
		&iphone)

	fmt.Println("Hej Verden", iphone, p)
}
