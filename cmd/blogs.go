package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/masonictemple4/masonictempl/db/models"
	"github.com/masonictemple4/masonictempl/internal/dtos"
	"github.com/masonictemple4/masonictempl/internal/parser"
)

var blogsCmd = &cobra.Command{
	Use:   "blog",
	Short: "Manage blogs within the masonictempl project.",
	Long: `Manage blogs within the masonictempl project. For example:
You canuse this command to list, create, update, and delete blogs.
The default will be to list all blogs.
masonictempl blogs [command]`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var blogCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new blog.",
	Long: `Create a new blog. For example:
masonictempl blog create <file path>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(blogsCmd)
}

func createBlog(ctx context.Context, path string, flags *pflag.FlagSet) error {
	// What if instead of passing a path to parser here and eventually would have to be
	// the filestore too if i read file here and pass bytes to the parsr and writer.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result dtos.BlogInput
	err = parser.ParseFile(path, &result)
	if err != nil {
		return err
	}

	fmt.Printf("\nThe parsed result is:  %+v\n", result)

	post := &models.Blog{
		Bucketname: os.Getenv("STORAGE_BUCKET"),
		State:      models.BlogStateDraft,
	}

	err = post.FromBlogInput(DB, &result)
	if err != nil {
		return err
	}

	// generate a slug because now we have a title.
	post.Slug = post.GenerateSlug("")

	fileHandler := filestore.NewGCPStore(false, 0)

	post.Docpath, err = post.GenerateDocPath()
	if err != nil {
		return err
	}

	written, err := fileHandler.Write(ctx, post.Docpath, data)
	if err != nil || len(data) != int(written) {
		return err
	}

	updateBody := map[string]any{"contenturl": post.GenerateContentUrl(), "docpath": post.Docpath, "state": models.BlogStatePublished, "slug": post.Slug}
	err = post.Update(DB, int(post.ID), updateBody)
	if err != nil {
		return err
	}

	return nil
}

func updateBlog(ctx context.Context, bid int, path string) error {
	var blog models.Blog

	err := blog.FindByID(DB, bid, nil)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var result dtos.BlogInput
	err = parser.ParseFile(path, &result)
	if err != nil {
		return err
	}

	fmt.Printf("\nThe parsed result is:  %+v\n", result)

	err = blog.FromBlogInput(DB, &result)
	if err != nil {
		return err
	}

	// generate a slug because now we have a title.
	blog.Slug = blog.GenerateSlug("")

	fileHandler := filestore.NewGCPStore(false, 0)

	// TODO: Should probably delete here to be safe we're not
	// gathering too unused files.
	blog.Docpath, err = blog.GenerateDocPath()
	if err != nil {
		return err
	}

	written, err := fileHandler.Write(ctx, blog.Docpath, data)
	if err != nil || len(data) != int(written) {
		return err
	}

	updateBody := map[string]any{
		"contenturl": blog.GenerateContentUrl(),
		"docpath":    blog.Docpath,
		"state":      models.BlogStatePublished,
		"slug":       blog.Slug,
	}
	err = blog.Update(DB, int(blog.ID), updateBody)
	if err != nil {
		return err
	}

	return nil

}
