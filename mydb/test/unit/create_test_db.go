package tests

import (
	"encoding/binary"
	"os"
)

// CreateTestDB crée un fichier SQLite de test complet
func CreateTestDB(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 1. Écrire l'en-tête SQLite (100 bytes)
	header := make([]byte, 100)
	
	// Magic string (16 bytes)
	copy(header[0:16], []byte("SQLite format 3\000"))
	
	// Page size (2 bytes) - 4096 bytes
	binary.BigEndian.PutUint16(header[16:18], 4096)
	
	// File format write version (1 byte)
	header[18] = 1
	
	// File format read version (1 byte)
	header[19] = 1
	
	// Reserved bytes (1 byte)
	header[20] = 0
	
	// Maximum embedded payload fraction (1 byte)
	header[21] = 64
	
	// Minimum embedded payload fraction (1 byte)
	header[22] = 32
	
	// Leaf payload fraction (1 byte)
	header[23] = 32
	
	// File change counter (4 bytes)
	binary.BigEndian.PutUint32(header[24:28], 1)
	
	// Size of the database file in pages (4 bytes)
	binary.BigEndian.PutUint32(header[28:32], 2) // 2 pages: header + table
	
	// Page number of the first freelist trunk page (4 bytes)
	binary.BigEndian.PutUint32(header[32:36], 0)
	
	// Total number of freelist pages (4 bytes)
	binary.BigEndian.PutUint32(header[36:40], 0)
	
	// Schema cookie (4 bytes)
	binary.BigEndian.PutUint32(header[40:44], 1)
	
	// Schema format number (4 bytes)
	binary.BigEndian.PutUint32(header[44:48], 4)
	
	// Default page cache size (4 bytes)
	binary.BigEndian.PutUint32(header[48:52], 0)
	
	// Largest root b-tree page (4 bytes)
	binary.BigEndian.PutUint32(header[52:56], 0)
	
	// Text encoding (4 bytes) - 1 for UTF-8
	binary.BigEndian.PutUint32(header[56:60], 1)
	
	// User version (4 bytes)
	binary.BigEndian.PutUint32(header[60:64], 0)
	
	// Incremental vacuum mode (4 bytes)
	binary.BigEndian.PutUint32(header[64:68], 0)
	
	// Application ID (4 bytes)
	binary.BigEndian.PutUint32(header[68:72], 0)
	
	// Reserved for expansion (20 bytes)
	// Les 20 derniers bytes sont réservés et mis à 0 par défaut

	if _, err := file.Write(header); err != nil {
		return err
	}

	// 2. Écrire la page de la table master (4096 bytes)
	masterPage := make([]byte, 4096)
	
	// En-tête de la page (8 bytes)
	// Type = table (0x0D)
	masterPage[0] = 0x0D
	// Nombre de colonnes = 5
	masterPage[1] = 5
	// Offset de la première cellule = 100
	binary.BigEndian.PutUint16(masterPage[3:5], 100)
	
	// Définition de la table users
	usersDef := []byte("CREATE TABLE users(id INTEGER PRIMARY KEY, name TEXT, age INTEGER, email TEXT, created_at DATETIME)")
	copy(masterPage[100:], usersDef)

	if _, err := file.Write(masterPage); err != nil {
		return err
	}

	return nil
} 