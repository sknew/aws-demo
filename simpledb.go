package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/simpledb"
	"io"
	"log"
	"net/http"
)

var mySession *session.Session
var simpleDB *simpledb.SimpleDB

func InitSimpleDB() {

	mySession = session.Must(session.NewSession())
	//simpleDB = simpledb.New(mySession, aws.NewConfig().WithRegion("us-east-1"))
	//simpleDB = simpledb.New(mySession, aws.NewConfig().WithRegion("us-west-1"))
	simpleDB = simpledb.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
}

func CreateDomainSimpleDB(domain string) {

	createDomainInput := simpledb.CreateDomainInput{DomainName: &domain}
	log.Println("createDomainInput", createDomainInput.GoString())

	createDomainOutput, err := simpleDB.CreateDomain(&createDomainInput)
	if err != nil {
		log.Println("createDomain err", err)
	}
	log.Println("createDomainOutput", createDomainOutput.GoString())
}

func ListDomainsSimpleDB() {

	listDomainsInput := simpledb.ListDomainsInput{}
	listDomainsOutput, err := simpleDB.ListDomains(&listDomainsInput)
	if err != nil {
		log.Println("listDomainsOutput err", err)
	}
	log.Println("listDomainsOutput", listDomainsOutput.GoString())
}

func PutPairSimpleDB(domain, key, value, time, server string, w http.ResponseWriter) {

	replace := true

	nameValue := "Value"
	replaceableAttributeValue := simpledb.ReplaceableAttribute{Name: &nameValue, Replace: &replace, Value: &value}
	log.Println("replaceableAttributeValue", replaceableAttributeValue.GoString())

	nameTime := "Time"
	replaceableAttributeTime := simpledb.ReplaceableAttribute{Name: &nameTime, Replace: &replace, Value: &time}
	log.Println("replaceableAttributeTime", replaceableAttributeTime.GoString())

	nameServer := "Server"
	replaceableAttributeServer := simpledb.ReplaceableAttribute{Name: &nameServer, Replace: &replace, Value: &server}
	log.Println("replaceableAttributeServer", replaceableAttributeServer.GoString())

	replaceableAttributes := make([]*simpledb.ReplaceableAttribute, 3)
	replaceableAttributes[0] = &replaceableAttributeValue
	replaceableAttributes[1] = &replaceableAttributeTime
	replaceableAttributes[2] = &replaceableAttributeServer
	//log.Println("replaceableAttributes[0]", replaceableAttributes[0].GoString())

	putAttributesInput := simpledb.PutAttributesInput{Attributes: replaceableAttributes, DomainName: &domain, ItemName: &key}
	log.Println("putAttributesInput", putAttributesInput.GoString())
	if w != nil {
		io.WriteString(w, "putAttributesInput"+putAttributesInput.GoString()+"\n")
	}

	putAttributesOutput, err := simpleDB.PutAttributes(&putAttributesInput)
	if err != nil {
		log.Println("PutAttributes err", err)
	}
	log.Println("putAttributesOutput", putAttributesOutput.GoString())
	if w != nil {
		io.WriteString(w, "putAttributesOutput"+putAttributesOutput.GoString()+"\n")
	}
}

func GetPairSimpleDB(domain, key string, w http.ResponseWriter) {

	attributeNames := make([]*string, 1)
	//attributeNames[0] = nil
	//log.Println("attributeNames[0]", *attributeNames[0])

	consistentRead := true
	getAttributesInput := simpledb.GetAttributesInput{AttributeNames: attributeNames, ConsistentRead: &consistentRead, DomainName: &domain, ItemName: &key}
	log.Println("getAttributesInput", getAttributesInput.GoString())
	if w != nil {
		io.WriteString(w, "getAttributesInput"+getAttributesInput.GoString()+"\n")
	}

	getAttributesOutput, err := simpleDB.GetAttributes(&getAttributesInput)
	if err != nil {
		log.Println("GetAttributes err", err)
	}
	log.Println("getAttributesOutput", getAttributesOutput.GoString())
	if w != nil {
		io.WriteString(w, "getAttributesOutput"+getAttributesOutput.GoString()+"\n")
	}
}

func DeletePairSimpleDB(domain, key string, w http.ResponseWriter) {

	deletableAttributes := make([]*simpledb.DeletableAttribute, 1)
	deletableAttributes[0] = nil
	//log.Println("deletableAttributes[0]", deletableAttributes[0].GoString())

	deleteAttributesInput := simpledb.DeleteAttributesInput{Attributes: deletableAttributes, DomainName: &domain, ItemName: &key}
	log.Println("deleteAttributesInput", deleteAttributesInput.GoString())
	if w != nil {
		io.WriteString(w, "deleteAttributesInput"+deleteAttributesInput.GoString()+"\n")
	}

	deleteAttributesOutput, err := simpleDB.DeleteAttributes(&deleteAttributesInput)
	if err != nil {
		log.Println("DeleteAttributes err", err)
	}
	log.Println("deleteAttributesOutput", deleteAttributesOutput.GoString())
	if w != nil {
		io.WriteString(w, "deleteAttributesOutput"+deleteAttributesOutput.GoString()+"\n")
	}
}

func SelectPairsSimpleDB(domain string, w http.ResponseWriter) {

	consistentRead := true
	selectExpression := "select * from " + domain
	log.Println("selectExpression", selectExpression)

	selectInput := simpledb.SelectInput{ConsistentRead: &consistentRead, SelectExpression: &selectExpression}
	log.Println("selectInput", selectInput.GoString())
	if w != nil {
		io.WriteString(w, "selectInput"+selectInput.GoString()+"\n")
	}

	selectOutput, err := simpleDB.Select(&selectInput)
	if err != nil {
		log.Println("Select err", err)
	}
	log.Println("selectOutput", selectOutput.GoString())
	if w != nil {
		io.WriteString(w, "selectOutput"+selectOutput.GoString()+"\n")
	}
}
