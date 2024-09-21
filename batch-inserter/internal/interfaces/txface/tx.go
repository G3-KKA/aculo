package txface

type (
	API any
	// # [API] provided by Tx() should never change while at least one Tx is open
	//
	// Restrict access to state consistency dependent API via Tx()
	Tx[T API] interface {
		// # T provided by Tx() should mutate it's own state thread-safe
		//
		// # Or do not mutate it at all while at least one Tx is open
		//
		// Operation(s) happening previous to any specific api calls
		Tx() (T, Commit, error)
	}

	//
	// If API accessed via Tx() -- it is clientside responsibility to [Commit] the transaction
	Commit func() error
)
