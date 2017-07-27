run: dataset insert

dataset:
	wget http://mlg.ucd.ie/files/datasets/bbcsport-fulltext.zip
	unzip bbcsport-fulltext.zip

insert:
	go run indexer.go --dataset=./bbcsport -es=http://localhost:9200 --textbox=http://localhost:8080