// Package service orchestrates one Telegram message end-to-end: resolve the
// user via chat ID, fetch their categories, ask the LLM to parse the message,
// re-validate its answer against the real category list, then write the
// transaction through the backend. Every method returns a user-facing Telegram
// reply string on the happy AND the handled-error path (e.g. "not linked
// yet") — only unexpected failures are returned as Go errors.
package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/afandimsr/cashbook-bot/internal/backend"
	"github.com/afandimsr/cashbook-bot/internal/llm"
)

// notLinkedMessage is shown whenever a chat isn't linked to a CashBook
// account yet, kept as one constant so every method gives the same complete
// instructions instead of drifting into shorter variants over time.
const notLinkedMessage = "Akun Telegram kamu belum terhubung ke CashBook. Buka aplikasi CashBook, buat kode link, lalu kirim /link <kode> di sini."

// BackendClient is satisfied by *backend.Client; declared here so tests can
// substitute a mock without hitting the real CashBook API.
type BackendClient interface {
	GetContext(ctx context.Context, chatID int64) (*backend.Context, error)
	LinkAccount(ctx context.Context, chatID int64, username, code string) error
	CreateTransaction(ctx context.Context, req backend.CreateTransactionRequest) error
	GetSummary(ctx context.Context, chatID int64) (*backend.DashboardSummary, error)
	GetMonthlyReport(ctx context.Context, chatID int64, month, year int) (*backend.MonthlyReport, error)
	GetBudgetStatus(ctx context.Context, chatID int64, month, year int) ([]backend.BudgetStatus, error)
	ListLinks(ctx context.Context) ([]backend.TelegramLink, error)
}

// LLMClient is satisfied by *llm.Client; declared here so tests can substitute
// a mock without calling the real Gemini API.
type LLMClient interface {
	ParseTransaction(ctx context.Context, message string, categories []backend.Category) (*llm.ParsedTransaction, error)
}

type Bot struct {
	backend BackendClient
	llm     LLMClient
}

func New(backendClient BackendClient, llmClient LLMClient) *Bot {
	return &Bot{backend: backendClient, llm: llmClient}
}

// Link consumes a one-time code the user generated in the CashBook app to
// bind this Telegram chat to their account.
func (b *Bot) Link(ctx context.Context, chatID int64, username, code string) (string, error) {
	err := b.backend.LinkAccount(ctx, chatID, username, code)
	if err == nil {
		return "Akun Telegram kamu berhasil terhubung ke CashBook. Sekarang kamu bisa langsung catat transaksi, misalnya: \"Beli bensin 50rb\".", nil
	}

	var apiErr *backend.APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == 400 {
		return "Kode tidak valid, sudah kedaluwarsa, atau sudah dipakai. Buat kode baru dari aplikasi CashBook lalu kirim /link <kode>.", nil
	}
	return "", err
}

// WhoAmI reports which CashBook account this Telegram chat is linked to, or
// prompts the user to link one if it isn't yet.
// lastName is optional in Telegram and may be empty.
func (b *Bot) WhoAmI(ctx context.Context, chatID int64, firstName, lastName string) (string, error) {
	name := displayName(firstName, lastName)
	botCtx, err := b.backend.GetContext(ctx, chatID)
	if err != nil {
		var apiErr *backend.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			return fmt.Sprintf("Halo %s! %s", name, notLinkedMessage), nil
		}
		return "", err
	}
	return fmt.Sprintf("Halo %s! Kamu sudah terhubung ke CashBook sebagai %s (%s).", name, botCtx.Name, botCtx.Email), nil
}

// HandleMessage parses a free-text message into a transaction and records it.
// Messages that can't be a transaction get HelpMessage without touching the
// backend or the LLM.
func (b *Bot) HandleMessage(ctx context.Context, chatID int64, message string) (string, error) {
	if !looksLikeTransaction(message) {
		return HelpMessage, nil
	}

	botCtx, err := b.backend.GetContext(ctx, chatID)
	if err != nil {
		var apiErr *backend.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			return notLinkedMessage, nil
		}
		return "", err
	}

	if len(botCtx.Categories) == 0 {
		return "Kamu belum punya kategori di CashBook. Buat kategori dulu di aplikasi sebelum mencatat transaksi lewat bot.", nil
	}

	parsed, err := b.llm.ParseTransaction(ctx, message, botCtx.Categories)
	if errors.Is(err, llm.ErrNoFunctionCall) || errors.Is(err, llm.ErrInvalidOutput) {
		return "Maaf, aku tidak bisa memahami itu sebagai transaksi. Coba tulis lebih spesifik, misalnya \"Beli bensin 50rb pakai dompet digital\".", nil
	}
	if err != nil {
		// Every LLM provider was unreachable/rate-limited: not the user's fault,
		// so let the handler reply with the generic "try again later".
		return "", err
	}

	category, ok := matchCategory(botCtx.Categories, parsed.CategoryName)
	if !ok {
		return fmt.Sprintf(
			"Aku mendeteksi transaksi sebesar %s, tapi kategori \"%s\" tidak cocok dengan kategori kamu. Coba sebutkan kategorinya lebih jelas.",
			formatRupiah(parsed.Amount), parsed.CategoryName,
		), nil
	}

	err = b.backend.CreateTransaction(ctx, backend.CreateTransactionRequest{
		ChatID:     chatID,
		CategoryID: category.ID,
		Amount:     parsed.Amount,
		Type:       category.Type, // trust the category's own type, never the model's guess
		Note:       parsed.Note,
		Date:       parsed.Date,
	})
	if err != nil {
		return "", err
	}

	verb := "Pemasukan"
	if category.Type == "expense" {
		verb = "Pengeluaran"
	}
	noteSuffix := ""
	if parsed.Note != "" {
		noteSuffix = fmt.Sprintf(" (%s)", parsed.Note)
	}
	return fmt.Sprintf("%s tercatat: %s untuk %s%s.", verb, formatRupiah(parsed.Amount), category.Name, noteSuffix), nil
}

// Summary returns the linked user's dashboard totals.
func (b *Bot) Summary(ctx context.Context, chatID int64) (string, error) {
	summary, err := b.backend.GetSummary(ctx, chatID)
	if err != nil {
		var apiErr *backend.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			return notLinkedMessage, nil
		}
		return "", err
	}
	return fmt.Sprintf(
		"Ringkasan kamu:\nPemasukan: %s\nPengeluaran: %s\nSaldo: %s",
		formatRupiah(summary.TotalIncome), formatRupiah(summary.TotalExpense), formatRupiah(summary.Balance),
	), nil
}

// MonthlyReport formats the per-category income and expense breakdown for the
// given month/year. It's used both by the scheduler (automatic report on the
// 1st of the month, for the month that just closed) and the on-demand /report
// command.
func (b *Bot) MonthlyReport(ctx context.Context, chatID int64, month, year int) (string, error) {
	report, err := b.backend.GetMonthlyReport(ctx, chatID, month, year)
	if err != nil {
		var apiErr *backend.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			return notLinkedMessage, nil
		}
		return "", err
	}

	monthName := time.Month(month).String()
	if len(report.Income) == 0 && len(report.Expense) == 0 {
		return fmt.Sprintf("Laporan %s %d: tidak ada transaksi tercatat.", monthName, year), nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Laporan %s %d:\n\n", monthName, year)
	totalIncome := writeReportSection(&sb, "Pemasukan", report.Income)
	sb.WriteString("\n")
	totalExpense := writeReportSection(&sb, "Pengeluaran", report.Expense)
	fmt.Fprintf(&sb, "\nSelisih: %s", formatRupiah(totalIncome-totalExpense))
	return sb.String(), nil
}

// writeReportSection writes one titled category list plus its total and
// returns that total. The backend already sorts rows by amount.
func writeReportSection(sb *strings.Builder, title string, rows []backend.CategoryReport) float64 {
	fmt.Fprintf(sb, "%s:\n", title)
	if len(rows) == 0 {
		sb.WriteString("- tidak ada\n")
	}
	var total float64
	for _, c := range rows {
		total += c.TotalAmount
		fmt.Fprintf(sb, "- %s: %s\n", c.CategoryName, formatRupiah(c.TotalAmount))
	}
	fmt.Fprintf(sb, "Total %s: %s\n", strings.ToLower(title), formatRupiah(total))
	return total
}

// Categories lists the linked user's categories grouped by type, reusing the
// same GetContext call that transaction parsing uses.
func (b *Bot) Categories(ctx context.Context, chatID int64) (string, error) {
	botCtx, err := b.backend.GetContext(ctx, chatID)
	if err != nil {
		var apiErr *backend.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			return notLinkedMessage, nil
		}
		return "", err
	}

	if len(botCtx.Categories) == 0 {
		return "Kamu belum punya kategori di CashBook. Buat kategori dulu di aplikasi.", nil
	}

	var income, expense []string
	for _, c := range botCtx.Categories {
		if c.Type == "income" {
			income = append(income, c.Name)
		} else {
			expense = append(expense, c.Name)
		}
	}

	var sb strings.Builder
	sb.WriteString("Kategori kamu:\n\n")
	writeCategorySection(&sb, "Pemasukan", income)
	sb.WriteString("\n")
	writeCategorySection(&sb, "Pengeluaran", expense)
	sb.WriteString("\nCatat transaksi dengan menyebut kategorinya, misalnya \"Makan siang 25rb\".")
	return sb.String(), nil
}

func writeCategorySection(sb *strings.Builder, title string, names []string) {
	fmt.Fprintf(sb, "%s:\n", title)
	if len(names) == 0 {
		sb.WriteString("- tidak ada\n")
		return
	}
	sort.Slice(names, func(i, j int) bool { return strings.ToLower(names[i]) < strings.ToLower(names[j]) })
	for _, n := range names {
		fmt.Fprintf(sb, "- %s\n", n)
	}
}

// Budgets shows each category budget for the month with how much is used,
// flagging the ones already over their limit.
func (b *Bot) Budgets(ctx context.Context, chatID int64, month, year int) (string, error) {
	statuses, err := b.backend.GetBudgetStatus(ctx, chatID, month, year)
	if err != nil {
		var apiErr *backend.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			return notLinkedMessage, nil
		}
		return "", err
	}

	monthName := time.Month(month).String()
	if len(statuses) == 0 {
		return fmt.Sprintf("Belum ada anggaran untuk %s %d. Atur anggaran per kategori di aplikasi CashBook.", monthName, year), nil
	}

	var totalLimit, totalSpent float64
	var sb strings.Builder
	fmt.Fprintf(&sb, "Anggaran %s %d:\n\n", monthName, year)
	for _, s := range statuses {
		totalLimit += s.Limit
		totalSpent += s.Spent
		fmt.Fprintf(&sb, "%s %s: %s / %s (%.0f%%)\n", budgetMarker(s), s.CategoryName,
			formatRupiah(s.Spent), formatRupiah(s.Limit), s.Percentage)
		if s.IsOver {
			fmt.Fprintf(&sb, "   lewat %s\n", formatRupiah(-s.Remaining))
		} else {
			fmt.Fprintf(&sb, "   sisa %s\n", formatRupiah(s.Remaining))
		}
	}
	fmt.Fprintf(&sb, "\nTotal terpakai: %s dari %s", formatRupiah(totalSpent), formatRupiah(totalLimit))
	return sb.String(), nil
}

// budgetMarker uses the same thresholds as the web app's budget progress bar
// (frontend BudgetPage): over the limit, above 80%, or fine.
func budgetMarker(s backend.BudgetStatus) string {
	switch {
	case s.IsOver:
		return "🔴"
	case s.Percentage > 80:
		return "🟡"
	default:
		return "🟢"
	}
}

func displayName(firstName, lastName string) string {
	return strings.TrimSpace(firstName + " " + lastName)
}

func matchCategory(categories []backend.Category, name string) (backend.Category, bool) {
	for _, c := range categories {
		if strings.EqualFold(c.Name, name) {
			return c, true
		}
	}
	return backend.Category{}, false
}
