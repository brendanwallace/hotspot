# Heterogeneous Risk Taking in SIR dynamics

_Hotspot model shows how location-based superspreading accelerates and reshapes epidemics. Wallace, B., Dimitrov, D., Hébert-Dufresne, L., Berdahl, A._. 

## Contents

### src
The `src` folder contains code for running models.

To run, install `golang` from [go.dev](go.dev), and then run from the command line with:

```
go run main.go
```

This runs a series of simulations, whose parameters should be adjusted by
modifying the file `main.go`, and saves the output to a .json file in the folder
`data`.

The go package `simulate` can be configured to run an ABM simulation,
a deterministic integro-differential-equation model, and a _difference_ equation
model.

Unit tests can be run with

```
go test
```

### figures

The `figures` folder contains notebooks and python files for generating the
figures used for the paper.

In general, these load .json files produced by the `src` files. There is one
notebook to produce each figure in the main paper.
