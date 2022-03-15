package hdwallet

import "fmt"

// https://github.com/satoshilabs/slips/blob/master/slip-0044.md

// m/purpose'/coin_type'/account'/change/address_index

func FilPath(index int) string {
	return fmt.Sprintf("m/44'/461'/0'/0/%d", index)
}
