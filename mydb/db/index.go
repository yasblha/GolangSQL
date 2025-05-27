package db

import (
	"encoding/binary"
	"errors"
	"io"
)

// IndexType représente le type d'index
type IndexType uint8

const (
	IndexTypeBTree IndexType = 2
)

// IndexHeader représente l'en-tête d'un index
type IndexHeader struct {
	Type        IndexType
	FirstPage   uint32
	LastPage    uint32
	KeySize     uint16
	PageSize    uint16
	MaxEntries  uint16
	MinEntries  uint16
}

// IndexNode représente un nœud dans l'arbre B
type IndexNode struct {
	IsLeaf     bool
	NumKeys    uint16
	Keys       []interface{}
	Values     []uint32 // Pointeurs vers les pages
	NextPage   uint32   // Pour les nœuds feuilles
}

// ReadIndexHeader lit l'en-tête d'un index
func ReadIndexHeader(file io.Reader) (*IndexHeader, error) {
	header := &IndexHeader{}
	
	// Lire le type d'index
	if err := binary.Read(file, binary.BigEndian, &header.Type); err != nil {
		return nil, err
	}

	// Vérifier que c'est bien un index B-tree
	if header.Type != IndexTypeBTree {
		return nil, errors.New("unsupported index type")
	}

	// Lire les autres champs
	if err := binary.Read(file, binary.BigEndian, &header.FirstPage); err != nil {
		return nil, err
	}
	if err := binary.Read(file, binary.BigEndian, &header.LastPage); err != nil {
		return nil, err
	}
	if err := binary.Read(file, binary.BigEndian, &header.KeySize); err != nil {
		return nil, err
	}
	if err := binary.Read(file, binary.BigEndian, &header.PageSize); err != nil {
		return nil, err
	}
	if err := binary.Read(file, binary.BigEndian, &header.MaxEntries); err != nil {
		return nil, err
	}
	if err := binary.Read(file, binary.BigEndian, &header.MinEntries); err != nil {
		return nil, err
	}

	return header, nil
}

// ReadIndexNode lit un nœud d'index
func ReadIndexNode(file io.Reader, header *IndexHeader) (*IndexNode, error) {
	node := &IndexNode{}
	
	// Lire le flag feuille
	var isLeaf uint8
	if err := binary.Read(file, binary.BigEndian, &isLeaf); err != nil {
		return nil, err
	}
	node.IsLeaf = isLeaf == 1

	// Lire le nombre de clés
	if err := binary.Read(file, binary.BigEndian, &node.NumKeys); err != nil {
		return nil, err
	}

	// TODO: Implémenter la lecture des clés et des valeurs
	// La lecture dépend du type de données stocké dans l'index

	return node, nil
} 