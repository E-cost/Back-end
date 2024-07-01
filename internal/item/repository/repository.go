package itemsRepository

import (
	"Ecost/internal/database/client/postgresql"
	"Ecost/internal/item/storage"
	"Ecost/pkg/logging"
	"strings"
)

type repository struct {
	logger *logging.Logger
	client postgresql.Client
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func NewRepository(logger *logging.Logger, client postgresql.Client) storage.ItemsRepository {
	return &repository{
		logger: logger,
		client: client,
	}
}
