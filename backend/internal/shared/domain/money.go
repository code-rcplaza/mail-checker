package domain

// Money is an amount in hundredths, never a float: CLP 5.990 => 599000, USD 20,00 => 2000.
type Money int64
