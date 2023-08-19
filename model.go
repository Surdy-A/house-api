package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type house struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Country     string  `json:"country"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Photo       string  `json:"photo"`
}

func getHouses(db *sql.DB) ([]house, error) {
	query := "SELECT id, name, address, country, description, price, photo from house"
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	houses := []house{}
	for rows.Next() {
		var h house
		err := rows.Scan(&h.ID, &h.Name, &h.Address, &h.Country, &h.Description, &h.Price, &h.Photo)
		if err != nil {
			return nil, err
		}
		houses = append(houses, h)
	}

	return houses, nil
}

func (h *house) getHouse(db *sql.DB) error {
	query := fmt.Sprintf("SELECT id, name, address, country, description, price, photo from house where id=%v", h.ID)
	row := db.QueryRow(query)
	err := row.Scan(&h.ID, &h.Name, &h.Address, &h.Country, &h.Description, &h.Price, &h.Photo)
	if err != nil {
		return err
	}
	return nil
}

func (h *house) createHouse(db *sql.DB) error {
	name := h.Name
	fmt.Println(name)
	query := fmt.Sprintf("insert into house(id, name, address, country, description, price, photo) values('%v', '%v', '%v', '%v', '%v', '%v', '%v')",
		h.ID, h.Name, h.Address, h.Country, h.Description, h.Price, h.Photo)

	result, err := db.Exec(query)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	h.ID = int(id)
	return nil

}

func (h *house) updateHouse(db *sql.DB) error {
	query := fmt.Sprintf("update house set name='%v', address='%v', country='%v', description='%v', price='%v', photo='%v' where id='%v'",
		h.Name, h.Address, h.Country, h.Description, h.Price, h.Photo, h.ID)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such row exists")
	}
	return err
}

func (h *house) deleteHouse(db *sql.DB) error {
	query := fmt.Sprintf("delete from house where id=%v", h.ID)
	_, err := db.Exec(query)
	return err
}
