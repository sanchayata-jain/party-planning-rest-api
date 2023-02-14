package tables

import (
	"context"
	"encoding/json"

	// "net/http"

	"gorm.io/gorm"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) GetTables(ctx context.Context) ([]byte, error) {
	tables := []*models.Table{}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Table{}).
			Find(&tables).
			WithContext(ctx).
			Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(tables)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r Repository) CreateTable(ctx context.Context, table models.Table) error {
	createTable := &models.Table{}
	// err := r.db.Transaction(func(tx *gorm.DB) error {
	// 	err := tx.Create(table).WithContext(ctx).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// return tx.First(createTable, "id = ?", table.ID).WithContext(ctx).Error
	// 	return nil
	// })
	// if err != nil {
	// 	return err
	// }

	// r.db.Raw("SELECT id, name, age FROM users WHERE name = ?", 3).Scan(&result)
	// tx := r.db.Raw("INSERT INTO tables (capacity, seats_empty) VALUES (30, 30);")
	// table := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	// if tx.Error != nil {
	// 	return tx.Error
	// }

	tx := r.db.Select("Capacity", "SeatsEmpty").Create(&table)
	if tx.Error != nil {
		return tx.Error
	}
	_, err := json.Marshal(*createTable)
	if err != nil {
		return err
	}

	// w.Write([]byte("New note created:\n"))
	// w.Write(b)
	return nil
}

// func (r Repository) AddTable() error {
// 	// r.db.Exec(`INSERT INTO %s ("capacity", "seats_empty")
// 	// VALUES (%d, %d);`, tablename, newTable.Capacity, newTable.SeatsEmpty)

// 	// sqlStatement := "INSERT INTO tables (capacity, seats_empty) VALUES (5, 5);"
// 	// _, err := r.db.Exec(sqlStatement)
// 	// if err != nil {
// 	// 	fmt.Println("error in the repo")
// 	// 	return err
// 	// }

// 	r.db.
// 	return nil
// }

// INSERT INTO table_name (column1, column2, column3, ...)
// VALUES (value1, value2, value3, ...);
