package main

import (
	"context"
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
	exists, err := client.IndexExists("mobilos").Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		createIndex, err := client.CreateIndex("mobilos").BodyString(mapping).Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			fmt.Println("Not acknowledged", createIndex)
		}
	}

	// p := client.Index().Index("mobilos")
	iphone := Mobile{"iphone 7", "12mp", "12mp", "included"}
	// json.Unmarshal([]byte(
	// 	`{ "name":"iphone 7", "camera":"12mp", "storage":"12mp", "battery":"included" }`),
	// 	&iphone)
	put1, err := client.Index().
		Index("mobilos").
		Type("phone").
		Id("1").
		BodyJson(iphone).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Hej Verden", put1)
}
