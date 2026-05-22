// printer/qr.go
package printer

import (
    "encoding/base64"
    "github.com/skip2/go-qrcode"
)

// GenerateQRDataURL returns a base64-encoded data URL of a QR code.
func GenerateQRDataURL(content string) (string, error) {
    png, err := qrcode.Encode(content, qrcode.Medium, 256)
    if err != nil {
        return "", err
    }
    return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png), nil
}