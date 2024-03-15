package components

import (
	"fmt"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

func CheckIfAccountExists(accountName string, client cosmosclient.Client, addressPrefix string, accountPath string) (bool, string) {

	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	account, err := registry.GetByName(accountName)
	if err != nil {
		return false, ""
	}

	addr, err := account.Address(addressPrefix)
	if err != nil {
		fmt.Println("Failed to get the Address:", err)
		return false, ""
	}

	return true, addr
}

func FetchAccount(accountName string, client cosmosclient.Client, addressPrefix string, accountPath string) (account cosmosaccount.Account, addr string, err error) {
	isExist, _ := CheckIfAccountExists(accountName, client, addressPrefix, accountPath)
	if isExist {
		registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
		if err != nil {
			fmt.Println(err)
			return cosmosaccount.Account{}, "", err
		}

		account, err := registry.GetByName(accountName)
		if err != nil {
			return cosmosaccount.Account{}, "", err
		}

		addr, err := account.Address(addressPrefix)
		if err != nil {
			fmt.Println("Failed to get the Address:", err)
			return cosmosaccount.Account{}, "", err
		}

		return account, addr, nil
	} else {
		return cosmosaccount.Account{}, "", fmt.Errorf("account not found")
	}
}

func CreateAccount(accountName string, accountPath string) {
	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		fmt.Println(err)
		return
	}
	account, mnemonic, err := registry.Create(accountName)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _ = mnemonic, account
	accountAddr, err := account.Address("air")
	if err != nil {
		fmt.Println(err.Error())
	}

	Logger.Info(fmt.Sprintf("Rollup Account Created > %s", accountAddr))
	Logger.Debug(mnemonic)
}
