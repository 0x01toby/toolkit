package rlp

type Marshaler interface {
	MarshalRLPTo(dst []byte) ([]byte, error)
	//MarshalRLPWith(a *)
}
