package main

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Periode struct {
	Id         int `json:"id"`
	TahunAwal  int `json:"tahun_awal"`
	TahunAkhir int `json:"tahun_akhir"`
}

type Opd struct {
	KodeOpd string `json:"kode_opd"`
	NamaOpd string `json:"nama_opd"`
}

type User struct {
	Nama     string `json:"nama_pegawai"`
	Nip      string `json:"nip"`
	Email    string `json:"email"`
	IsActive int    `json:"-"`
	Status   string `json:"status"`
	Roles    string `json:"roles"`
}
