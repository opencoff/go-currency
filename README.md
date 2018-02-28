# go-currency - Currency implementation using big.Int

## What is it?
Currency implementation that uses atto-dollars (1e-18) as the basic 
unit. All currency is normalized to atto-dollars for computation.
Computation uses big.Int.

## Notes
- The test cases cover the parsing and output only; I assume
  big.Int covers the boring arithmetic tests.
