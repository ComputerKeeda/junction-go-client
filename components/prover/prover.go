package prover

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/components"
	"github.com/ComputerKeeda/junction-go-client/types"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"os"
)

type MyCircuit struct {
	To              [components.PodSize]frontend.Variable `gnark:",public"`
	From            [components.PodSize]frontend.Variable `gnark:",public"`
	Amount          [components.PodSize]frontend.Variable `gnark:",public"`
	TransactionHash [components.PodSize]frontend.Variable `gnark:",public"`
	FromBalances    [components.PodSize]frontend.Variable `gnark:",public"`
	ToBalances      [components.PodSize]frontend.Variable `gnark:",public"`
}

type TransactionSecond struct {
	To              string
	From            string
	Amount          string
	FromBalances    string
	ToBalances      string
	TransactionHash string
}

func getTransactionHash(tx TransactionSecond) string {
	record := tx.To + tx.From + tx.Amount + tx.FromBalances + tx.ToBalances + tx.TransactionHash
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}
func GetMerkleRootSecond(transactions []TransactionSecond) string {
	var merkleTree []string

	for _, tx := range transactions {
		merkleTree = append(merkleTree, getTransactionHash(tx))
	}

	for len(merkleTree) > 1 {
		var tempTree []string
		for i := 0; i < len(merkleTree); i += 2 {
			if i+1 == len(merkleTree) {
				tempTree = append(tempTree, merkleTree[i])
			} else {
				combinedHash := merkleTree[i] + merkleTree[i+1]
				h := sha256.New()
				h.Write([]byte(combinedHash))
				tempTree = append(tempTree, hex.EncodeToString(h.Sum(nil)))
			}
		}
		merkleTree = tempTree
	}

	return merkleTree[0]
}

func (circuit *MyCircuit) Define(api frontend.API) error {
	for i := 0; i < components.PodSize; i++ {
		api.AssertIsLessOrEqual(circuit.Amount[i], circuit.FromBalances[i])

		api.Sub(circuit.FromBalances[i], circuit.Amount[i])
		api.Add(circuit.ToBalances[i], circuit.Amount[i])

		updatedFromBalance := api.Sub(circuit.FromBalances[i], circuit.Amount[i])
		updatedToBalance := api.Add(circuit.ToBalances[i], circuit.Amount[i])

		api.AssertIsEqual(updatedFromBalance, api.Sub(circuit.FromBalances[i], circuit.Amount[i]))
		api.AssertIsEqual(updatedToBalance, api.Add(circuit.ToBalances[i], circuit.Amount[i]))
	}

	return nil
}

func ComputeCCS() constraint.ConstraintSystem {
	var circuit MyCircuit
	ccs, _ := frontend.Compile(ecc.BLS12_381.ScalarField(), r1cs.NewBuilder, &circuit)

	return ccs
}

func GenerateKeyPair() (groth16.ProvingKey, groth16.VerifyingKey, error) {
	ccs := ComputeCCS()
	pk, vk, error := groth16.Setup(ccs)
	return pk, vk, error
}

func GenerateProof(inputData types.PodStruct, PodNumber int) (any, string, []byte, error) {
	fmt.Println("Generating Proof")
	ccs := ComputeCCS()

	var transactions []TransactionSecond
	for i := 0; i < components.PodSize; i++ {
		transaction := TransactionSecond{
			To:              inputData.To[i],
			From:            inputData.From[i],
			Amount:          inputData.Amounts[i],
			FromBalances:    inputData.SenderBalances[i],
			ToBalances:      inputData.ReceiverBalances[i],
			TransactionHash: inputData.TransactionHash[i],
		}
		transactions = append(transactions, transaction)
	}
	currentStatusHash := GetMerkleRootSecond(transactions)
	fmt.Println("currentStatusHash Merkle root")
	fmt.Println(currentStatusHash)
	if _, err := os.Stat("provingKey.txt"); os.IsNotExist(err) {
		fmt.Println("Proving key does not exist. Please run the command 'sequencer-sdk create-vk-pk' to generate the proving key")
		return nil, "", nil, err
	}

	pk, err := ReadProvingKeyFromFile("provingKey.txt")

	if err != nil {
		fmt.Println("Error reading proving key:", err)
		return nil, "", nil, err
	}

	if err != nil {
		fmt.Println("Error getting snark field")
		return nil, "", nil, err
	}
	var inputValueLength int

	fromLength := len(inputData.From)
	toLength := len(inputData.To)
	amountsLength := len(inputData.Amounts)
	txHashLength := len(inputData.TransactionHash)
	senderBalancesLength := len(inputData.SenderBalances)
	receiverBalancesLength := len(inputData.ReceiverBalances)
	messagesLength := len(inputData.Messages)
	txNoncesLength := len(inputData.TransactionNonces)
	accountNoncesLength := len(inputData.AccountNonces)

	if fromLength == toLength &&
		fromLength == amountsLength &&
		fromLength == txHashLength &&
		fromLength == senderBalancesLength &&
		fromLength == receiverBalancesLength &&
		fromLength == messagesLength &&
		fromLength == txNoncesLength &&
		fromLength == accountNoncesLength {
		inputValueLength = fromLength
	} else {
		fmt.Println("Error: Input data is not correct")
		return nil, "", nil, fmt.Errorf("input data is not correct")
	}

	if inputValueLength < components.PodSize {
		leftOver := components.PodSize - inputValueLength
		for i := 0; i < leftOver; i++ {
			inputData.From = append(inputData.From, "0")
			inputData.To = append(inputData.To, "0")
			inputData.Amounts = append(inputData.Amounts, "0")
			inputData.TransactionHash = append(inputData.TransactionHash, "0")
			inputData.SenderBalances = append(inputData.SenderBalances, "0")
			inputData.ReceiverBalances = append(inputData.ReceiverBalances, "0")
			inputData.Messages = append(inputData.Messages, "0")
			inputData.TransactionNonces = append(inputData.TransactionNonces, "0")
			inputData.AccountNonces = append(inputData.AccountNonces, "0")
		}
	}

	inputs := MyCircuit{
		To:              [components.PodSize]frontend.Variable{},
		From:            [components.PodSize]frontend.Variable{},
		Amount:          [components.PodSize]frontend.Variable{},
		TransactionHash: [components.PodSize]frontend.Variable{},
		FromBalances:    [components.PodSize]frontend.Variable{},
		ToBalances:      [components.PodSize]frontend.Variable{},
	}

	for i := 0; i < components.PodSize; i++ {
		// var amount string
		if inputData.Amounts[i] > inputData.SenderBalances[i] {
			fmt.Println("Amount value give below")
			fmt.Println(inputData.Amounts[i])
			fmt.Println("From balance value given below")
			fmt.Println(inputData.SenderBalances[i])
			os.Exit(90)
		}
		inputs.To[i] = frontend.Variable(inputData.To[i])
		inputs.From[i] = frontend.Variable(inputData.From[i])
		inputs.Amount[i] = frontend.Variable(inputData.Amounts[i])
		inputs.TransactionHash[i] = frontend.Variable(inputData.TransactionHash[i])
		inputs.FromBalances[i] = frontend.Variable(inputData.SenderBalances[i])
		inputs.ToBalances[i] = frontend.Variable(inputData.ReceiverBalances[i])
	}

	witness, err := frontend.NewWitness(&inputs, ecc.BLS12_381.ScalarField())
	if err != nil {
		// fmt.Println(inputs.From)
		// fmt.Println(inputs.To)
		fmt.Println(inputs.Amount)
		fmt.Println(inputs.FromBalances)
		fmt.Println(inputs.ToBalances)

		fmt.Printf("Error creating a witness: %v\n", err)
		return nil, "", nil, err
	}

	witnessVector := witness.Vector()

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		fmt.Printf("Error generating proof: %v\n", err)
		return nil, "", nil, err
	}

	proofDbValue, err := json.Marshal(proof)
	if err != nil {
		fmt.Println("Error marshalling proof:", err)
		return nil, "", nil, err
	}

	return witnessVector, currentStatusHash, proofDbValue, nil
}

func ReadProvingKeyFromFile(filename string) (groth16.ProvingKey, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pk := groth16.NewProvingKey(ecc.BLS12_381)
	_, err = pk.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return pk, nil
}
