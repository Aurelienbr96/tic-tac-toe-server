// interfaces/connection.go
package interfaces

type Connection interface {
	ReadMessage() ([]byte, error)
	WriteMessage([]byte) error
	Close() error
	GetRemoteAddress() string
}
