package main

type Response struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Periode `json:"data"`
}

type Periode struct {
	Id         int `json:"id"`
	TahunAwal  int `json:"tahun_awal"`
	TahunAkhir int `json:"tahun_akhir"`
}
