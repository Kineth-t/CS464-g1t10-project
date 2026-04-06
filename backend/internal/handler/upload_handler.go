package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"log/slog"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadHandler handles media uploads to Cloudinary.
type UploadHandler struct{}

// NewUploadHandler creates a new UploadHandler.
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// Upload accepts a multipart/form-data POST with an "image" field,
// uploads it to Cloudinary, and returns the secure URL.
//
// @Summary      Upload product image
// @Description  Upload an image file to Cloudinary and receive a public URL.
// @Tags         admin
// @Accept       mpfd
// @Produce      json
// @Param        image  formData  file  true  "Image file"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /upload [post]
func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Limit to 10 MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		slog.Warn("upload failed: multipart form error", "error", err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "file too large or invalid multipart form"})
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		slog.Warn("upload failed: missing image field")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing 'image' field in form"})
		return
	}
	defer file.Close()

	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		slog.Error("upload failed: CLOUDINARY_URL not set")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "image upload not configured"})
		return
	}

	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		slog.Error("cloudinary connection error", "error", err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "cloudinary configuration error"})
		return
	}

	slog.Info("starting cloudinary upload", "filename", header.Filename, "size", header.Size)

	resp, err := cld.Upload.Upload(r.Context(), file, uploader.UploadParams{
		Folder: "ringr",
	})
	if err != nil {
		slog.Error("cloudinary upload failed", "error", err.Error(), "filename", header.Filename)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "upload failed: " + err.Error()})
		return
	}

	slog.Info("image upload successful", "url", resp.SecureURL, "public_id", resp.PublicID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": resp.SecureURL})
}
