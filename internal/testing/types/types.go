/*
Copyright 2024-present Volodymyr Konstanchuk and contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

// These are the internal types that are needed to test the framework.

type CustomString string

type AInterface interface {
	Method()
}

type AStruct struct {
	Value int
}

func (a *AStruct) Method() {}

type BStruct struct {
	AStruct *AStruct
}

type CStruct struct {
	*AStruct      `componego:"inject"`
	Value         int
	privateField  AInterface `componego:"inject"`
	PublicField1  *AStruct   `componego:"inject,otherValue"`
	PublicField2  *BStruct
	IncorrectTag1 *AStruct `componego:"INJECT"`
	IncorrectTag2 *AStruct `COMPONEGO:"inject"`
}

type DStruct struct {
	PublicField2 AInterface `componego:"inject"`
	PublicField1 *CStruct   `componego:"inject"`
}

func (c *CStruct) GetPrivateField() AInterface {
	return c.privateField
}

var _ AInterface = (*AStruct)(nil)
