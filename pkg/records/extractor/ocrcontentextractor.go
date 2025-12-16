package extractor

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kazemisoroush/assistant/pkg/records"
	"github.com/otiai10/gosseract/v2"
)

// OCRContentExtractor extracts records from images using OCR
type OCRContentExtractor struct {
	typeExtractor TypeExtractor
}

// NewOCRContentExtractor creates a new OCRExtractor instance
func NewOCRContentExtractor(typeExtractor TypeExtractor) ContentExtractor {
	return &OCRContentExtractor{
		typeExtractor: typeExtractor,
	}
}

// Extract processes raw content (image or text) and returns a Record
func (o *OCRContentExtractor) Extract(rawContent string) (records.Record, error) {
	now := time.Now()

	// 1) Try to OCR if rawContent looks like an image input; otherwise treat it as already-text.
	text, meta, err := o.toText(rawContent)
	if err != nil {
		return records.Record{}, fmt.Errorf("OCR extraction failed: %w", err)
	}

	// 2) Classify based on extracted text
	recordType := o.typeExtractor.GetType(text)

	rec := records.Record{
		ID:        fmt.Sprintf("ocr-%d", now.UnixNano()),
		Type:      recordType,
		Content:   text,
		CreatedAt: now,
		UpdatedAt: now,
		Metadata:  meta,
		Tags:      []string{"TBA"},
	}
	return rec, nil
}

// toText tries to OCR if rawContent is image-ish; otherwise returns rawContent as text.
// Metadata returned is useful for debugging (source/type, OCR used, etc.).
func (o *OCRContentExtractor) toText(rawContent string) (string, map[string]interface{}, error) {
	meta := map[string]interface{}{
		"source": "ocr",
	}

	s := strings.TrimSpace(rawContent)
	if s == "" {
		return "", meta, errors.New("rawContent is empty")
	}

	// Case A) data URL: data:image/png;base64,xxxx
	if looksLikeDataURL(s) {
		meta["input_kind"] = "data_url"
		mime, b64 := splitDataURL(s)
		meta["mime"] = mime

		imgBytes, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return "", meta, fmt.Errorf("failed to decode data URL base64: %w", err)
		}
		text, err := o.ocrBytesToText(imgBytes, mimeToExt(mime))
		if err != nil {
			return "", meta, err
		}
		meta["ocr_used"] = true
		return text, meta, nil
	}

	// Case B) looks like a file path to an image
	if looksLikeImagePath(s) {
		meta["input_kind"] = "file_path"
		text, err := o.ocrFileToText(s)
		if err != nil {
			return "", meta, err
		}
		meta["ocr_used"] = true
		return text, meta, nil
	}

	// Case C) raw base64 (no prefix) — common if you store images in DB
	// We only attempt this if it's plausibly base64 and fairly large.
	if looksLikeBase64ImageBlob(s) {
		meta["input_kind"] = "base64_blob"

		imgBytes, err := base64.StdEncoding.DecodeString(stripBase64Whitespace(s))
		if err != nil {
			// If decode fails, treat as text (some OCR output can look base64-ish)
			meta["ocr_used"] = false
			return rawContent, meta, nil
		}

		// We don’t know the type; assume png by default (you can sniff magic bytes if you want).
		// Better: sniff header and choose ext. We'll do a tiny sniff.
		ext := sniffImageExt(imgBytes)
		text, err := o.ocrBytesToText(imgBytes, ext)
		if err != nil {
			return "", meta, err
		}
		meta["ocr_used"] = true
		meta["sniffed_ext"] = ext
		return text, meta, nil
	}

	// Case D) already text
	meta["input_kind"] = "text"
	meta["ocr_used"] = false
	return rawContent, meta, nil
}

func looksLikeDataURL(s string) bool {
	return strings.HasPrefix(s, "data:image/") && strings.Contains(s, ";base64,")
}

func splitDataURL(s string) (mime string, b64 string) {
	// format: data:image/png;base64,XXXX
	// mime part between "data:" and ";base64,"
	prefix := "data:"
	i := strings.Index(s, ";base64,")
	mime = s[len(prefix):i]
	b64 = s[i+len(";base64,"):]
	return mime, b64
}

func looksLikeImagePath(s string) bool {
	ext := strings.ToLower(filepath.Ext(s))
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp" || ext == ".tif" || ext == ".tiff" {
		_, err := os.Stat(s)
		return err == nil
	}
	return false
}

func looksLikeBase64ImageBlob(s string) bool {
	// Heuristic: base64 strings are usually long, only base64 chars, and length%4==0 often.
	// We keep this conservative so we don't mis-detect random text.
	ss := stripBase64Whitespace(s)
	if len(ss) < 500 { // images are usually larger than this
		return false
	}
	if len(ss)%4 != 0 {
		return false
	}
	for _, r := range ss {
		if (r < 'a' || r > 'z') &&
			(r < 'A' || r > 'Z') &&
			(r < '0' || r > '9') &&
			r != '+' && r != '/' && r != '=' {
			return false
		}
	}
	return true
}

func stripBase64Whitespace(s string) string {
	// remove common whitespace/newlines
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}
func (o *OCRContentExtractor) ocrBytesToText(img []byte, ext string) (string, error) {
	// Tesseract/gosseract prefers a file path, so we write a temp file.
	tmpDir := os.TempDir()
	if ext == "" {
		ext = ".png"
	}
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("ocr-%d%s", time.Now().UnixNano(), ext))

	if err := os.WriteFile(tmpFile, img, 0600); err != nil {
		return "", fmt.Errorf("failed to write temp image: %w", err)
	}
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	return o.ocrFileToText(tmpFile)
}

func mimeToExt(mime string) string {
	switch strings.ToLower(mime) {
	case "image/png":
		return ".png"
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	default:
		return ".png"
	}
}

func sniffImageExt(b []byte) string {
	// Minimal magic-byte sniffing
	if len(b) >= 8 && string(b[:8]) == "\x89PNG\r\n\x1a\n" {
		return ".png"
	}
	if len(b) >= 3 && b[0] == 0xFF && b[1] == 0xD8 && b[2] == 0xFF {
		return ".jpg"
	}
	// WEBP: "RIFF....WEBP"
	if len(b) >= 12 && string(b[:4]) == "RIFF" && string(b[8:12]) == "WEBP" {
		return ".webp"
	}
	return ".png"
}

func (o *OCRContentExtractor) ocrFileToText(path string) (string, error) {
	client := gosseract.NewClient()
	defer func() {
		if err := client.Close(); err != nil {
			fmt.Printf("warning: failed to close tesseract client: %v\n", err)
		}
	}()

	// Optional: set languages. Requires language packs installed.
	// client.SetLanguage("eng") // or "eng+fas" if you install Persian traineddata
	if err := client.SetImage(path); err != nil {
		return "", fmt.Errorf("failed to set image: %w", err)
	}
	return client.Text()
}
