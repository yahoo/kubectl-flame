//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package cmd

import "fmt"

func PrintSuccess() {
	fmt.Printf("✔\n")
}

func PrintError() {
	fmt.Printf("❌\n")
}
