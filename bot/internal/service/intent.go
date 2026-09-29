package service

import (
	"strings"
	"unicode"
)

// HelpMessage is the reply for anything the bot doesn't handle: unknown
// commands, /help, and free text that can't be a transaction. Add new
// commands here as the bot grows so every entry point lists them.
const HelpMessage = `Aku bisa bantu kamu mencatat keuangan di CashBook.

Catat transaksi dengan pesan biasa, misalnya:
- "Beli bensin 50rb"
- "Gaji 5jt"

Perintah:
/summary - ringkasan pemasukan, pengeluaran, dan saldo
/report - laporan pemasukan & pengeluaran bulan ini
/budget - pemakaian anggaran bulan ini
/kategori - lihat daftar kategori kamu
/link <kode> - hubungkan akun CashBook
/start - cek akun yang terhubung
/help - tampilkan bantuan ini`

// amountWords are tokens that signal an amount written without digits
// ("lima puluh ribu") or with a unit suffix split off ("50 rb"). Matched as
// whole tokens so words like "kok" or "okay" don't count.
var amountWords = map[string]bool{
	"rp": true, "k": true, "rb": true, "ribu": true, "jt": true, "juta": true,
	"puluh": true, "belas": true, "ratus": true,
	"sepuluh": true, "sebelas": true, "seratus": true, "seribu": true, "sejuta": true,
}

// looksLikeTransaction is a cheap local pre-filter run before any backend or
// LLM call: every transaction needs an amount, so a message with no digit and
// no amount word ("halo", "makasih") can't be one and is answered locally.
// False positives (e.g. "jam 5") just fall through to the LLM as before.
func looksLikeTransaction(text string) bool {
	if strings.IndexFunc(text, unicode.IsDigit) >= 0 {
		return true
	}
	tokens := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	for _, t := range tokens {
		if amountWords[t] {
			return true
		}
	}
	return false
}
