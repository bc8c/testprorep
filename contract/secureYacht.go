package main

import (
	"encoding/json"
	"fmt"	

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type YachtStruct struct {
	Name string `json:NAME`
	Color string `json:COLOR`
	Size int `json:SIZE`
	Owner string `json:OWNER`
	Price int `json:PRICE`	
}

type SecureYacht struct {
	contractapi.Contract
}

func (s *SecureYacht) YachtExists(ctx contractapi.TransactionContextInterface, name string) (bool, error) {
	// 원장에 name을 key로 갖는 데이터가 존재하면 true, 아니면 false 반환
	yachtbytes, err := ctx.GetStub().GetState(name)
	if err != nil {
		return false, err
	}
	return yachtbytes != nil, nil
}

 
func (s *SecureYacht) RegistYacht(ctx contractapi.TransactionContextInterface, username string, name string, color string, size int, price int) error {

	// name 으로 원장에 값이 존재하는 체크!
	exists, err := s.YachtExists(ctx, name)
	if err != nil {
		return err
	}
	if exists != false {
		return fmt.Errorf(`name(%s) is already exists`, name)
	}

	if username != "YachtMaker" {
		return fmt.Errorf(`username is not "YachtMaker"`)
	}
	
	yacht := YachtStruct{
		Name: name,
		Color: color,
		Size: size,
		Owner: "YachtMaker",
		Price: price,
	}

	yachtbytes, err := json.Marshal(yacht)

	if err != nil {
		return err
	}
	
	return ctx.GetStub().PutState(name, yachtbytes)
}

func (s *SecureYacht) ReadYacht(ctx contractapi.TransactionContextInterface, name string) (*YachtStruct, error) {
	yachtbytes, err := ctx.GetStub().GetState(name)
	if err != nil {
		return nil, err
	}
	if yachtbytes == nil {
		return nil, fmt.Errorf(`yacht(%s) is not exists`, name)
	}

	var yacht YachtStruct
	err = json.Unmarshal(yachtbytes, &yacht)
	if err != nil {
		return nil, err
	}

	return &yacht, nil
}

func (s *SecureYacht) DeleteYacht(ctx contractapi.TransactionContextInterface, name string) error {
	
	// name 으로 원장에 값이 존재하는 체크!
	exists, err := s.YachtExists(ctx, name)
	if err != nil {
		return err
	}
	if exists == false {
		return fmt.Errorf(`name(%s) does not exists`, name)
	}

	return ctx.GetStub().DelState(name)
}

func (s *SecureYacht) SellYacht(ctx contractapi.TransactionContextInterface, name string, fromuser string, touser string) error {
	// name으로 요트 조회
	yacht, err := s.ReadYacht(ctx, name)
	if err != nil {
		return err
	}
	// 요트의 owner 정보가 == fromuser 가 맞는지 체크
	if yacht.Owner != fromuser{
		return fmt.Errorf(`can not sell the yacht !!`)
	}	
	// 요트 정보의 owner를 touser -> 원장에 다시 덮어씌우기
	yacht.Owner = touser

	yachtbytes, err := json.Marshal(yacht)

	if err != nil {
		return err
	}	
	return ctx.GetStub().PutState(name, yachtbytes)
}

//  name string, color string, size int, price int => 하나의 구조체로 묶는다 =>  YachtStruct => JSON 마샬링 => 하나의 bytes 코드

func main() {

	chaincode, err := contractapi.NewChaincode(new(SecureYacht))

	if err != nil {
		fmt.Printf("Error create teamate chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting teamate chaincode: %s", err.Error())
	}
}

