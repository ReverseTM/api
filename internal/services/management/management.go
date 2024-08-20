package management

import (
	"fmt"
	"log/slog"
)

type DataReader interface {
	Data(keys []string) (map[string]any, error)
}

type DataWriter interface {
	SaveData(data map[string]any) error
}

type ManagementService struct {
	log        *slog.Logger
	dataReader DataReader
	dataWriter DataWriter
}

func NewManagementService(
	log *slog.Logger,
	dataReader DataReader,
	dataWriter DataWriter,
) *ManagementService {
	return &ManagementService{
		log:        log,
		dataReader: dataReader,
		dataWriter: dataWriter,
	}
}

func (m *ManagementService) Read(keys []string) (map[string]any, error) {
	const op = "services.management.Read"

	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("reading data")

	result, err := m.dataReader.Data(keys)
	if err != nil {
		log.Error("failed reading data")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (m *ManagementService) Write(data map[string]any) error {
	const op = "services.management.Write"

	log := m.log.With(
		slog.String("op", op),
	)

	log.Info("writing data")

	err := m.dataWriter.SaveData(data)
	if err != nil {
		log.Error("failed writing data")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
