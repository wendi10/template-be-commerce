package doku

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/template-be-commerce/internal/domain"
	"github.com/template-be-commerce/internal/payment"
	"github.com/template-be-commerce/pkg/logger"
	"go.uber.org/zap"
)

const providerName = domain.PaymentProviderDoku

type Provider struct {
	clientID  string
	secretKey string
	baseURL   string
	httpClient *http.Client
}

func NewProvider(clientID, secretKey, baseURL string) *Provider {
	return &Provider{
		clientID:  clientID,
		secretKey: secretKey,
		baseURL:   baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *Provider) ProviderName() domain.PaymentProvider {
	return providerName
}

// CreateTransaction calls the DOKU payment creation API.
// Reference: https://developers.doku.com/accept-payment/
func (p *Provider) CreateTransaction(ctx context.Context, req payment.CreateTransactionRequest) (*payment.TransactionResponse, error) {
	requestID := fmt.Sprintf("%s-%d", req.OrderNumber, time.Now().UnixMilli())
	requestDate := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	payload := map[string]interface{}{
		"order": map[string]interface{}{
			"amount":      req.Amount,
			"invoice_number": req.OrderNumber,
			"currency":    "IDR",
			"callback_url": req.CallbackURL,
			"auto_redirect": false,
		},
		"payment": map[string]interface{}{
			"payment_due_date": 60, // minutes
		},
		"customer": map[string]interface{}{
			"name":  req.CustomerName,
			"email": req.CustomerEmail,
			"phone": req.CustomerPhone,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("doku: marshal request: %w", err)
	}

	signature := p.generateSignature(requestID, requestDate, string(body))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/checkout/v1/payment", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("doku: create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Client-Id", p.clientID)
	httpReq.Header.Set("Request-Id", requestID)
	httpReq.Header.Set("Request-Timestamp", requestDate)
	httpReq.Header.Set("Signature", "HMACSHA256="+signature)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("doku: send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	logger.Debug("doku create transaction response", zap.Int("status", resp.StatusCode), zap.String("body", string(respBody)))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("doku: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Response struct {
			Checkout struct {
				URL string `json:"url"`
			} `json:"checkout"`
			Invoice struct {
				ID string `json:"id"`
			} `json:"invoice"`
		} `json:"response"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("doku: parse response: %w", err)
	}

	expiredAt := time.Now().Add(60 * time.Minute).Unix()
	return &payment.TransactionResponse{
		TransactionID: result.Response.Invoice.ID,
		PaymentURL:    result.Response.Checkout.URL,
		ExpiredAt:     expiredAt,
		RawResponse:   string(respBody),
	}, nil
}

// HandleCallback parses and validates the DOKU payment notification.
func (p *Provider) HandleCallback(ctx context.Context, payload []byte, headers map[string]string) (*payment.CallbackResult, error) {
	var data struct {
		Order struct {
			InvoiceNumber string `json:"invoice_number"`
			Amount        string `json:"amount"`
		} `json:"order"`
		Transaction struct {
			Status string `json:"status"`
			ID     string `json:"id"`
		} `json:"transaction"`
		Security struct {
			CheckSum string `json:"check_sum"`
		} `json:"security"`
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, fmt.Errorf("doku callback: unmarshal: %w", err)
	}

	// Verify checksum
	expectedChecksum := p.computeCallbackChecksum(data.Order.Amount, data.Order.InvoiceNumber, data.Transaction.Status)
	if data.Security.CheckSum != expectedChecksum {
		return nil, fmt.Errorf("doku callback: invalid checksum")
	}

	status := mapDokuStatus(data.Transaction.Status)

	return &payment.CallbackResult{
		TransactionID: data.Transaction.ID,
		OrderID:       data.Order.InvoiceNumber,
		Status:        status,
		Amount:        data.Order.Amount,
		RawPayload:    string(payload),
	}, nil
}

// generateSignature computes the HMAC-SHA256 request signature.
func (p *Provider) generateSignature(requestID, requestDate, body string) string {
	digest := sha256.Sum256([]byte(body))
	digestHex := hex.EncodeToString(digest[:])

	component := fmt.Sprintf("Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s\nDigest:%s",
		p.clientID, requestID, requestDate, "/checkout/v1/payment", digestHex)

	mac := sha256.New()
	mac.Write([]byte(p.secretKey))
	mac.Write([]byte(component))
	return hex.EncodeToString(mac.Sum(nil))
}

// computeCallbackChecksum verifies callback data integrity.
func (p *Provider) computeCallbackChecksum(amount, invoiceNumber, status string) string {
	raw := fmt.Sprintf("%s%s%s%s", p.clientID, amount, invoiceNumber, status)
	h := sha256.Sum256([]byte(raw + p.secretKey))
	return hex.EncodeToString(h[:])
}

func mapDokuStatus(dokuStatus string) domain.PaymentStatus {
	switch dokuStatus {
	case "SUCCESS":
		return domain.PaymentStatusSuccess
	case "FAILED", "DECLINED":
		return domain.PaymentStatusFailed
	case "EXPIRED":
		return domain.PaymentStatusExpired
	default:
		return domain.PaymentStatusPending
	}
}
