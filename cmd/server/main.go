package main

import (
	"bookmarks/internal/archiver"
	"bookmarks/internal/models"
	"bookmarks/internal/storage"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Initialize Database
	db, err := storage.New("bookmarks.db")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Archiver
	arc := archiver.New()

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Setup Templates
	t := &Template{
		templates: template.Must(template.New("").Funcs(template.FuncMap{
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
		}).ParseGlob("web/templates/*.html")),
	}
	e.Renderer = t

	// API Routes (for iOS Shortcut)
	e.POST("/api/bookmarks", func(c echo.Context) error {
		type Request struct {
			URL     string `json:"url" form:"url"`
			Comment string `json:"comment" form:"comment"`
		}
		var req Request
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Sanitize URL (iOS Shortcuts sometimes sends multiple URLs separated by newline)
		req.URL = strings.TrimSpace(req.URL)
		if idx := strings.Index(req.URL, "\n"); idx != -1 {
			req.URL = strings.TrimSpace(req.URL[:idx])
		}

		// Archive the URL
		bookmark, err := arc.Archive(req.URL)
		if err != nil {
			// Fallback if archiving fails: just save the URL
			bookmark = &models.Bookmark{
				URL:   req.URL,
				Title: req.URL, // Temporary title
			}
			log.Printf("Failed to archive %s: %v", req.URL, err)
		}

		bookmark.Comment = req.Comment

		// Save to DB
		id, err := db.CreateBookmark(bookmark)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save bookmark"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"id": id})
	})

	// Web Routes
	e.GET("/", func(c echo.Context) error {
		bookmarks, err := db.ListBookmarks()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"Bookmarks": bookmarks,
		})
	})

	e.GET("/bookmarks/:id", func(c echo.Context) error {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		bookmark, err := db.GetBookmark(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Bookmark not found")
		}

		if bookmark.Deleted {
			return c.Render(http.StatusOK, "deleted.html", map[string]interface{}{
				"Bookmark": bookmark,
			})
		}

		return c.Render(http.StatusOK, "detail.html", map[string]interface{}{
			"Bookmark": bookmark,
		})
	})

	e.POST("/bookmarks/:id/delete", func(c echo.Context) error {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		if err := db.DeleteBookmark(id); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to delete bookmark")
		}
		return c.Redirect(http.StatusSeeOther, "/")
	})

	e.POST("/bookmarks/:id/comment", func(c echo.Context) error {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		comment := c.FormValue("comment")
		if err := db.UpdateComment(id, comment); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update comment")
		}
		return c.Redirect(http.StatusSeeOther, "/bookmarks/"+c.Param("id"))
	})

	e.Logger.Fatal(e.Start(":8080"))
}
