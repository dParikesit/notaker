package routes

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"server/app/models"
	mid "server/pkg/middleware"
	"server/pkg/utils"
	"server/platform/database"
)

type GoogleToken struct {
	ClientId   string `json:"clientId"`
	Credential string `json:"credential"`
	Select_by  string `json:"select_by"`
}

func Login(c *fiber.Ctx) error {
	body := new(GoogleToken)
	if err := c.BodyParser(body); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Bad Token",
		})

		return err
	}

	if body.ClientId != utils.Dotenv("OAUTH") {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid Client ID",
		})
	}

	google, errGoogle := utils.ExtractGoogle(body.Credential)
	if !errGoogle {
		c.SendStatus(fiber.StatusInternalServerError)
	}

	sess, err := mid.Store.Get(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Session error",
		})
		log.Panic(err)
	}

	sess.Set("email", google.Email)
	if err := sess.Save(); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Session save error",
		})
		log.Panic(err)
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"error":   false,
		"picture": google.Picture,
	})
}

func Logout(c*fiber.Ctx) error {
	sess, err := mid.Store.Get(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Session error",
		})
  	log.Panic(err)
  }

	if err := sess.Destroy(); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Session destroy error",
		})
  	log.Panic(err)
  }

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"message": "Logout successful",
	})

	return nil
}

func Add(c *fiber.Ctx) error {
	var note models.Note

	if form, err := c.MultipartForm(); err == nil {

		if title := form.Value["title"]; len(title) > 0 {
			note.Title = title[0]
		}
		if text := form.Value["text"]; len(text) > 0 {
			note.Text = text[0]
		}

		// Get all files from "documents" key:
		images := form.File["images"]
		// => []*multipart.FileHeader

		var imagePaths []string
		// Loop through files:
		for _, image := range images {
			// Save the files to disk:
			if err := c.SaveFile(image, fmt.Sprintf("./images/%s", image.Filename)); err != nil {
				return err
			}
			imagePaths = append(imagePaths, (image.Filename))
		}
		note.Images = imagePaths
	}

	noteResult, err := database.Notes.InsertOne(c.Context(), note)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"note":  noteResult,
	})

	return nil
}

func GetNotes(c *fiber.Ctx) error {
	cursor, err := database.Notes.Find(c.Context(), bson.M{})
	if err != nil {
		log.Panic(err)
	}
	var notes []bson.M
	if err = cursor.All(c.Context(), &notes); err != nil {
		log.Panic(err)
	}

	c.Status(fiber.StatusFound).JSON(fiber.Map{
		"error": false,
		"notes": notes,
	})

	return nil
}

func GetNote(c *fiber.Ctx) error {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	var note bson.M
	if err := database.Notes.FindOne(c.Context(), bson.M{"_id": id}).Decode(&note); err != nil {
		log.Panic(err)
		c.SendStatus(fiber.StatusNotFound)
	}
	c.Status(fiber.StatusFound).JSON(fiber.Map{
		"error": false,
		"note":  note,
	})
	return nil
}

func PatchNote(c *fiber.Ctx) error {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	var update models.Note
	if form, err := c.MultipartForm(); err == nil {

		if title := form.Value["title"]; len(title) > 0 {
			update.Title = title[0]
		}
		if text := form.Value["text"]; len(text) > 0 {
			update.Text = text[0]
		}

		// Get all files from "documents" key:
		if images := form.File["images"]; len(images) > 0 {
			var imagePaths []string
			// Loop through files:
			for _, image := range images {
				// Save the files to disk:
				if err := c.SaveFile(image, fmt.Sprintf("./images/%s", image.Filename)); err != nil {
					c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error":  true,
						"status": "Save error",
					})
					return err
				}
				imagePaths = append(imagePaths, (image.Filename))
			}
			update.Images = imagePaths
		}
	}
	var beforeDocument models.Note
	if len(update.Images) > 0 {
		result := database.Notes.FindOneAndUpdate(c.Context(), bson.M{"_id": id}, bson.D{
			{"$set", bson.D{{"title", update.Title}, {"text", update.Text}, {"images", update.Images}}},
		}).Decode(&beforeDocument)
		if result != nil {
			c.Status(fiber.StatusNotModified).JSON(fiber.Map{
				"error":  true,
				"status": "Patch failed",
			})
			log.Panic(result)
		}
	} else {
		result := database.Notes.FindOneAndUpdate(c.Context(), bson.M{"_id": id}, bson.D{
			{"$set", bson.D{{"title", update.Title}, {"text", update.Text}}},
		}).Decode(&beforeDocument)

		if result != nil {
			c.Status(fiber.StatusNotModified).JSON(fiber.Map{
				"error":  true,
				"status": "Patch failed",
			})
			log.Panic(result)
		}
	}
	imgBefore := beforeDocument.Images

	var exists bool
	for _, before := range imgBefore {
		exists = false
		for _, after := range update.Images {
			if before == after {
				exists = true
			}
		}
		if !exists {
			os.Remove(filepath.Join("images", before))
		}
	}
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
	})
	return nil
}

func DeleteNote(c *fiber.Ctx) error {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	var beforeDocument models.Note
	err := database.Notes.FindOneAndDelete(c.Context(), bson.M{"_id": id}).Decode(&beforeDocument)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":  true,
			"status": "Note not found",
		})
		log.Panic(err)
	}

	for _, before := range beforeDocument.Images {
		os.Remove(filepath.Join("images", before))
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":  false,
		"result": "Delete successful",
	})
	return nil
}
