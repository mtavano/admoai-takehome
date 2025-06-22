package api

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mtavano/admoai-takehome/internal/store"
)

// DashboardData contiene los datos para el template
type DashboardData struct {
	Ads         []*store.AdvertiseRecord
	TotalAds    int
	ActiveAds   int
	InactiveAds int
	ExpiredAds  int
}

// AdsDashboardHandler maneja el endpoint para mostrar el dashboard de anuncios
func AdsDashboardHandler(c *gin.Context, ctx *Context) (any, int, error) {
	// Obtener todos los anuncios
	var ads []*store.AdvertiseRecord
	err := ctx.Db.Select(&ads, "SELECT * FROM ads ORDER BY created_at DESC")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Calcular estadísticas
	stats := calculateStats(ads)

	// Preparar datos para el template
	data := DashboardData{
		Ads:         ads,
		TotalAds:    stats.TotalAds,
		ActiveAds:   stats.ActiveAds,
		InactiveAds: stats.InactiveAds,
		ExpiredAds:  stats.ExpiredAds,
	}

	// Cargar y renderizar el template
	tmpl, err := template.New("ads_table.html").Funcs(template.FuncMap{
		"formatTime": formatTime,
	}).ParseFiles("internal/api/templates/ads_table.html")
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Establecer el header de contenido HTML
	c.Header("Content-Type", "text/html; charset=utf-8")

	// Renderizar el template
	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return nil, http.StatusOK, nil
}

// Stats contiene las estadísticas de los anuncios
type Stats struct {
	TotalAds    int
	ActiveAds   int
	InactiveAds int
	ExpiredAds  int
}

// calculateStats calcula las estadísticas de los anuncios
func calculateStats(ads []*store.AdvertiseRecord) Stats {
	stats := Stats{}

	for _, ad := range ads {
		stats.TotalAds++

		// Calcular si está expirado
		if ad.ExpiresAt != nil {
			ad.CalculateAndSetExpired()
			if ad.Expired {
				stats.ExpiredAds++
			}
		}

		// Contar por estado
		switch ad.Status {
		case store.AdvertiseStatusActive:
			stats.ActiveAds++
		case store.AdvertiseStatusInactive:
			stats.InactiveAds++
		}
	}

	return stats
}

// formatTime formatea un timestamp Unix a una fecha legible
func formatTime(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("02/01/2006 15:04:05")
}
