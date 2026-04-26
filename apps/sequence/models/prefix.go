package models


//this table has all system prefix_sequences : for employees invoicess orders (prefix_sequences)
type PrefixSequence struct {
    ID           int    `f:"number, primary, auto"`
    Name         string `f:"text"`   // e.g., "Employee ID", "Sales Invoice"
    Prefix       string `f:"text"`   // e.g., "EMP-", "INV-"
    NextValue    int    `f:"number"` // The next number to be assigned
    Digits       int    `f:"number"` // Total length of the number part (for padding)
    Step         int    `f:"number"` // Usually 1, but some systems skip numbers
}