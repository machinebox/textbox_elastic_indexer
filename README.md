# Textbox to increase relevance in your search

This is an example of use case of textbox to extract keywords, places and people from a news article, so I can increase relevance, on any search that I perform and power visualizations.

![Dashboard](plots/dashboard.png)

## Requirements

* Textbox  (https://machinebox.io/docs/textbox) 
* Elastic Search and Kibana (https://www.elastic.co)

## Dataset 

We are going to use a subset of BBCSport News Dataset that could be found here.

    http://mlg.ucd.ie/datasets/bbc.html
    http://mlg.ucd.ie/files/datasets/bbcsport-fulltext.zip


## Index the data

`indexer.go` will pre-process the articles using `textbox` and index the dataset into Elastic Search

You can download and run the inserting using the `Makefile`:

```
$ make run
```

Alternative you can do it manually:

```
# Get the dataset
$ wget http://mlg.ucd.ie/files/datasets/bbcsport-fulltext.zip
$ unzip bbcsport-fulltext.zip

# Run the indexer
$ go run indexer.go
```

# Structure of the document

Pre-Processing with textbox allows to have more structured document.

```
{
   id: "123",
   title: "Radcliffe will compete in London",
   content: "Paula Radcliffe will compete in the Flora London Marathon...",
   keywords: ["race director david bedford", "25th anniversary", "..."],
   places: ["London"],
   people: ["Paula Radcliffe"]
}
```

# Power a better search

Now that we have more structure data we can for example query by `place`:

```
GET news_textbox/_search
{
  "query": {
    "term": {
      "places.keyword": "London"
    }
  }
}
```

Or by people:
```
GET news_textbox/_search
{
  "query": {
    "term": {
      "people.keyword": "Paula Radcliffe"
    }
  }
}
```

# Visualize with Kibana

And you can plot visualizations to get trends and quick feedback.

## Tag clouds 
### People tag cloud
![People](plots/people.png)
### Keywords tag cloud
![Keywords](plots/keywords.png)
### Dashboard
![Dashboard](plots/dashboard.png)
