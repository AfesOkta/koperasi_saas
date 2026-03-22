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
			DoNothing: true,
		}).
		Create(&accounts).Error; err != nil {
		return err
	}
	// setelah proses insert data coa, saya ingin update parent_id dari account yang memiliki referensi code characted 1 sampai ke-3
	// menjadi parent_id account yang memiliki referensi code characted 1 sampai ke-2
	// contoh: account dengan code 1100 memiliki parent_id 1000
	// account dengan code 1110 memiliki parent_id 1100
	// account dengan code 1111 memiliki parent_id 1110
	// account dengan code 11111 memiliki parent_id 1111
	// account dengan code 111111 memiliki parent_id 11111
	// account dengan code 1111111 memiliki parent_id 111111

	// 1. Fetch all accounts for the current orgID to get their database IDs
	var allAccounts []accountingModel.Account
	if err := db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&allAccounts).Error; err != nil {
		return err
	}

	// 2. Build map of Code to ID
	codeToID := make(map[string]uint)
	for _, acc := range allAccounts {
		codeToID[acc.Code] = acc.ID
	}

	// 3. Update ParentIDs based on code structure
	for _, acc := range allAccounts {
		parentCode := ""
		if len(acc.Code) > 4 {
			parentCode = acc.Code[:len(acc.Code)-1]
		} else if len(acc.Code) == 4 {
			if acc.Code[3] != '0' {
				parentCode = acc.Code[:3] + "0"
			} else if acc.Code[2] != '0' {
				parentCode = acc.Code[:2] + "00"
			} else if acc.Code[1] != '0' {
				parentCode = acc.Code[:1] + "000"
			}
		}

		if parentCode != "" {
			if parentID, ok := codeToID[parentCode]; ok {
				if err := db.WithContext(ctx).Model(&acc).Update("parent_id", parentID).Error; err != nil {
					return err
				}
			}
		}
	}

	log.Printf("✅ Seeded and linked %d SAK ETAP accounts for org %d", len(accounts), orgID)
	return nil
}

// SeedAllCOA iterates through all organizations and runs SeedCOASAKETAP for each.
func SeedAllCOA(ctx context.Context, db *gorm.DB) error {
	var orgIDs []uint
	if err := db.WithContext(ctx).Table("organizations").Pluck("id", &orgIDs).Error; err != nil {
		return err
	}

	for _, id := range orgIDs {
		_ = SeedCOASAKETAP(ctx, db, id)
	}
	return nil
}
