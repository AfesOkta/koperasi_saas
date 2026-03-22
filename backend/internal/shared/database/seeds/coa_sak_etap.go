package seeds

import (
	"context"
	"log"

	accountingModel "github.com/koperasi-gresik/backend/internal/modules/accounting/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SeedCOASAKETAP seeds the Chart of Accounts based on SAK ETAP (Standar Akuntansi Keuangan untuk Entitas Tanpa Akuntabilitas Publik)
// which is the standard for Indonesian cooperatives.
func SeedCOASAKETAP(ctx context.Context, db *gorm.DB, orgID uint) error {
	accounts := []accountingModel.Account{
		// -- ASSETS (1000) --
		{Code: "1000", Name: "AKTIVA", Type: "Asset", NormalBalance: "debit"},
		{Code: "1100", Name: "AKTIVA LANCAR", Type: "Asset", NormalBalance: "debit"},
		{Code: "1110", Name: "Kas", Type: "Asset", NormalBalance: "debit"},
		{Code: "1120", Name: "Bank", Type: "Asset", NormalBalance: "debit"},
		{Code: "1150", Name: "Piutang Pinjaman Anggota", Type: "Asset", NormalBalance: "debit"},
		{Code: "1151", Name: "Penyisihan Penghapusan Piutang", Type: "Asset", NormalBalance: "credit"},
		{Code: "1200", Name: "PERSEDIAAN", Type: "Asset", NormalBalance: "debit"},
		{Code: "1300", Name: "AKTIVA TETAP", Type: "Asset", NormalBalance: "debit"},
		{Code: "1310", Name: "Tanah & Bangunan", Type: "Asset", NormalBalance: "debit"},
		{Code: "1311", Name: "Akumulasi Penyusutan Bangunan", Type: "Asset", NormalBalance: "credit"},
		{Code: "1320", Name: "Inventaris & Peralatan Kantor", Type: "Asset", NormalBalance: "debit"},
		{Code: "1321", Name: "Akumulasi Penyusutan Inventaris", Type: "Asset", NormalBalance: "credit"},

		// -- LIABILITIES (2000) --
		{Code: "2000", Name: "KEWAJIBAN", Type: "Liability", NormalBalance: "credit"},
		{Code: "2100", Name: "KEWAJIBAN JANGKA PENDEK", Type: "Liability", NormalBalance: "credit"},
		{Code: "2110", Name: "Simpanan Sukarela", Type: "Liability", NormalBalance: "credit"},
		{Code: "2120", Name: "Titipan Dana Anggota", Type: "Liability", NormalBalance: "credit"},
		{Code: "2130", Name: "Hutang Usaha", Type: "Liability", NormalBalance: "credit"},
		{Code: "2140", Name: "Beban yang Masih Harus Dibayar", Type: "Liability", NormalBalance: "credit"},

		// -- EQUITY (3000) --
		{Code: "3000", Name: "EKUITAS", Type: "Equity", NormalBalance: "credit"},
		{Code: "3100", Name: "Simpanan Pokok", Type: "Equity", NormalBalance: "credit"},
		{Code: "3110", Name: "Simpanan Wajib", Type: "Equity", NormalBalance: "credit"},
		{Code: "3200", Name: "Dana Cadangan", Type: "Equity", NormalBalance: "credit"},
		{Code: "3300", Name: "Hibat / Donasi", Type: "Equity", NormalBalance: "credit"},
		{Code: "3900", Name: "Sisa Hasil Usaha (SHU)", Type: "Equity", NormalBalance: "credit"},

		// -- REVENUE (4000) --
		{Code: "4000", Name: "PENDAPATAN", Type: "Revenue", NormalBalance: "credit"},
		{Code: "4100", Name: "Pendapatan Bunga Pinjaman", Type: "Revenue", NormalBalance: "credit"},
		{Code: "4200", Name: "Pendapatan Provisi & Administrasi", Type: "Revenue", NormalBalance: "credit"},
		{Code: "4300", Name: "Pendapatan Partisipasi Anggota", Type: "Revenue", NormalBalance: "credit"},
		{Code: "4400", Name: "Pendapatan Penjualan Toko / Unit Usaha", Type: "Revenue", NormalBalance: "credit"},

		// -- EXPENSES (5000) --
		{Code: "5000", Name: "BEBAN", Type: "Expense", NormalBalance: "debit"},
		{Code: "5100", Name: "Beban Operasional", Type: "Expense", NormalBalance: "debit"},
		{Code: "5110", Name: "Beban Gaji & Tunjangan", Type: "Expense", NormalBalance: "debit"},
		{Code: "5120", Name: "Beban Listrik, Air & Telpon", Type: "Expense", NormalBalance: "debit"},
		{Code: "5130", Name: "Beban ATK & Cetakan", Type: "Expense", NormalBalance: "debit"},
		{Code: "5200", Name: "Beban Bunga Simpanan", Type: "Expense", NormalBalance: "debit"},
		{Code: "5300", Name: "Beban Penyusutan Aktiva Tetap", Type: "Expense", NormalBalance: "debit"},
	}

	for i := range accounts {
		accounts[i].OrganizationID = orgID
		accounts[i].IsActive = true
	}

	if err := db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "organization_id"}, {Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "type", "normal_balance", "is_active"}),
		}).
		Create(&accounts).Error; err != nil {
		return err
	}

	log.Printf("✅ Seeded %d SAK ETAP accounts for org %d", len(accounts), orgID)
	return nil
}
