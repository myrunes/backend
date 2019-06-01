package database

// Middleware describes the structure of a
// database middleware.
type Middleware interface {
	// Connect to the database server or file or
	// whatever you are about to use.
	Connect(params ...interface{}) error
	// Close the connection to the database.
	Close()
}
