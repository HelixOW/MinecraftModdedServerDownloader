package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/download", func(c *fiber.Ctx) error {
		// Creating a buffer to hold the zip archive.
		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)

		// This is the folder we will zip and send.
		folderToZip := "/home/minecraft/server"

		err := filepath.Walk(folderToZip, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Creating the new file in the zip archive.
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			header.Name = path
			header.Method = zip.Deflate

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			// If it's not a directory, open the file to write its contents to the zip archive.
			fileToZip, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fileToZip.Close()

			_, err = io.Copy(writer, fileToZip)
			return err
		})
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError).SendString("Something went wrong")
			return err
		}

		// Make sure to check the error on Close.
		err = zipWriter.Close()
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError).SendString("Something went wrong")
			return err
		}

		// Set the headers and send the buffer.
		c.Type("application/zip", "utf-8")
		c.Set(fiber.HeaderContentDisposition, `attachment; filename="archive.zip"`)
		return c.SendStream(bytes.NewReader(buf.Bytes()), buf.Len())
	})

	app.Listen(":12345")
}
