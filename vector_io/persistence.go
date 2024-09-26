package vector_io

// *** Potentially move this code to vector_space.go ***

import (
	"bytes"
	"eigen_db/constants"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"eigen_db/cfg"
)

func (store *vectorStore) persistToDisk(db_persist_path string) error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(store)
	if err != nil {
		return err
	}
	serializedData := buf.Bytes()

	return os.WriteFile(db_persist_path, serializedData, constants.DB_PERSIST_CHMOD)
}

func (store *vectorStore) loadPersistedVectors(db_persist_path string) error {
	serializedVectors, err := os.ReadFile(db_persist_path)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(serializedVectors)
	decoder := gob.NewDecoder(buf)
	err = decoder.Decode(store)
	if err != nil {
		return err
	}

	for id, v := range store.StoredVectors { // load deserialized stored vectors into the vector space
		err := store.vectorSpace.InsertVector(v.Components, uint32(id))
		if err != nil {
			return err // should probably panic since vectors are not properly loaded into memory
		}
	}

	return nil
}

func StartPersistenceLoop(config *cfg.Config) error {
	if _, err := os.Stat(constants.DB_PERSIST_PATH); os.IsNotExist(err) {
		if err = os.MkdirAll(constants.EIGEN_DIR, constants.DB_PERSIST_CHMOD); err != nil { // perm should maybe be switched to 600 instead of 400
			return err
		}
	}

	go func() {
		for {
			err := vectorStoreInstance.persistToDisk(constants.DB_PERSIST_PATH)
			if err != nil {
				fmt.Printf("Failed to persist data to disk: %s\n", err)
			}

			time.Sleep((&cfg.ConfigFactory{}).GetConfig().GetPersistenceTimeInterval())
		}
	}()

	return nil
}
