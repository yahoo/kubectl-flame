//: Copyright Verizon Media
//: Licensed under the terms of the Apache 2.0 License. See LICENSE file in the project root for terms.
package cmd

import "fmt"

type Printer interface {
	Print(str string)
	PrintSuccess()
	PrintError()
}

type dryRunPrinter struct {
	dryRun bool
}

func (p *dryRunPrinter) Print(str string) {
	if !p.dryRun {
		fmt.Print(str)
	}
}

func (p *dryRunPrinter) PrintSuccess() {
	if !p.dryRun {
		fmt.Printf("✔\n")
	}
}

func (p *dryRunPrinter) PrintError() {
	fmt.Printf("❌\n")
}

func NewPrinter(dryRun bool) Printer {
	return &dryRunPrinter{
		dryRun: dryRun,
	}
}
