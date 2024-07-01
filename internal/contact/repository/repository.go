package contactsRepository

import (
	"Ecost/internal/contact/storage"
	"Ecost/internal/database/client/postgresql"
	"Ecost/pkg/logging"
	"strings"

	"github.com/jackc/pgconn"
)

type repository struct {
	logger *logging.Logger
	client postgresql.Client
}

var (
	pgErr *pgconn.PgError
)

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func NewRepository(logger *logging.Logger, client postgresql.Client) storage.ContactsRepository {
	return &repository{
		logger: logger,
		client: client,
	}
}
