package tests

import (
	"encoding/binary"
	"mydb/internal/db"
	"mydb/pkg/sql"
	"os"
	"testing"
)

// Crée un fichier SQLite de test avec des données
func createTestDB(t *testing.T) string {
	// Créer un fichier temporaire
	tmpfile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpfile.Close()

	// Créer l'en-tête SQLite complet (100 bytes)
	header := make([]byte, 100)

	// 1. Magic string (16 bytes)
	copy(header[0:16], []byte("SQLite format 3\000"))

	// 2. Page size (2 bytes)
	binary.BigEndian.PutUint16(header[16:18], 4096)

	// 3. File format write version (1 byte)
	header[18] = 1

	// 4. File format read version (1 byte)
	header[19] = 1

	// 5. Bytes of unused "reserved" space at the end of each page (1 byte)
	header[20] = 0

	// 6. Maximum embedded payload fraction (1 byte)
	header[21] = 64

	// 7. Minimum embedded payload fraction (1 byte)
	header[22] = 32

	// 8. Leaf payload fraction (1 byte)
	header[23] = 32

	// 9. File change counter (4 bytes)
	binary.BigEndian.PutUint32(header[24:28], 0)

	// 10. Size of the database file in pages (4 bytes)
	binary.BigEndian.PutUint32(header[28:32], 1)

	// 11. Page number of the first freelist trunk page (4 bytes)
	binary.BigEndian.PutUint32(header[32:36], 0)

	// 12. Total number of freelist pages (4 bytes)
	binary.BigEndian.PutUint32(header[36:40], 0)

	// 13. The schema cookie (4 bytes)
	binary.BigEndian.PutUint32(header[40:44], 1)

	// 14. The schema format number (4 bytes)
	binary.BigEndian.PutUint32(header[44:48], 4)

	// 15. Default page cache size (4 bytes)
	binary.BigEndian.PutUint32(header[48:52], 0)

	// 16. The page number of the largest root b-tree page (4 bytes)
	binary.BigEndian.PutUint32(header[52:56], 0)

	// 17. The text encoding (4 bytes) - 1 for UTF-8
	binary.BigEndian.PutUint32(header[56:60], 1)

	// 18. The "user version" (4 bytes)
	binary.BigEndian.PutUint32(header[60:64], 0)

	// 19. True (non-zero) for incremental-vacuum mode (4 bytes)
	binary.BigEndian.PutUint32(header[64:68], 0)

	// 20. The "Application ID" (4 bytes)
	binary.BigEndian.PutUint32(header[68:72], 0)

	// 21. Reserved for expansion (20 bytes)
	// Les 20 derniers bytes sont réservés et mis à 0 par défaut

	// Écrire l'en-tête
	if _, err := tmpfile.Write(header); err != nil {
		t.Fatal(err)
	}

	// Écrire la table master (première page)
	masterPage := make([]byte, 4096)
	// Type = table
	masterPage[0] = 0x0D
	// Nombre de colonnes = 5
	masterPage[1] = 5
	// Offset de la première cellule = 100
	binary.BigEndian.PutUint16(masterPage[3:5], 100)

	// Écrire la définition de la table users
	usersDef := []byte("CREATE TABLE users(id INTEGER PRIMARY KEY, name TEXT)")
	copy(masterPage[100:], usersDef)

	if _, err := tmpfile.Write(masterPage); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name()
}

func TestBinaryOperations(t *testing.T) {
	// Créer la base de test
	dbPath := createTestDB(t)
	defer os.Remove(dbPath)

	// Ouvrir le fichier
	file, err := os.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Tester la lecture de l'en-tête
	info := db.ParseHeader(file)
	if info.PageSize != 4096 {
		t.Errorf("Page size attendu: 4096, obtenu: %d", info.PageSize)
	}

	// Tester la lecture de la table master
	tables := db.ReadMasterTable(file, info)
	if len(tables) == 0 {
		t.Error("Aucune table trouvée")
	}

	// Tester le parsing d'une requête SELECT
	query := "SELECT id, name FROM users WHERE id = 1"
	selectQuery, err := sql.Parse(query)
	if err != nil {
		t.Errorf("Erreur de parsing SELECT: %v", err)
	}
	if selectQuery.GetTable() != "users" {
		t.Errorf("Table attendue: users, obtenue: %s", selectQuery.GetTable())
	}

	// Tester le parsing d'une requête INSERT
	insertQuery := "INSERT INTO users (id, name) VALUES (1, 'John')"
	insert, err := sql.Parse(insertQuery)
	if err != nil {
		t.Errorf("Erreur de parsing INSERT: %v", err)
	}
	if insert.GetTable() != "users" {
		t.Errorf("Table attendue: users, obtenue: %s", insert.GetTable())
	}
}

// Test des opérations sur les données binaires
func TestBinaryDataOperations(t *testing.T) {
	dbPath := createTestDB(t)
	defer os.Remove(dbPath)

	file, err := os.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Tester l'insertion de données binaires
	// TODO: Implémenter l'insertion de données binaires
	// Pour l'instant, on vérifie juste que le fichier existe
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Le fichier de test n'existe pas")
	}
}
